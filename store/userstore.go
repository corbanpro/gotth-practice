package store

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) UserStore {
	return &userStore{
		db: db,
	}
}

func (s *userStore) GetByUsername(username string) (User, error) {
	var user User
	if err := s.db.First(&user, "username = ?", username).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *userStore) Create(username, password, firstName, lastName string) (User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return User{}, err
	}

	user := User{
		Username:  username,
		Password:  hashedPassword,
		FirstName: firstName,
		LastName:  lastName,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *userStore) Update(username, password, firstName, lastName string) (User, error) {
	err := s.db.Model(&User{Username: username}).Updates(User{
		Username:  username,
		Password:  password,
		FirstName: firstName, LastName: lastName,
	}).Error

	if err != nil {
		return User{}, err
	}

	return s.GetByUsername(username)
}

func (s *userStore) Delete(username string) error {
	err := s.db.Delete(&User{Username: username}).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *userStore) ValidatePassword(username, password string) (bool, error) {
	user, err := s.GetByUsername(username)

	if err != nil {
		return false, err
	}

	if !comparePasswords(user.Password, password) {
		return false, nil
	}

	return true, nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func comparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
