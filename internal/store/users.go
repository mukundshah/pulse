package store

import (
	"pulse/internal/models"

	"github.com/google/uuid"
)

func (s *Store) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) CreateUser(user *models.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	return s.db.Create(user).Error
}

func (s *Store) UpdateUser(user *models.User) error {
	return s.db.Save(user).Error
}

