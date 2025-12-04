package store

import (
	"pulse/internal/models"

	"github.com/google/uuid"
)

func (s *Store) CreateTag(tag *models.Tag) error {
	return s.db.Create(tag).Error
}

func (s *Store) GetTag(id uuid.UUID) (*models.Tag, error) {
	var tag models.Tag
	if err := s.db.First(&tag, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (s *Store) GetTagByName(name string, projectID uuid.UUID) (*models.Tag, error) {
	var tag models.Tag
	if err := s.db.Where("name = ? AND project_id = ?", name, projectID).First(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (s *Store) ListTags() ([]models.Tag, error) {
	var tags []models.Tag
	if err := s.db.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (s *Store) GetTagsByProject(projectID uuid.UUID) ([]models.Tag, error) {
	var tags []models.Tag
	if err := s.db.Where("project_id = ?", projectID).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (s *Store) AddTagToProject(projectID uuid.UUID, tagID uuid.UUID) error {
	// Verify project exists
	var project models.Project
	if err := s.db.First(&project, "id = ?", projectID).Error; err != nil {
		return err
	}
	// Update tag's ProjectID
	return s.db.Model(&models.Tag{}).Where("id = ?", tagID).Update("project_id", projectID).Error
}

func (s *Store) RemoveTagFromProject(projectID uuid.UUID, tagID uuid.UUID) error {
	// Verify tag belongs to the project
	var tag models.Tag
	if err := s.db.Where("id = ? AND project_id = ?", tagID, projectID).First(&tag).Error; err != nil {
		return err
	}
	// Since ProjectID is required, we delete the tag when removing from project
	return s.db.Delete(&tag).Error
}

func (s *Store) AddTagToCheck(checkID uuid.UUID, tagID uuid.UUID) error {
	var check models.Check
	if err := s.db.First(&check, "id = ?", checkID).Error; err != nil {
		return err
	}
	var tag models.Tag
	if err := s.db.First(&tag, "id = ?", tagID).Error; err != nil {
		return err
	}
	return s.db.Model(&check).Association("Tags").Append(&tag)
}

func (s *Store) RemoveTagFromCheck(checkID uuid.UUID, tagID uuid.UUID) error {
	var check models.Check
	if err := s.db.First(&check, "id = ?", checkID).Error; err != nil {
		return err
	}
	var tag models.Tag
	if err := s.db.First(&tag, "id = ?", tagID).Error; err != nil {
		return err
	}
	return s.db.Model(&check).Association("Tags").Delete(&tag)
}
