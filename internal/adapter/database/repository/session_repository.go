package repository

import (
	"context"
	"errors"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) ports.SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, session *domain.Session) (*domain.Session, error) {
	sessionModel := &model.SessionModel{}
	sessionModel.FromDomain(session)

	if err := r.db.WithContext(ctx).Create(sessionModel).Error; err != nil {
		return nil, err
	}

	return sessionModel.ToDomain(), nil
}

func (r *SessionRepository) FindByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error) {
	var sessionModel model.SessionModel

	err := r.db.WithContext(ctx).
		Where("refresh_token = ?", refreshToken).
		First(&sessionModel).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("session not found")
		}
		return nil, err
	}

	return sessionModel.ToDomain(), nil
}

func (r *SessionRepository) FindByUserID(ctx context.Context, userID uint) ([]*domain.Session, error) {
	var sessionModels []model.SessionModel

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&sessionModels).Error

	if err != nil {
		return nil, err
	}

	sessions := make([]*domain.Session, len(sessionModels))
	for i, sessionModel := range sessionModels {
		sessions[i] = sessionModel.ToDomain()
	}

	return sessions, nil
}

func (r *SessionRepository) Update(ctx context.Context, session *domain.Session) (*domain.Session, error) {
	sessionModel := &model.SessionModel{}
	sessionModel.FromDomain(session)

	err := r.db.WithContext(ctx).Save(sessionModel).Error
	if err != nil {
		return nil, err
	}

	return sessionModel.ToDomain(), nil
}

func (r *SessionRepository) Delete(ctx context.Context, sessionID uint) error {
	err := r.db.WithContext(ctx).
		Where("id = ?", sessionID).
		Delete(&model.SessionModel{}).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *SessionRepository) DeleteByUserID(ctx context.Context, userID uint) error {
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.SessionModel{}).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *SessionRepository) DeleteByRefreshToken(ctx context.Context, refreshToken string) error {
	err := r.db.WithContext(ctx).
		Where("refresh_token = ?", refreshToken).
		Delete(&model.SessionModel{}).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *SessionRepository) DeleteExpiredSessions(ctx context.Context) error {
	err := r.db.WithContext(ctx).
		Where("expires_at < ?", gorm.Expr("NOW()")).
		Delete(&model.SessionModel{}).Error

	if err != nil {
		return err
	}

	return nil
}
