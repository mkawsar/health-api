package services

import (
	"errors"
	db "health/models/db"
	models "health/models"
)

// CreateNote creates a new note with the given title and content belonging to the user with the given userId.
// The note is created with a unique ID and the current time as the createdAt and updatedAt timestamps.
// If the note cannot be created, an error is returned.
func CreateNote(userId uint, title string, content string) (*db.Note, error) {
	note := db.NewNote(userId, title, content)
	err := DB.Create(note).Error
	if err != nil {
		return nil, errors.New("cannot create new note")
	}
	return note, nil
}

// GetNotes retrieves all notes belonging to a user with the given userId,
// paginated to the given page and limit.
// The result is a slice of Note documents.
// If the page is out of bounds, an error is returned.
func GetNotes(userId uint, page int, limit int) ([]db.Note, error) {
	var notes []db.Note
	offset := page * limit
	err := DB.Where("author_id = ?", userId).
		Offset(offset).
		Limit(limit + 1).
		Find(&notes).Error

	if err != nil {
		return nil, errors.New("cannot find notes")
	}
	return notes, nil
}

// GetNoteById retrieves a note from the PostgreSQL database by the given noteId, only if the user with the given userId is the author.
// If the note does not exist, an error is returned.
func GetNoteById(userId uint, noteId uint) (*db.Note, error) {
	note := &db.Note{}
	err := DB.Where("id = ? AND author_id = ?", noteId, userId).First(note).Error
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
func UpdateNote(userId uint, noteId uint, request *models.NoteRequest) error {
	note := &db.Note{}
	err := DB.Where("id = ? AND author_id = ?", noteId, userId).First(note).Error
	if err != nil {
		return errors.New("cannot find note")
	}

	note.Title = request.Title
	note.Content = request.Content
	err = DB.Save(note).Error

	if err != nil {
		return errors.New("cannot update")
	}

	return nil
}

// DeleteNote deletes a note with the given noteId if the user with the given userId is the author.
// If the note does not exist or the deletion fails, an error is returned.
func DeleteNote(userId uint, noteId uint) error {
	result := DB.Where("id = ? AND author_id = ?", noteId, userId).Delete(&db.Note{})
	if result.Error != nil || result.RowsAffected == 0 {
		return errors.New("cannot delete note")
	}

	return nil
}
