package services

import "mts_analytics/internal/domain"

type repository interface {
	Save(event domain.Event) error
	GetSignedCount() (int, error)
	GetNotSignedYetCount() (int, error)
	GetSignitionTotalTime(taskUUID string) (seconds int, err error)
}

type service struct {
	repo repository
}

func New(repo repository) *service {
	return &service{repo: repo}
}

func (s service) Save(event domain.Event) error {
	return s.repo.Save(event)
}

func (s service) GetSignedCount() (int, error) {
	return s.repo.GetSignedCount()
}

func (s service) GetNotSignedYetCount() (int, error) {
	return s.repo.GetNotSignedYetCount()
}

func (s service) GetSignitionTotalTime(taskUUID string) (int, error) {
	return s.repo.GetSignitionTotalTime(taskUUID)
}
