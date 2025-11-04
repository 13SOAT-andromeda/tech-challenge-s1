package services

import (
	"context"
	"errors"
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/maintenance"
	omodel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order_maintenance"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMaintenanceService_Create(t *testing.T) {
	ctx := context.Background()
	mrepo := &mocks.MockMaintenanceRepository{}
	omrepo := &mocks.MockGenericRepository[omodel.Model]{}
	svc := NewMaintenanceService(mrepo, omrepo)

	input := domain.Maintenance{Name: "Oil change", Price: 1000, CategoryId: domain.MaintenanceCategory("standard")}
	var inModel maintenance.Model
	inModel.FromDomain(&input)

	returned := inModel
	returned.ID = 10

	mrepo.On("Create", ctx, &inModel).Return(&returned, nil)

	res, err := svc.Create(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, uint(10), res.ID)
	mrepo.AssertExpectations(t)
}

func TestMaintenanceService_GetByID(t *testing.T) {
	ctx := context.Background()
	mrepo := &mocks.MockMaintenanceRepository{}
	omrepo := &mocks.MockGenericRepository[omodel.Model]{}
	svc := NewMaintenanceService(mrepo, omrepo)

	model := maintenance.Model{}
	model.ID = 5
	model.Name = "Brake"
	model.Price = 2000
	model.CategoryId = "standard"

	mrepo.On("FindByID", mock.Anything, uint(5)).Return(&model, nil)

	res, err := svc.GetByID(ctx, 5)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, uint(5), res.ID)
	assert.NotEmpty(t, res.CategoryId)
	mrepo.AssertExpectations(t)
}

func TestMaintenanceService_GetByIDs_Empty(t *testing.T) {
	ctx := context.Background()
	mrepo := &mocks.MockMaintenanceRepository{}
	omrepo := &mocks.MockGenericRepository[omodel.Model]{}
	svc := NewMaintenanceService(mrepo, omrepo)

	res, err := svc.GetByIDs(ctx, []uint{})
	assert.NoError(t, err)
	assert.Empty(t, res)
}

func TestMaintenanceService_GetByIDs(t *testing.T) {
	ctx := context.Background()
	mrepo := &mocks.MockMaintenanceRepository{}
	omrepo := &mocks.MockGenericRepository[omodel.Model]{}
	svc := NewMaintenanceService(mrepo, omrepo)

	m1 := maintenance.Model{Name: "A", Price: 150, CategoryId: "standard"}
	m1.ID = 1
	m2 := maintenance.Model{Name: "B", Price: 200, CategoryId: "standard"}
	m2.ID = 2
	mrepo.On("FindByIDs", ctx, []uint{1, 2}).Return([]maintenance.Model{m1, m2}, nil)

	res, err := svc.GetByIDs(ctx, []uint{1, 2})
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, uint(1), res[0].ID)
	assert.Equal(t, uint(2), res[1].ID)
	mrepo.AssertExpectations(t)
}

func TestMaintenanceService_UpdateByID(t *testing.T) {
	ctx := context.Background()
	mrepo := &mocks.MockMaintenanceRepository{}
	omrepo := &mocks.MockGenericRepository[omodel.Model]{}
	svc := NewMaintenanceService(mrepo, omrepo)

	input := domain.Maintenance{ID: 7, Name: "X", Price: 500, CategoryId: domain.MaintenanceCategory("standard")}
	var model maintenance.Model
	model.FromDomain(&input)

	mrepo.On("Update", ctx, &model).Return(nil)

	err := svc.UpdateByID(ctx, 7, input)
	assert.NoError(t, err)
	mrepo.AssertExpectations(t)
}

func TestMaintenanceService_UpdateByID_Error(t *testing.T) {
	ctx := context.Background()
	mrepo := &mocks.MockMaintenanceRepository{}
	omrepo := &mocks.MockGenericRepository[omodel.Model]{}
	svc := NewMaintenanceService(mrepo, omrepo)

	input := domain.Maintenance{ID: 7, Name: "X"}
	var model maintenance.Model
	model.FromDomain(&input)

	errExpected := errors.New("update failed")
	mrepo.On("Update", ctx, &model).Return(errExpected)

	err := svc.UpdateByID(ctx, 7, input)
	assert.Error(t, err)
	assert.Equal(t, errExpected, err)
	mrepo.AssertExpectations(t)
}

func TestMaintenanceService_DeleteByID(t *testing.T) {
	ctx := context.Background()
	mrepo := &mocks.MockMaintenanceRepository{}
	omrepo := &mocks.MockGenericRepository[omodel.Model]{}
	svc := NewMaintenanceService(mrepo, omrepo)

	model := maintenance.Model{Name: "Del", Price: 300, CategoryId: "standard"}
	model.ID = 9

	mrepo.On("FindByID", ctx, uint(9)).Return(&model, nil)
	mrepo.On("Delete", ctx, uint(9)).Return(nil)

	res, err := svc.DeleteByID(ctx, 9)
	assert.NoError(t, err)
	assert.Equal(t, uint(9), res.ID)
	mrepo.AssertExpectations(t)
}

func TestMaintenanceService_DeleteByID_FindError(t *testing.T) {
	ctx := context.Background()
	mrepo := &mocks.MockMaintenanceRepository{}
	omrepo := &mocks.MockGenericRepository[omodel.Model]{}
	svc := NewMaintenanceService(mrepo, omrepo)

	errExpected := errors.New("not found")
	mrepo.On("FindByID", ctx, uint(99)).Return((*maintenance.Model)(nil), errExpected)

	res, err := svc.DeleteByID(ctx, 99)
	assert.Error(t, err)
	assert.Nil(t, res)
	mrepo.AssertExpectations(t)
}

func TestMaintenanceService_CreateOrderMaintenances_Success(t *testing.T) {
	ctx := context.Background()
	mrepo := &mocks.MockMaintenanceRepository{}
	omrepo := &mocks.MockGenericRepository[omodel.Model]{}
	svc := NewMaintenanceService(mrepo, omrepo)

	orderId := uint(3)
	ids := []uint{11, 12}

	omrepo.On("Create", ctx, mock.MatchedBy(func(m *omodel.Model) bool {
		return m != nil && m.MaintenanceId == 11 && m.OrderId == orderId
	})).Return(&omodel.Model{}, nil)
	omrepo.On("Create", ctx, mock.MatchedBy(func(m *omodel.Model) bool {
		return m != nil && m.MaintenanceId == 12 && m.OrderId == orderId
	})).Return(&omodel.Model{}, nil)

	err := svc.CreateOrderMaintenances(ctx, orderId, ids)
	assert.NoError(t, err)
	omrepo.AssertExpectations(t)
}

func TestMaintenanceService_CreateOrderMaintenances_Error(t *testing.T) {
	ctx := context.Background()
	mrepo := &mocks.MockMaintenanceRepository{}
	omrepo := &mocks.MockGenericRepository[omodel.Model]{}
	svc := NewMaintenanceService(mrepo, omrepo)

	orderId := uint(3)
	ids := []uint{21, 22}

	errExpected := errors.New("create fail")
	omrepo.On("Create", ctx, mock.MatchedBy(func(m *omodel.Model) bool {
		return m != nil && m.MaintenanceId == 21 && m.OrderId == orderId
	})).Return((*omodel.Model)(nil), errExpected)

	err := svc.CreateOrderMaintenances(ctx, orderId, ids)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create order maintenance")
	omrepo.AssertExpectations(t)
}
