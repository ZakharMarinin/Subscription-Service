package usecase

import (
	"context"
	"errors"
	"log/slog"
	"testovoe/internal/config"
	"testovoe/internal/domain"
	"time"

	"github.com/google/uuid"
)

type Storage interface {
	CreateSub(ctx context.Context, userSub domain.UserSub) error
	UpdateSub(ctx context.Context, userSub domain.UserSub) error
	DeleteSub(ctx context.Context, subID uuid.UUID, userID uuid.UUID) error
	GetSubs(ctx context.Context) ([]*domain.UserSub, error)
	GetUserSubs(ctx context.Context, userID uuid.UUID) ([]*domain.UserSub, error)
	GetUserSub(ctx context.Context, subID uuid.UUID) (*domain.UserSub, error)
	GetTotalCost(ctx context.Context, userID uuid.UUID, serviceName string, from, to time.Time) (int, error)
}

type UseCase struct {
	log     *slog.Logger
	storage Storage
	cfg     *config.Config
}

func New(log *slog.Logger, storage Storage, cfg *config.Config) *UseCase {
	return &UseCase{
		log:     log,
		storage: storage,
		cfg:     cfg,
	}
}

func (u *UseCase) CreateSub(ctx context.Context, userSub domain.UserSub) error {
	const op = "usecase.CreateSub"

	if err := validatePrice(userSub.ServicePrice); err != nil {
		u.log.Error("Validation failed", "op", op, "error", err)
		return err
	}

	err := u.storage.CreateSub(ctx, userSub)
	if err != nil {
		u.log.Error("Failed to create subscription", "op", op, "error", err)
		return err
	}

	return nil
}

func (u *UseCase) UpdateSub(ctx context.Context, userSub domain.UserSub) error {
	const op = "usecase.UpdateSub"

	if err := validatePrice(userSub.ServicePrice); err != nil {
		u.log.Error("Validation failed", "op", op, "error", err)
		return err
	}

	err := u.storage.UpdateSub(ctx, userSub)
	if err != nil {
		u.log.Error("Failed to update subscription", "op", op, "error", err)
		return err
	}

	return nil
}

func (u *UseCase) DeleteSub(ctx context.Context, subID, userID uuid.UUID) error {
	const op = "usecase.DeleteSub"

	err := u.storage.DeleteSub(ctx, subID, userID)
	if err != nil {
		u.log.Error("Failed to delete subscription", "op", op, "error", err)
		return err
	}

	return nil
}

func (u *UseCase) GetSubs(ctx context.Context) ([]*domain.UserSub, error) {
	const op = "usecase.GetSubs"

	subs, err := u.storage.GetSubs(ctx)
	if err != nil {
		u.log.Error("Failed to get subscriptions", "op", op, "error", err)
		return nil, err
	}

	return subs, nil
}

func (u *UseCase) GetUserSub(ctx context.Context, subID uuid.UUID) (*domain.UserSub, error) {
	const op = "usecase.GetUserSub"

	sub, err := u.storage.GetUserSub(ctx, subID)
	if err != nil {
		u.log.Error("Failed to get subscriptions", "op", op, "error", err)
		return nil, err
	}

	return sub, nil
}

func (u *UseCase) GetUserSubs(ctx context.Context, userID uuid.UUID) ([]*domain.UserSub, error) {
	const op = "usecase.GetUserSubs"

	subs, err := u.storage.GetUserSubs(ctx, userID)
	if err != nil {
		u.log.Error("Failed to get subscriptions", "op", op, "error", err)
		return nil, err
	}

	return subs, nil
}

func (u *UseCase) GetTotalCost(ctx context.Context, userID uuid.UUID, serviceName, fromStr, toStr string) (int, error) {
	const op = "usecase.GetTotalCost"

	log := u.log.With(
		slog.String("op", op),
		slog.String("user_id", userID.String()),
		slog.String("service", serviceName),
	)

	from, err := time.Parse("01-2006", fromStr)
	if err != nil {
		log.Error("invalid from_date format", slog.String("val", fromStr))
		return 0, err
	}

	toRaw, err := time.Parse("01-2006", toStr)
	if err != nil {
		log.Error("invalid to_date format", slog.String("val", toStr))
		return 0, err
	}

	to := toRaw.AddDate(0, 1, 0).Add(-time.Second)

	cost, err := u.storage.GetTotalCost(ctx, userID, serviceName, from, to)
	if err != nil {
		log.Error("failed to get total cost from storage", slog.Any("err", err))
		return 0, err
	}

	log.Info("total cost calculated", slog.Int("result", cost))
	return cost, nil
}

func validatePrice(price int) error {
	if price < 0 {
		return errors.New("price cannot be negative")
	}

	return nil
}
