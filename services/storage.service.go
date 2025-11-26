package services

import (
	"context"
	"errors"
	"fmt"
	models "health/models/db"
	"log"
	"sync"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var dbOnce sync.Once

// ConnectDB creates a new database connection (for CLI tools)
func ConnectDB(dsn string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

// InitMySQL initializes the GORM library with the MySQL connection string.
// It is called once during application startup.
func InitMySQL() {
	dbOnce.Do(func() {
		charset := Config.MySQLCharset
		if charset == "" {
			charset = "utf8mb4"
		}

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
			Config.MySQLUser,
			Config.MySQLPassword,
			Config.MySQLHost,
			Config.MySQLPort,
			Config.MySQLDatabase,
			charset,
		)

		var err error
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(fmt.Sprintf("Failed to connect to MySQL: %v", err))
		}

		// Run manual migrations
		if err = RunMigrations(DB); err != nil {
			panic(fmt.Sprintf("Failed to run migrations: %v", err))
		}

		log.Println("Connected to MySQL!")
	})
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
// The key is constructed using the user's ID and the note's ID,
// and is formatted as "req.cache.note:<userId>:<noteId>".

func getNoteCacheKey(userId uint, noteId uint) string {
	return fmt.Sprintf("req.cache.note:%d:%d", userId, noteId)
}

// CacheOneNote stores a single note in Redis cache with a TTL of 1 minute,
// using a cache key constructed from the user's ID and the note's ID.
// The function does nothing if the UseRedis configuration option is disabled.
func CacheOneNote(userId uint, note *models.Note) {
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

func GetNoteFromCache(userId uint, noteId uint) (*models.Note, error) {
	if !Config.UseRedis {
		return nil, errors.New("no redis client, set USE_REDIS in .env")
	}
	note := &models.Note{}
	noteCachekey := getNoteCacheKey(userId, noteId)
	err := GetRedisCache().Get(context.TODO(), noteCachekey, note)
	return note, err
}
