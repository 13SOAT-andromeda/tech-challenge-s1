package services

import (
	"context"
	"fmt"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/maintenance"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order_maintenance"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type MaintenanceService struct {
	repo   ports.MaintenanceRepository
	repoOm ports.OrderMaintenanceRepository
}

func NewMaintenanceService(repo ports.MaintenanceRepository, repoOm ports.OrderMaintenanceRepository) *MaintenanceService {
	return &MaintenanceService{repo: repo, repoOm: repoOm}
}

func (s *MaintenanceService) Create(ctx context.Context, c domain.Maintenance) (*domain.Maintenance, error) {
	var model maintenance.Model
	model.FromDomain(&c)

	response, err := s.repo.Create(ctx, &model)
	if err != nil {
		return nil, err
	}

	result := response.ToDomain()
	return result, nil
}

func (s *MaintenanceService) GetByID(ctx context.Context, id uint) (*domain.Maintenance, error) {
	response, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	result := response.ToDomain()
	result.CategoryId = domain.ParseCategoryName(string(result.CategoryId))

	return result, nil
}

func (s *MaintenanceService) GetByIDs(ctx context.Context, maintenanceIDs []uint) ([]domain.Maintenance, error) {
	if len(maintenanceIDs) == 0 {
		return []domain.Maintenance{}, nil
	}

	records, err := s.repo.FindByIDs(ctx, maintenanceIDs)
	if err != nil {
		return nil, err
	}

	var result []domain.Maintenance
	for _, record := range records {
		toDomainRecord := *record.ToDomain()
		toDomainRecord.CategoryId = domain.ParseCategoryName(string(toDomainRecord.CategoryId))
		result = append(result, toDomainRecord)
	}

	return result, nil
}

func (s *MaintenanceService) GetAll(ctx context.Context) ([]domain.Maintenance, error) {
	records, err := s.repo.FindAll(ctx, false)
	if err != nil {
		return nil, err
	}

	var result []domain.Maintenance
	for _, record := range records {
		result = append(result, *record.ToDomain())
	}

	return result, nil
}

func (s *MaintenanceService) UpdateByID(ctx context.Context, id uint, c domain.Maintenance) error {
	var model maintenance.Model
	model.FromDomain(&c)

	err := s.repo.Update(ctx, &model)
	if err != nil {
		return err
	}

	return nil
}

func (s *MaintenanceService) DeleteByID(ctx context.Context, id uint) (*domain.Maintenance, error) {
	response, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return nil, err
	}

	result := response.ToDomain()
	return result, nil
}

func (s *MaintenanceService) CreateOrderMaintenances(ctx context.Context, orderId uint, maintenanceIds []uint) error {
	for _, maintenanceId := range maintenanceIds {
		var model order_maintenance.Model
		model.Maintenance.ID = maintenanceId
		model.Order.ID = orderId

		_, err := s.repoOm.Create(ctx, &model)
		if err != nil {
			return fmt.Errorf("failed to create order maintenance: %w", err)
		}
	}

	return nil
}
