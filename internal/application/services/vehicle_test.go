package services

import (
	"context"
	"errors"
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/vehicle"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVehicleService_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockVehicleRepository)
	service := NewVehicleService(mockRepo)

	ctx := context.Background()
	plate, _ := domain.NewPlate("ABC1234")
	v := domain.Vehicle{
		Plate: plate,
		Name:  "Honda Civic",
		Year:  2020,
		Brand: "Honda",
		Color: "Branco",
	}

	mockRepo.On("GetByPlate", ctx, plate.GetPlate()).Return(nil, nil)

	var mockModel vehicle.Model
	mockModel.FromDomain(&v)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*vehicle.Model")).Return(&mockModel, nil)

	result, err := service.Create(ctx, v)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, v.Name, result.Name)
	assert.Equal(t, v.Brand, result.Brand)
	assert.Equal(t, v.Color, result.Color)
	mockRepo.AssertExpectations(t)
}

func TestVehicleService_Create_PlateAlreadyExists(t *testing.T) {
	mockRepo := new(mocks.MockVehicleRepository)
	service := NewVehicleService(mockRepo)

	ctx := context.Background()
	plate, _ := domain.NewPlate("ABC1234")
	v := domain.Vehicle{
		Plate: plate,
		Name:  "Honda Civic",
		Year:  2020,
		Brand: "Honda",
		Color: "Branco",
	}

	var existingModel vehicle.Model
	existingModel.FromDomain(&v)
	mockRepo.On("GetByPlate", ctx, plate.GetPlate()).Return(&existingModel, nil)

	result, err := service.Create(ctx, v)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "Vehicle already exists")
	mockRepo.AssertExpectations(t)
}

func TestVehicleService_Create_GetByPlateError(t *testing.T) {
	mockRepo := new(mocks.MockVehicleRepository)
	service := NewVehicleService(mockRepo)

	ctx := context.Background()
	plate, _ := domain.NewPlate("ABC1234")
	v := domain.Vehicle{
		Plate: plate,
		Name:  "Honda Civic",
		Year:  2020,
		Brand: "Honda",
		Color: "Branco",
	}

	expectedError := errors.New("database error")
	mockRepo.On("GetByPlate", ctx, plate.GetPlate()).Return(nil, expectedError)

	result, err := service.Create(ctx, v)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestVehicleService_Create_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.MockVehicleRepository)
	service := NewVehicleService(mockRepo)

	ctx := context.Background()
	plate, _ := domain.NewPlate("ABC1234")
	v := domain.Vehicle{
		Plate: plate,
		Name:  "Honda Civic",
		Year:  2020,
		Brand: "Honda",
		Color: "Branco",
	}

	expectedError := errors.New("database connection error")
	mockRepo.On("GetByPlate", ctx, plate.GetPlate()).Return(nil, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*vehicle.Model")).Return(nil, expectedError)

	result, err := service.Create(ctx, v)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestVehicleService_GetByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockVehicleRepository)
	service := NewVehicleService(mockRepo)

	ctx := context.Background()
	id := uint(1)

	plate, _ := domain.NewPlate("ABC1234")
	v := domain.Vehicle{
		ID:    id,
		Plate: plate,
		Name:  "Honda Civic",
		Year:  2020,
		Brand: "Honda",
		Color: "Branco",
	}

	var repoResp vehicle.Model
	repoResp.FromDomain(&v)

	mockRepo.On("FindByID", ctx, id).Return(&repoResp, nil)

	result, err := service.GetByID(ctx, id)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, v.Name, result.Name)
	assert.Equal(t, v.Brand, result.Brand)
	mockRepo.AssertExpectations(t)
}

