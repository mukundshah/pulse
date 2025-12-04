package store

import (
	"time"

	"pulse/internal/models"
	"pulse/internal/redis"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store struct {
	db    *gorm.DB
	redis *redis.Client
}

func New(db *gorm.DB) *Store {
	return &Store{db: db}
}

func NewWithCache(db *gorm.DB, redisClient *redis.Client) *Store {
	return &Store{db: db, redis: redisClient}
}

// Project methods
func (s *Store) CreateProject(project *models.Project) error {
	return s.db.Create(project).Error
}

func (s *Store) GetProject(id uuid.UUID) (*models.Project, error) {
	var project models.Project
	if err := s.db.Preload("Tags").First(&project, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (s *Store) ListProjects() ([]models.Project, error) {
	var projects []models.Project
	if err := s.db.Preload("Tags").Find(&projects).Error; err != nil {
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

// Check methods
func (s *Store) CreateCheck(check *models.Check) error {
	return s.db.Create(check).Error
}

func (s *Store) GetCheck(id uuid.UUID) (*models.Check, error) {
	var check models.Check
	if err := s.db.Preload("Project").Preload("Tags").Preload("Regions").First(&check, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &check, nil
}

func (s *Store) GetChecksByProject(projectID uuid.UUID) ([]models.Check, error) {
	var checks []models.Check
	if err := s.db.Preload("Tags").Preload("Regions").Where("project_id = ?", projectID).Find(&checks).Error; err != nil {
		return nil, err
	}
	return checks, nil
}

func (s *Store) UpdateCheck(check *models.Check) error {
	return s.db.Save(check).Error
}

func (s *Store) DeleteCheck(id uuid.UUID) error {
	return s.db.Delete(&models.Check{}, "id = ?", id).Error
}

func (s *Store) GetDueChecks() ([]models.Check, error) {
	var checks []models.Check
	now := time.Now()
	if err := s.db.Where("is_enabled = ? AND (next_run_at IS NULL OR next_run_at <= ?)", true, now).Find(&checks).Error; err != nil {
		return nil, err
	}
	return checks, nil
}

func (s *Store) UpdateCheckStatus(checkID uuid.UUID, nextRun time.Time, consecutiveFails int, lastStatus string) error {
	now := time.Now()
	return s.db.Model(&models.Check{}).Where("id = ?", checkID).Updates(map[string]interface{}{
		"last_run_at":       now,
		"next_run_at":       nextRun,
		"consecutive_fails": consecutiveFails,
		"last_status":       lastStatus,
	}).Error
}

// CheckRun methods
func (s *Store) CreateCheckRun(run *models.CheckRun) error {
	return s.db.Create(run).Error
}

func (s *Store) GetCheckRun(id uuid.UUID) (*models.CheckRun, error) {
	var run models.CheckRun
	if err := s.db.Preload("Check").Preload("Region").First(&run, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &run, nil
}

func (s *Store) GetCheckRunsByCheck(checkID uuid.UUID, limit int) ([]models.CheckRun, error) {
	var runs []models.CheckRun
	query := s.db.Preload("Region").Where("check_id = ?", checkID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&runs).Error; err != nil {
		return nil, err
	}
	return runs, nil
}

// Tag methods
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

func (s *Store) GetTagByName(name string) (*models.Tag, error) {
	var tag models.Tag
	if err := s.db.Where("name = ?", name).First(&tag).Error; err != nil {
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

func (s *Store) AddTagToProject(projectID uuid.UUID, tagID uuid.UUID) error {
	var project models.Project
	if err := s.db.First(&project, "id = ?", projectID).Error; err != nil {
		return err
	}
	var tag models.Tag
	if err := s.db.First(&tag, "id = ?", tagID).Error; err != nil {
		return err
	}
	return s.db.Model(&project).Association("Tags").Append(&tag)
}

func (s *Store) RemoveTagFromProject(projectID uuid.UUID, tagID uuid.UUID) error {
	var project models.Project
	if err := s.db.First(&project, "id = ?", projectID).Error; err != nil {
		return err
	}
	var tag models.Tag
	if err := s.db.First(&tag, "id = ?", tagID).Error; err != nil {
		return err
	}
	return s.db.Model(&project).Association("Tags").Delete(&tag)
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
