package storage

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"testovoe/internal/domain"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Storage struct {
	DB *pgxpool.Pool
}

func New(ctx context.Context, storagePath string) (*Storage, error) {
	const op = "storage.postgresql.SQL.NEW"

	if err := runMigrations(storagePath); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	poolConfig, err := pgxpool.ParseConfig(storagePath)
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	return &Storage{DB: db}, nil
}

func runMigrations(dbURL string) error {
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return err
	}

	db := stdlib.OpenDB(*config.ConnConfig)
	defer db.Close()

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Close() error {
	s.DB.Close()
	return nil
}

func (s *Storage) CreateSub(ctx context.Context, userSub domain.UserSub) error {
	const op = "storage.storage.CreateSub"

	query, args, err := sq.
		Insert("subscriptions").
		Columns("service_name", "sub_price", "user_id", "started_at", "ended_at").
		Values(userSub.ServiceName, userSub.ServicePrice, userSub.UserID, userSub.StartedAt, userSub.EndedAt).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateSub(ctx context.Context, userSub domain.UserSub) error {
	const op = "storage.storage.UpdateSub"

	query, args, err := sq.
		Update("subscriptions").
		SetMap(map[string]interface{}{
			"service_name": userSub.ServiceName,
			"sub_price":    userSub.ServicePrice,
			"ended_at":     userSub.EndedAt,
		}).
		Where(sq.Eq{"id": userSub.ID, "user_id": userSub.UserID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteSub(ctx context.Context, subID, userID uuid.UUID) error {
	const op = "storage.storage.DeleteSub"

	query, args, err := sq.
		Delete("subscriptions").
		Where(sq.Eq{"id": subID, "user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetSubs(ctx context.Context) ([]*domain.UserSub, error) {
	const op = "storage.storage.GetSubs"

	query, args, err := sq.
		Select("service_name", "sub_price", "user_id", "started_at", "ended_at").
		From("subscriptions").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := s.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var userSubs []*domain.UserSub

	for rows.Next() {
		var userSub domain.UserSub
		if err := rows.Scan(&userSub.ServiceName, &userSub.ServicePrice, &userSub.UserID, &userSub.StartedAt, &userSub.EndedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		userSubs = append(userSubs, &userSub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return userSubs, nil
}

func (s *Storage) GetUserSubs(ctx context.Context, userID uuid.UUID) ([]*domain.UserSub, error) {
	const op = "storage.storage.GetSubs"

	query, args, err := sq.
		Select("service_name", "sub_price", "user_id", "started_at", "ended_at").
		From("subscriptions").
		Where(sq.Eq{"user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := s.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var userSubs []*domain.UserSub

	for rows.Next() {
		var userSub domain.UserSub
		if err := rows.Scan(&userSub.ServiceName, &userSub.ServicePrice, &userSub.UserID, &userSub.StartedAt, &userSub.EndedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		userSubs = append(userSubs, &userSub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return userSubs, nil
}

func (s *Storage) GetUserSub(ctx context.Context, subID uuid.UUID) (*domain.UserSub, error) {
	const op = "storage.storage.GetUserSub"

	query, args, err := sq.
		Select("service_name", "sub_price", "user_id", "started_at", "ended_at").
		From("subscriptions").
		Where(sq.Eq{"id": subID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var userSub domain.UserSub

	err = s.DB.QueryRow(ctx, query, args...).Scan(
		&userSub.ServiceName,
		&userSub.ServicePrice,
		&userSub.UserID,
		&userSub.StartedAt,
		&userSub.EndedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &userSub, nil
}

func (s *Storage) GetTotalCost(ctx context.Context, userID uuid.UUID, serviceName string, from, to time.Time) (int, error) {
	const op = "storage.storage.GetTotalCost"

	query, args, err := sq.
		Select("COALESCE(SUM(sub_price), 0)").
		From("subscriptions").
		Where(sq.Eq{"user_id": userID}).
		Where(sq.Eq{"service_name": serviceName}).
		Where(sq.GtOrEq{"started_at": from}).
		Where(sq.LtOrEq{"started_at": to}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var total int
	err = s.DB.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return total, nil
}
