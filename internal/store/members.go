package store

import (
	"pulse/internal/models"

	"github.com/google/uuid"
)

func (s *Store) GetProjectMembers(projectID uuid.UUID) ([]*models.ProjectMember, error) {
	var members []*models.ProjectMember
	err := s.db.Where("project_id = ? AND deleted_at IS NULL", projectID).
		Preload("User").
		Find(&members).Error
	return members, err
}

func (s *Store) GetProjectMember(projectID, userID uuid.UUID) (*models.ProjectMember, error) {
	var member models.ProjectMember
	err := s.db.Where("project_id = ? AND user_id = ? AND deleted_at IS NULL", projectID, userID).
		First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (s *Store) CreateProjectMember(member *models.ProjectMember) error {
	if member.ID == uuid.Nil {
		member.ID = uuid.New()
	}
	return s.db.Create(member).Error
}

func (s *Store) UpdateProjectMemberRole(projectID, userID uuid.UUID, role string) error {
	return s.db.Model(&models.ProjectMember{}).
		Where("project_id = ? AND user_id = ? AND deleted_at IS NULL", projectID, userID).
		Update("role", role).Error
}

func (s *Store) RemoveProjectMember(projectID, userID uuid.UUID) error {
	return s.db.Where("project_id = ? AND user_id = ?", projectID, userID).
		Delete(&models.ProjectMember{}).Error
}

func (s *Store) IsProjectMember(projectID, userID uuid.UUID) (bool, error) {
	var count int64
	err := s.db.Model(&models.ProjectMember{}).
		Where("project_id = ? AND user_id = ? AND deleted_at IS NULL", projectID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *Store) IsProjectAdmin(projectID, userID uuid.UUID) (bool, error) {
	var member models.ProjectMember
	err := s.db.Where("project_id = ? AND user_id = ? AND role = ? AND deleted_at IS NULL", projectID, userID, "admin").
		First(&member).Error
	if err != nil {
		return false, nil // Not an admin if not found
	}
	return true, nil
}

// GetProjectMemberUserIDs returns a list of user IDs who are members of a project
func (s *Store) GetProjectMemberUserIDs(projectID uuid.UUID) ([]uuid.UUID, error) {
	var userIDs []uuid.UUID
	err := s.db.Model(&models.ProjectMember{}).
		Where("project_id = ? AND deleted_at IS NULL", projectID).
		Pluck("user_id", &userIDs).Error
	return userIDs, err
}
