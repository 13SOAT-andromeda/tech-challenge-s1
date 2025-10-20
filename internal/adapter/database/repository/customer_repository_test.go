package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/address"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	require.NoError(t, err)

	return gormDB, mock, sqlDB
}

func TestCustomerRepository_FindByEmail(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewCustomerRepository(db)
	ctx := context.Background()

	t.Run("success - customer found", func(t *testing.T) {
		expectedCustomer := customer.Model{
			Email:    "test@example.com",
			Name:     "Test User",
			Document: "12345678900",
			Type:     "CPF",
			Contact:  "11999999999",
			Address: address.Model{
				City:          "New York",
				Country:       "US",
				Address:       "Teste",
				Neighborhood:  "New York",
				ZipCode:       "12345",
				AddressNumber: "12345",
			},
		}

		rows := sqlmock.NewRows([]string{"id", "name", "email", "document", "type", "contact", "address", "address_number", "zip_code", "city", "country", "neighborhood"}).
			AddRow(
				expectedCustomer.ID,
				expectedCustomer.Name,
				expectedCustomer.Email,
				expectedCustomer.Document,
				expectedCustomer.Type,
				expectedCustomer.Contact,
				expectedCustomer.Address.Address,
				expectedCustomer.Address.AddressNumber,
				expectedCustomer.Address.ZipCode,
				expectedCustomer.Address.City,
				expectedCustomer.Address.Country,
				expectedCustomer.Address.Neighborhood)

		expectedSQL := `SELECT * FROM "Customer" WHERE email = $1 AND "Customer"."deleted_at" IS NULL ORDER BY "Customer"."id" LIMIT $2`

		mock.ExpectQuery(regexp.QuoteMeta(
			expectedSQL)).
			WithArgs(expectedCustomer.Email, 1).
			WillReturnRows(rows)

		customer, err := repo.FindByEmail(ctx, expectedCustomer.Email)

		assert.NoError(t, err)
		assert.NotNil(t, customer)
		assert.Equal(t, expectedCustomer.Email, customer.Email)
		assert.Equal(t, expectedCustomer.Name, customer.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - customer not found", func(t *testing.T) {
		expectedSQL := regexp.QuoteMeta(`SELECT * FROM "Customer" WHERE email = $1 AND "Customer"."deleted_at" IS NULL ORDER BY "Customer"."id" LIMIT $2`)
		emailNaoEncontrado := "notfound@example.com"

		mock.ExpectQuery(expectedSQL).
			WithArgs(emailNaoEncontrado, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		customer, err := repo.FindByEmail(ctx, emailNaoEncontrado)

		assert.Error(t, err)
		assert.Nil(t, customer)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBaseRepository_FindByID(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewCustomerRepository(db)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedCustomer := domain.Customer{
			ID:       1,
			Email:    "test@example.com",
			Name:     "Test User",
			Document: "12345678900",
			Type:     "CPF",
			Contact:  "11999999999",
		}

		rows := sqlmock.NewRows([]string{"id", "name", "email", "document", "type", "contact"}).
			AddRow(expectedCustomer.ID, expectedCustomer.Name, expectedCustomer.Email,
				expectedCustomer.Document, expectedCustomer.Type, expectedCustomer.Contact)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "Customer" WHERE "Customer"."id" = $1 AND "Customer"."deleted_at" IS NULL ORDER BY "Customer"."id" LIMIT $2`)).
			WithArgs(expectedCustomer.ID, 1).
			WillReturnRows(rows)

		customer, err := repo.FindByID(ctx, expectedCustomer.ID)

		assert.NoError(t, err)
		assert.NotNil(t, customer)
		assert.Equal(t, expectedCustomer.ID, customer.ID)
		assert.Equal(t, expectedCustomer.Email, customer.Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - customer not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "Customer" WHERE "Customer"."id" = $1 AND "Customer"."deleted_at" IS NULL ORDER BY "Customer"."id" LIMIT $2`)).
			WithArgs(uint(999), 1).
			WillReturnError(gorm.ErrRecordNotFound)

		customer, err := repo.FindByID(ctx, uint(999))

		assert.Error(t, err)
		assert.Nil(t, customer)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBaseRepository_Create(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewCustomerRepository(db)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		customer := &customer.Model{
			Email:    "new@example.com",
			Name:     "New User",
			Document: "12345678900",
			Type:     "CPF",
			Contact:  "11999999999",
		}

		expectedSQL := `INSERT INTO "Customer" ("created_at","updated_at","deleted_at","name","email","document","type","contact","address","address_number","neighborhood","city","country","zip_code") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING "id"`

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				customer.Name,
				customer.Email,
				customer.Document,
				customer.Type,
				customer.Contact,
				customer.Address.Address,
				customer.Address.AddressNumber,
				customer.Address.Neighborhood,
				customer.Address.City,
				customer.Address.Country,
				customer.Address.ZipCode,
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		createdCustomer, err := repo.Create(ctx, customer)

		assert.NoError(t, err)
		assert.NotNil(t, createdCustomer)
		assert.Equal(t, uint(1), createdCustomer.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database error", func(t *testing.T) {
		customer := &customer.Model{
			Email:    "error@example.com",
			Name:     "Error User",
			Document: "12345678900",
			Type:     "CPF",
			Contact:  "11999999999",
		}
		expectedSQL := `INSERT INTO "Customer" ("created_at","updated_at","deleted_at","name","email","document","type","contact","address","address_number","neighborhood","city","country","zip_code") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING "id"`

		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				customer.Name,
				customer.Email,
				customer.Document,
				customer.Type,
				customer.Contact,
				customer.Address.Address,
				customer.Address.AddressNumber,
				customer.Address.Neighborhood,
				customer.Address.City,
				customer.Address.Country,
				customer.Address.ZipCode,
			).
			WillReturnError(sql.ErrConnDone)

		_, err := repo.Create(ctx, customer)

		assert.Error(t, err)
		assert.Equal(t, sql.ErrConnDone, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBaseRepository_Delete(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewCustomerRepository(db)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		customerID := uint(1)

		mock.ExpectExec(regexp.QuoteMeta(
			`UPDATE "Customer" SET "deleted_at"=$1 WHERE "Customer"."id" = $2 AND "Customer"."deleted_at" IS NULL`)).
			WithArgs(sqlmock.AnyArg(), customerID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(ctx, customerID)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database error", func(t *testing.T) {
		customerID := uint(1)

		mock.ExpectExec(regexp.QuoteMeta(
			`UPDATE "Customer" SET "deleted_at"=$1 WHERE "Customer"."id" = $2 AND "Customer"."deleted_at" IS NULL`)).
			WithArgs(sqlmock.AnyArg(), customerID).
			WillReturnError(sql.ErrConnDone)

		err := repo.Delete(ctx, customerID)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBaseRepository_FindAll(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewCustomerRepository(db)
	ctx := context.Background()

	t.Run("success - multiple Customer", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "document", "type", "contact"}).
			AddRow(1, "User 1", "user1@example.com", "12345678900", "CPF", "11999999999").
			AddRow(2, "User 2", "user2@example.com", "98765432100", "CPF", "11888888888")

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "Customer"`)).
			WillReturnRows(rows)

		Customer, err := repo.FindAll(ctx)

		assert.NoError(t, err)
		assert.Len(t, Customer, 2)
		assert.Equal(t, "user1@example.com", Customer[0].Email)
		assert.Equal(t, "user2@example.com", Customer[1].Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success - empty result", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "document", "type", "contact"})

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "Customer"`)).
			WillReturnRows(rows)

		Customer, err := repo.FindAll(ctx)

		assert.NoError(t, err)
		assert.Len(t, Customer, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "Customer"`)).
			WillReturnError(sql.ErrConnDone)

		Customer, err := repo.FindAll(ctx)

		assert.Error(t, err)
		assert.Nil(t, Customer)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
