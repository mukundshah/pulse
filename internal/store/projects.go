package store

import (
	"pulse/internal/models"

	"github.com/google/uuid"
)

func (s *Store) CreateProject(project *models.Project) error {
	return s.db.Create(project).Error
}

func (s *Store) GetProject(id uuid.UUID) (*models.Project, error) {
	var project models.Project
	if err := s.db.First(&project, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (s *Store) ListProjects() ([]models.Project, error) {
	var projects []models.Project
	if err := s.db.Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (s *Store) UpdateProject(project *models.Project) error {
	return s.db.Save(project).Error
}

func (s *Store) DeleteProject(id uuid.UUID) error {
	return s.db.Delete(&models.Project{}, "id = ?", id).Error
}
