package services

import (
	"context"
	"errors"
	"strconv"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/vehicle"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/converters"
)

type vehicleService struct {
	repo ports.VehicleRepository
}

func NewVehicleService(repo ports.VehicleRepository) *vehicleService {
	return &vehicleService{repo: repo}
}

func (s *vehicleService) Create(ctx context.Context, v domain.Vehicle) (*domain.Vehicle, error) {

	if user, err := s.GetByPlate(ctx, v.Plate.GetPlate()); err != nil {
		return nil, err
	} else if user != nil {
		return nil, errors.New("Vehicle already exists")
	}

	vehicle := &vehicle.Model{}
	vehicle.FromDomain(&v)

	_, err := s.repo.Create(ctx, vehicle)

	if err != nil {
		return nil, err
	}

	created := vehicle.ToDomain()

	return created, nil
}

func (s *vehicleService) GetByID(ctx context.Context, id uint) (*domain.Vehicle, error) {

	vehicle, err := s.repo.FindByID(ctx, id)

	if err != nil {
		return nil, err
	}
	v := vehicle.ToDomain()

	return v, nil
}

func (s *vehicleService) GetByPlate(ctx context.Context, plate string) (*domain.Vehicle, error) {

	v, err := s.repo.GetByPlate(ctx, plate)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	return v.ToDomain(), nil
}

func (s *vehicleService) GetAll(ctx context.Context, params map[string]interface{}) (*[]domain.Vehicle, error) {

	vSearch := ports.VehicleSearch{Name: "", Plate: "", Brand: "", Color: "", Year: 0, Status: true}
	if params["name"] != nil {
		vSearch.Name = params["name"].(string)
	}

	if params["plate"] != nil {
		vSearch.Plate = params["plate"].(string)
	}

	if params["brand"] != nil {
		vSearch.Brand = params["brand"].(string)
	}

	if params["color"] != nil {
		vSearch.Color = params["color"].(string)
	}

	if params["year"] != nil {
		vSearch.Year, _ = strconv.Atoi(params["year"].(string))
	}

	if params["status"] != nil {
		vSearch.Status, _ = strconv.ParseBool(params["status"].(string))
	}

	vehicles := s.repo.Search(ctx, vSearch)
	vehicleD := make([]domain.Vehicle, 0, len(vehicles))

	for _, v := range vehicles {
		vehicleD = append(vehicleD, *v.ToDomain())
	}

	return &vehicleD, nil
}

func (s *vehicleService) Update(ctx context.Context, v domain.Vehicle) (*domain.Vehicle, error) {

	existing, err := s.repo.FindByID(ctx, v.ID)

	if err != nil || existing == nil {
		return nil, errors.New("vehicle not found")
	}

	existingDomain := existing.ToDomain()

	if v.Plate.GetPlate() != "" && v.Plate.GetPlate() != existingDomain.Plate.GetPlate() {
		if vehicle, err := s.GetByPlate(ctx, v.Plate.GetPlate()); err != nil {
			return nil, err
		} else if vehicle != nil && vehicle.ID != existing.ID {
			return nil, errors.New("Vehicle already exists")
		}
	}

	mergedVehicle := converters.MergeStructs(existingDomain, v).(domain.Vehicle)

	vModel := &vehicle.Model{}
	vModel.FromDomain(&mergedVehicle)

	err = s.repo.Update(ctx, vModel)
	if err != nil {
		return nil, err
	}

	updated := vModel.ToDomain()

	return updated, nil
}

func (s *vehicleService) Delete(ctx context.Context, id uint) error {

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return errors.New("Vehicle not found")
	}

	return nil
}
