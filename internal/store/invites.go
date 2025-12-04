package store

import (
	"pulse/internal/models"

	"github.com/google/uuid"
)

func (s *Store) GetProjectInvitationByToken(token string) (*models.ProjectInvitation, error) {
	var invite models.ProjectInvitation
	err := s.db.Where("token = ? AND deleted_at IS NULL", token).First(&invite).Error
	if err != nil {
		return nil, err
	}
	return &invite, nil
}

func (s *Store) GetProjectInvitations(projectID uuid.UUID) ([]*models.ProjectInvitation, error) {
	var invites []*models.ProjectInvitation
	err := s.db.Where("project_id = ? AND deleted_at IS NULL", projectID).Find(&invites).Error
	return invites, err
}

func (s *Store) CreateProjectInvitation(invite *models.ProjectInvitation) error {
	if invite.ID == uuid.Nil {
		invite.ID = uuid.New()
	}
	return s.db.Create(invite).Error
}

func (s *Store) DeleteProjectInvitation(id uuid.UUID) error {
	return s.db.Delete(&models.ProjectInvitation{}, "id = ?", id).Error
}

