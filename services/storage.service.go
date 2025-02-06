package services

import (
	"context"
	"errors"
	models "health/models/db"
	"log"
	"sync"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitMongoDB initializes the mgm library with the MongoDB URI and database name.
// It is called once during application startup.
func InitMongoDB() {
	// Setup the mgm default config
	err := mgm.SetDefaultConfig(nil, Config.MongodbDatabase, options.Client().ApplyURI(Config.MongodbUri))
	if err != nil {
		panic(err)
	}

	log.Println("Connected to MongoDB!")
}

var redisDefaultClient *redis.Client
var redisDefaultOnce sync.Once

var redisCache *cache.Cache
var redisCacheOnce sync.Once

// GetRedisDefaultClient returns the default Redis client instance.
// The client is created during the first call to this function.
// Subsequent calls will return the same instance.
// The client is configured with the default address specified in the configuration.
func GetRedisDefaultClient() *redis.Client {
	redisDefaultOnce.Do(func() {
		redisDefaultClient = redis.NewClient(&redis.Options{
			Addr: Config.RedisDefaultAddr,
		})
	})

	return redisDefaultClient
}

// GetRedisCache returns the default Redis cache instance.
// The cache is created during the first call to this function with a Redis client
// and a local cache using a TinyLFU algorithm. Subsequent calls will return
// the same cache instance.
func GetRedisCache() *cache.Cache {
	redisCacheOnce.Do(func() {
		redisCache = cache.New(&cache.Options{
			Redis:      GetRedisDefaultClient(),
			LocalCache: cache.NewTinyLFU(1000, time.Minute),
		})
	})

	return redisCache
}

// RedisCacheConnection attempts to ping the Redis server using the default client instance.
// If the ping fails, it will panic. Otherwise, it will log a message indicating that the connection was successful.
func CheckRedisCacheConnection() {
	redisClient := GetRedisDefaultClient()
	err := redisClient.Ping(context.Background()).Err()

	if err != nil {
		panic(err)
	}

	log.Println("Connected to Redis!")
}

// getNoteCacheKey generates a cache key for a note belonging to a specific user.
// The key is constructed using the user's ObjectID and the note's ObjectID,
// and is formatted as "req.cache.note:<userId>:<noteId>".

func getNoteCacheKey(userId primitive.ObjectID, noteId primitive.ObjectID) string {
	return "req.cache.note:" + userId.Hex() + ":" + noteId.Hex()
}

// CacheOneNote stores a single note in Redis cache with a TTL of 1 minute,
// using a cache key constructed from the user's ObjectID and the note's ObjectID.
// The function does nothing if the UseRedis configuration option is disabled.
func CacheOneNote(userId primitive.ObjectID, note *models.Note) {
	if !Config.UseRedis {
		return
	}

	noteCacheKey := getNoteCacheKey(userId, note.ID)

	_ = GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   noteCacheKey,
		Value: note,
		TTL:   time.Minute,
	})
}

func GetNoteFromCache(userId primitive.ObjectID, noteId primitive.ObjectID) (*models.Note, error) {
	if !Config.UseRedis {
		return nil, errors.New("no redis client, set USE_REDIS in .env")
	}
	note := &models.Note{}
	noteCachekey := getNoteCacheKey(userId, noteId)
	err := GetRedisCache().Get(context.TODO(), noteCachekey, note)
	return note, err
}
