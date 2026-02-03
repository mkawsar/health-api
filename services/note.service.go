package services

import (
	"errors"
	models "health/models"
	db "health/models/db"

	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateNote creates a new note with the given title and content belonging to the user with the given userId.
// The note is created with a unique ID and the current time as the createdAt and updatedAt timestamps.
// If the note cannot be created, an error is returned.
func CreateNote(userId primitive.ObjectID, title string, content string) (*db.Note, error) {
	note := db.NewNote(userId, title, content)
	err := mgm.Coll(note).Create(note)
	if err != nil {
		return nil, errors.New("cannot create new note")
	}
	return note, nil
}

// GetNotes retrieves all notes belonging to a user with the given userId,
// paginated to the given page and limit.
// The result is a slice of Note documents.
// If the page is out of bounds, an error is returned.
func GetNotes(userId primitive.ObjectID, page int, limit int) ([]db.Note, error) {
	var notes []db.Note
	findOptions := options.Find().SetSkip(int64(page * limit)).SetLimit(int64(limit + 1))
	err := mgm.Coll(&db.Note{}).SimpleFind(
		&notes,
		bson.M{"author": userId.Hex()},
		findOptions,
	)

	if err != nil {
		return nil, errors.New("cannot find notes")
	}
	return notes, nil
}

// GetNoteById retrieves a note from the MongoDB database by the given noteId, only if the user with the given userId is the author.
// If the note does not exist, an error is returned.
func GetNoteById(userId primitive.ObjectID, noteId primitive.ObjectID) (*db.Note, error) {
	note := &db.Note{}
	err := mgm.Coll(note).First(bson.M{field.ID: noteId, "author": userId.Hex()}, note)
	if err != nil {
		return nil, errors.New("cannot find note")
	}

	return note, nil
}

// UpdateNote updates a note with the given noteId, only if the user with the given userId is the author.
// The note is updated with the given title and content.
// If the note does not exist, an error is returned.
// If the user is not the author, an error is returned.
// If the note cannot be updated, an error is returned.
func UpdateNote(userId primitive.ObjectID, noteId primitive.ObjectID, request *models.NoteRequest) error {
	note := &db.Note{}
	err := mgm.Coll(note).FindByID(noteId, note)
	if err != nil {
		return errors.New("cannot find note")
	}

	if note.Author != userId.Hex() {
		return errors.New("you cannot update this note")
	}

	note.Title = request.Title
	note.Content = request.Content
	err = mgm.Coll(note).Update(note)

	if err != nil {
		return errors.New("cannot update")
	}

	return nil
}

// DeleteNote deletes a note with the given noteId if the user with the given userId is the author.
// If the note does not exist or the deletion fails, an error is returned.
func DeleteNote(userId primitive.ObjectID, noteId primitive.ObjectID) error {
	deleteResult, err := mgm.Coll(&db.Note{}).DeleteOne(mgm.Ctx(), bson.M{field.ID: noteId, "author": userId.Hex()})
	if err != nil || deleteResult.DeletedCount <= 0 {
		return errors.New("cannot delete note")
	}

	return nil
}