func TestVehicleService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockVehicleRepository)
	service := NewVehicleService(mockRepo)

	ctx := context.Background()
	id := uint(999)

	mockRepo.On("FindByID", ctx, id).Return(nil, errors.New("vehicle not found"))

	result, err := service.GetByID(ctx, id)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestVehicleService_GetByPlate_Success(t *testing.T) {
	mockRepo := new(mocks.MockVehicleRepository)
	service := NewVehicleService(mockRepo)

	ctx := context.Background()
	plateStr := "ABC1234"
	plate, _ := domain.NewPlate(plateStr)

	v := domain.Vehicle{
		ID:    1,
		Plate: plate,
		Name:  "Honda Civic",
		Year:  2020,
		Brand: "Honda",
		Color: "Branco",
	}

	var repoResp vehicle.Model
	repoResp.FromDomain(&v)

	mockRepo.On("GetByPlate", ctx, plateStr).Return(&repoResp, nil)

	result, err := service.GetByPlate(ctx, plateStr)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, v.Name, result.Name)
	mockRepo.AssertExpectations(t)
}

func TestVehicleService_GetByPlate_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockVehicleRepository)
	service := NewVehicleService(mockRepo)

	ctx := context.Background()
	plateStr := "XYZ9999"

	mockRepo.On("GetByPlate", ctx, plateStr).Return(nil, nil)

	result, err := service.GetByPlate(ctx, plateStr)

	assert.NoError(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestVehicleService_GetByPlate_Error(t *testing.T) {
	mockRepo := new(mocks.MockVehicleRepository)
	service := NewVehicleService(mockRepo)

	ctx := context.Background()
	plateStr := "ABC1234"

	expectedError := errors.New("database error")
	mockRepo.On("GetByPlate", ctx, plateStr).Return(nil, expectedError)

	result, err := service.GetByPlate(ctx, plateStr)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestVehicleService_GetAll_Success(t *testing.T) {
	mockRepo := new(mocks.MockVehicleRepository)
	service := NewVehicleService(mockRepo)

	ctx := context.Background()

	plate1, _ := domain.NewPlate("ABC1234")
	plate2, _ := domain.NewPlate("XYZ5678")

	expectedVehicles := []vehicle.Model{
		{Plate: plate1.GetPlate(), Name: "Honda Civic", Year: 2020, Brand: "Honda", Color: "Branco"},
		{Plate: plate2.GetPlate(), Name: "Toyota Corolla", Year: 2021, Brand: "Toyota", Color: "Preto"},
	}

	vSearch := ports.VehicleSearch{Name: "", Plate: "", Brand: "", Color: "", Year: 0, Status: true}
	mockRepo.On("Search", ctx, vSearch).Return(expectedVehicles)

	result, err := service.GetAll(ctx, map[string]interface{}{})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 2)
	assert.Equal(t, "Honda Civic", (*result)[0].Name)
	assert.Equal(t, "Toyota Corolla", (*result)[1].Name)
	mockRepo.AssertExpectations(t)
}

func TestVehicleService_GetAll_WithFilters(t *testing.T) {
	mockRepo := new(mocks.MockVehicleRepository)
	service := NewVehicleService(mockRepo)

	ctx := context.Background()
	params := map[string]interface{}{
		"name":   "Civic",
		"plate":  "ABC",
		"brand":  "Honda",
		"color":  "Branco",
		"year":   "2020",
		"status": "true",
	}

	plate, _ := domain.NewPlate("ABC1234")
	expectedVehicles := []vehicle.Model{
		{Plate: plate.GetPlate(), Name: "Honda Civic", Year: 2020, Brand: "Honda", Color: "Branco"},
	}

	vSearch := ports.VehicleSearch{
		Name:   "Civic",
		Plate:  "ABC",
		Brand:  "Honda",
		Color:  "Branco",
		Year:   2020,
		Status: true,
	}
	mockRepo.On("Search", ctx, vSearch).Return(expectedVehicles)

	result, err := service.GetAll(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 1)
	assert.Equal(t, "Honda Civic", (*result)[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestVehicleService_GetAll_WithPartialFilters(t *testing.T) {
	mockRepo := new(mocks.MockVehicleRepository)
	service := NewVehicleService(mockRepo)

	ctx := context.Background()
	params := map[string]interface{}{
		"name": "Civic",
	}

	plate, _ := domain.NewPlate("ABC1234")
	expectedVehicles := []vehicle.Model{
		{Plate: plate.GetPlate(), Name: "Honda Civic", Year: 2020, Brand: "Honda", Color: "Branco"},
	}

	vSearch := ports.VehicleSearch{
		Name:   "Civic",
		Plate:  "",
		Brand:  "",
		Color:  "",
		Year:   0,
		Status: true,
	}
	mockRepo.On("Search", ctx, vSearch).Return(expectedVehicles)

	result, err := service.GetAll(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 1)
	mockRepo.AssertExpectations(t)
}

func TestVehicleService_Update_Success(t *testing.T) {
	mockRepo := new(mocks.MockVehicleRepository)
	service := NewVehicleService(mockRepo)

	ctx := context.Background()
	plate, _ := domain.NewPlate("ABC1234")
	existingVehicle := domain.Vehicle{
		ID:    1,
		Plate: plate,
		Name:  "Honda Civic",
		Year:  2020,
		Brand: "Honda",
		Color: "Branco",
	}

	var existingModel vehicle.Model
	existingModel.FromDomain(&existingVehicle)

	updatedVehicle := domain.Vehicle{
		ID:    1,
		Plate: plate,
		Name:  "Honda Civic LX",
		Year:  2020,
		Brand: "Honda",
		Color: "Branco",
	}

	mockRepo.On("FindByID", ctx, uint(1)).Return(&existingModel, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*vehicle.Model")).Return(nil)

	result, err := service.Update(ctx, updatedVehicle)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Honda Civic LX", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestVehicleService_Update_VehicleNotFound(t *testing.T) {
	mockRepo := new(mocks.MockVehicleRepository)
	service := NewVehicleService(mockRepo)

	ctx := context.Background()
	plate, _ := domain.NewPlate("ABC1234")
	updatedVehicle := domain.Vehicle{
		ID:    999,
		Plate: plate,
		Name:  "Honda Civic",
	}

	mockRepo.On("FindByID", ctx, uint(999)).Return(nil, errors.New("vehicle not found"))

	result, err := service.Update(ctx, updatedVehicle)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "vehicle not found")
	mockRepo.AssertExpectations(t)
}

func TestVehicleService_Update_PlateAlreadyExists(t *testing.T) {
	mockRepo := new(mocks.MockVehicleRepository)
	service := NewVehicleService(mockRepo)

	ctx := context.Background()
	existingPlate, _ := domain.NewPlate("ABC1234")
	newPlate, _ := domain.NewPlate("XYZ5678")

	existingVehicle := domain.Vehicle{
		ID:    1,
		Plate: existingPlate,
		Name:  "Honda Civic",
	}

	anotherVehicle := domain.Vehicle{
		ID:    2,
		Plate: newPlate,
		Name:  "Toyota Corolla",
	}

	updatedVehicle := domain.Vehicle{
		ID:    1,
		Plate: newPlate,
		Name:  "Honda Civic",
	}

	var existingModel vehicle.Model
	existingModel.FromDomain(&existingVehicle)

	var anotherModel vehicle.Model
	anotherModel.FromDomain(&anotherVehicle)

	mockRepo.On("FindByID", ctx, uint(1)).Return(&existingModel, nil)
	mockRepo.On("GetByPlate", ctx, newPlate.GetPlate()).Return(&anotherModel, nil)

	result, err := service.Update(ctx, updatedVehicle)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "Vehicle already exists")
	mockRepo.AssertExpectations(t)
}

func TestVehicleService_Update_PlateUnchanged(t *testing.T) {
	mockRepo := new(mocks.MockVehicleRepository)
	service := NewVehicleService(mockRepo)

	ctx := context.Background()
	plate, _ := domain.NewPlate("ABC1234")
	existingVehicle := domain.Vehicle{
		ID:    1,
		Plate: plate,
		Name:  "Honda Civic",
		Year:  2020,
		Brand: "Honda",
		Color: "Branco",
	}

	var existingModel vehicle.Model
	existingModel.FromDomain(&existingVehicle)

	updatedVehicle := domain.Vehicle{
		ID:    1,
		Plate: plate,
		Name:  "Honda Civic LX",
		Year:  2020,
		Brand: "Honda",
		Color: "Branco",
	}

	mockRepo.On("FindByID", ctx, uint(1)).Return(&existingModel, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*vehicle.Model")).Return(nil)

	result, err := service.Update(ctx, updatedVehicle)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestVehicleService_Update_WithoutChangingPlate(t *testing.T) {
	mockRepo := new(mocks.MockVehicleRepository)
	service := NewVehicleService(mockRepo)

	ctx := context.Background()
	plate, _ := domain.NewPlate("ABC1234")
	existingVehicle := domain.Vehicle{
		ID:    1,
		Plate: plate,
		Name:  "Honda Civic",
		Year:  2020,
		Brand: "Honda",
		Color: "Branco",
	}

	var existingModel vehicle.Model
	existingModel.FromDomain(&existingVehicle)

	updatedVehicle := domain.Vehicle{
		ID:    1,
		Plate: plate,
		Name:  "Honda Civic LX",
		Year:  2020,
		Brand: "Honda",
		Color: "Branco",
	}

	mockRepo.On("FindByID", ctx, uint(1)).Return(&existingModel, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*vehicle.Model")).Return(nil)

	result, err := service.Update(ctx, updatedVehicle)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Honda Civic LX", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestVehicleService_Update_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.MockVehicleRepository)
	service := NewVehicleService(mockRepo)

	ctx := context.Background()
	plate, _ := domain.NewPlate("ABC1234")
	existingVehicle := domain.Vehicle{
		ID:    1,
		Plate: plate,
		Name:  "Honda Civic",
	}

	var existingModel vehicle.Model
	existingModel.FromDomain(&existingVehicle)

	updatedVehicle := domain.Vehicle{
		ID:    1,
		Plate: plate,
		Name:  "Honda Civic LX",
	}

	expectedError := errors.New("database error")
	mockRepo.On("FindByID", ctx, uint(1)).Return(&existingModel, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*vehicle.Model")).Return(expectedError)

	result, err := service.Update(ctx, updatedVehicle)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestVehicleService_Delete_Success(t *testing.T) {
	mockRepo := new(mocks.MockVehicleRepository)
	service := NewVehicleService(mockRepo)

	ctx := context.Background()
	id := uint(1)

	mockRepo.On("Delete", ctx, id).Return(nil)

	err := service.Delete(ctx, id)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVehicleService_Delete_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.MockVehicleRepository)
	service := NewVehicleService(mockRepo)

	ctx := context.Background()
	id := uint(1)

	expectedError := errors.New("database error")
	mockRepo.On("Delete", ctx, id).Return(expectedError)

	err := service.Delete(ctx, id)

	assert.Error(t, err)
	assert.EqualError(t, err, "Vehicle not found")
	mockRepo.AssertExpectations(t)
}

func TestVehicleService_GetAll_WithYearAndStatusConversion(t *testing.T) {
	mockRepo := new(mocks.MockVehicleRepository)
	service := NewVehicleService(mockRepo)

	ctx := context.Background()
	params := map[string]interface{}{
		"year":   "2020",
		"status": "false",
	}

	vSearch := ports.VehicleSearch{
		Name:   "",
		Plate:  "",
		Brand:  "",
		Color:  "",
		Year:   2020,
		Status: false,
	}
	mockRepo.On("Search", ctx, vSearch).Return([]vehicle.Model{})

	result, err := service.GetAll(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestVehicleService_GetAll_WithInvalidYear(t *testing.T) {
	mockRepo := new(mocks.MockVehicleRepository)
	service := NewVehicleService(mockRepo)

	ctx := context.Background()
	params := map[string]interface{}{
		"year": "invalid",
	}

	vSearch := ports.VehicleSearch{
		Name:   "",
		Plate:  "",
		Brand:  "",
		Color:  "",
		Year:   0,
		Status: true,
	}
	mockRepo.On("Search", ctx, vSearch).Return([]vehicle.Model{})

	result, err := service.GetAll(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestVehicleService_GetAll_WithInvalidStatus(t *testing.T) {
	mockRepo := new(mocks.MockVehicleRepository)
	service := NewVehicleService(mockRepo)

	ctx := context.Background()
	params := map[string]interface{}{
		"status": "invalid",
	}

	vSearch := ports.VehicleSearch{
		Name:   "",
		Plate:  "",
		Brand:  "",
		Color:  "",
		Year:   0,
		Status: false,
	}
	mockRepo.On("Search", ctx, vSearch).Return([]vehicle.Model{})

	result, err := service.GetAll(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}
