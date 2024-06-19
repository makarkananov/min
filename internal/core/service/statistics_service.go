package service

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"min/internal/core/domain"
	"min/internal/core/port"
)

type StatisticsService struct {
	repo port.StatisticsRepository
}

func NewStatisticsService(repo port.StatisticsRepository) *StatisticsService {
	return &StatisticsService{repo: repo}
}

func (s *StatisticsService) AddEvent(ctx context.Context, event domain.Event) error {
	log.Infof("Adding event: %v", event)

	err := s.repo.AddEvent(ctx, event)
	if err != nil {
		log.Errorf("Failed to add event: %v", err)
		return fmt.Errorf("failed to add event: %w", err)
	}

	log.Infof("Successfully added event: %v", event)
	return nil
}
