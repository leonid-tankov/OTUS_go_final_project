package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/leonid-tankov/OTUS_go_final_project/internal/config"
)

var (
	ErrNoRows         = pgx.ErrNoRows
	ErrNoRowsAffected = errors.New("no rows updated")
	ErrHasRotation    = errors.New("this banner rotation already exists in the slot")
)

type Storage struct {
	Dsn string
	DB  *pgx.Conn
}

func New(cfg *config.Config) *Storage {
	return &Storage{
		Dsn: fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			cfg.Postgres.Username,
			cfg.Postgres.Password,
			cfg.Postgres.Host,
			cfg.Postgres.Port,
			cfg.Postgres.Database,
		),
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	var err error
	s.DB, err = pgx.Connect(ctx, s.Dsn)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.DB.Close(ctx)
}

func (s *Storage) GetSocialDemGroups(ctx context.Context) ([]int, error) {
	statement := `
SELECT id FROM social_dem_groups
`
	rows, err := s.DB.Query(ctx, statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (s *Storage) GetClicks(ctx context.Context, click Click) ([]Counter, error) {
	statement := `
SELECT banner_id, count FROM clicks
WHERE slot_id = $1 and social_dem_group_id = $2
`
	rows, err := s.DB.Query(ctx, statement, click.SlotID, click.SocialDemGroupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var counters []Counter
	for rows.Next() {
		var id, count int
		err = rows.Scan(&id, &count)
		if err != nil {
			return nil, err
		}
		counters = append(counters, Counter{int64(id), count})
	}
	if counters == nil {
		return nil, ErrNoRows
	}
	return counters, nil
}

func (s *Storage) CheckBanner(ctx context.Context, id int64) error {
	statement := `
SELECT id FROM banners
WHERE id = $1
`
	var result int
	err := s.DB.QueryRow(ctx, statement, id).Scan(&result)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) CheckSlot(ctx context.Context, id int64) error {
	statement := `
SELECT id FROM slots
WHERE id = $1
`
	var result int
	err := s.DB.QueryRow(ctx, statement, id).Scan(&result)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) CheckSocialDemGroup(ctx context.Context, id int64) error {
	statement := `
SELECT id FROM social_dem_groups
WHERE id = $1
`
	var result int
	err := s.DB.QueryRow(ctx, statement, id).Scan(&result)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) CheckClicks(ctx context.Context, click Click) error {
	statement := `
SELECT COUNT(*) FROM clicks
WHERE banner_id = $1 and slot_id = $2
`
	var count int
	err := s.DB.QueryRow(ctx, statement, click.BannerID, click.SlotID).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrHasRotation
	}
	return nil
}

func (s *Storage) InsertClicks(ctx context.Context, click Click) error {
	statement := `
INSERT INTO clicks (banner_id, slot_id, social_dem_group_id, count)
VALUES ($1, $2, $3, 0)
`
	_, err := s.DB.Exec(ctx, statement, click.BannerID, click.SlotID, click.SocialDemGroupID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) UpdateClicks(ctx context.Context, click Click) error {
	statement := `
UPDATE clicks
SET count = count + 1
WHERE banner_id = $1 and slot_id = $2 and social_dem_group_id = $3
`
	tag, err := s.DB.Exec(ctx, statement, click.BannerID, click.SlotID, click.SocialDemGroupID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNoRowsAffected
	}
	return nil
}

func (s *Storage) DeleteClicks(ctx context.Context, click Click) error {
	statement := `
DELETE FROM clicks
WHERE banner_id = $1 and slot_id = $2
`
	tag, err := s.DB.Exec(ctx, statement, click.BannerID, click.SlotID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNoRowsAffected
	}
	return nil
}
