package store

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type todoStore struct {
	db *gorm.DB
}

func NewTodoStore(db *gorm.DB) TodoStore {
	return &todoStore{
		db: db,
	}
}

func (s *todoStore) GetAll(username string) ([]TodoItem, error) {
	var items []TodoItem
	err := s.db.Find(&items, "username = ?", username).Error

	if err != nil {
		return []TodoItem{}, err
	}

	return items, nil
}

func (s *todoStore) GetById(username, id string) (TodoItem, error) {
	var item TodoItem
	err := s.db.First(&item, "username = ? AND id = ?", username, id).Error
	if err != nil {
		return TodoItem{}, err
	}
	return item, nil
}

func (s *todoStore) Update(username, id, task, dueDate string) (TodoItem, error) {
	err := s.db.Model(&TodoItem{Id: id, Username: username}).Updates(TodoItem{Task: task, DueDate: dueDate}).Error
	if err != nil {
		return TodoItem{}, err
	}
	return s.GetById(username, id)
}

func (s *todoStore) Delete(username, id string) error {
	err := s.db.Delete(&TodoItem{}, "username = ? AND id = ?", username, id).Error
	return err
}

func (s *todoStore) Create(username, task, dueDate string) (TodoItem, error) {
	id := uuid.New().String()

	item := TodoItem{
		Username: username,
		Id:       id,
		Task:     task,
		DueDate:  dueDate,
		Status:   "not-started",
	}

	err := s.db.Create(&item).Error
	if err != nil {
		return TodoItem{}, err
	}

	return item, nil
}
