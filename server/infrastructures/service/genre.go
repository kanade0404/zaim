package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kanade0404/zaim/server/domains"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

type GenreService struct {
	baseService
}

func (g GenreService) FindB43Genre(ctx context.Context, input domains.FindB43Genre) (*domains.Genre, error) {
	g.logger.Infof("Find genre start. input: %+v", input)
	defer g.logger.Infof("Find genre end. input: %+v", input)
	var result struct {
		GenrePK    uuid.UUID `bun:"genre_pk"`
		GenreID    int64     `bun:"genre_id"`
		CategoryPK uuid.UUID `bun:"category_pk"`
		CategoryID int64     `bun:"category_id"`
	}
	query := g.db.NewSelect().
		Model((*domains.Genre)(nil)).
		Join("LEFT JOIN b43_genre ON b43_genre.genre_id = g.id").
		Join("LEFT JOIN active_genre ag ON ag.genre_id = g.id").
		Join("LEFT JOIN inactive_genre iag ON iag.genre_id = g.id").
		Join("LEFT JOIN category c ON c.id = g.category_id").
		Where("user_id = ?", input.UserID).
		Where("b43_genre.id IS NOT NULL").
		Where("ag.activated IS NOT NULL").
		ColumnExpr("g.id AS genre_pk, g.genre_id, c.id AS category_pk, c.category_id, max(activated) as last_activated, max(inactivated) as last_inactivated").
		Group("g.id", "g.genre_id", "c.id", "c.category_id")
	if err := g.db.NewSelect().
		TableExpr("(?) AS genres", query).
		Where("genres.last_activated is not null").
		WhereGroup(" AND ", func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.Where("genres.last_inactivated IS NULL").
				WhereOr("genres.last_activated > genres.last_inactivated")
		}).
		Column("genre_pk", "genre_id", "category_pk", "category_id").
		Scan(ctx, &result); err != nil {
		return nil, fmt.Errorf("failed to select genre: %v", err)
	}
	genre := &domains.Genre{
		ID:         result.GenrePK,
		GenreID:    result.GenreID,
		CategoryID: result.CategoryPK,
		Category: &domains.Category{
			ID:         result.CategoryPK,
			CategoryID: result.CategoryID,
		},
	}
	return genre, nil
}

func (g GenreService) Save(ctx context.Context, input domains.SaveGenreInput) error {
	g.logger.Infof("Save genre start. input: %+v", input)
	defer g.logger.Infof("Save genre end. input: %+v", input)
	genre := &domains.Genre{
		GenreID:    int64(input.ID),
		CategoryID: input.CategoryID,
	}
	genreExists, err := g.db.NewSelect().
		Model(genre).
		Where("category_id = ?", input.CategoryID).
		Where("genre_id = ?", input.ID).
		Exists(ctx)
	if err != nil {
		return fmt.Errorf("failed to select exists genre: %v", err)
	}
	if !genreExists {
		if _, err := g.db.NewInsert().
			Model(genre).
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert genre: %v", err)
		}
	}
	if err := g.db.NewSelect().
		Model(genre).
		Where("category_id = ?", input.CategoryID).
		Where("genre_id = ?", input.ID).
		Limit(1).
		Scan(ctx); err != nil {
		return fmt.Errorf("failed to select genre: %v", err)
	}
	if input.Active == 1 {
		activeGenre := &domains.ActiveGenre{
			GenreID: genre.ID,
		}
		if _, err := g.db.NewInsert().
			Model(activeGenre).
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert active genre: %v", err)
		}
	} else if input.Active == -1 {
		inActiveGenre := &domains.InActiveGenre{
			GenreID: genre.ID,
		}
		if _, err := g.db.NewInsert().
			Model(inActiveGenre).
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert inactive genre: %v", err)
		}
	}
	genreModifiedEvent := &domains.GenreModifiedEvent{
		Name:    input.Name,
		Sort:    input.Sort,
		GenreID: genre.ID,
	}
	if _, err := g.db.NewInsert().
		Model(genreModifiedEvent).
		Exec(ctx); err != nil {
		return fmt.Errorf("failed to insert genre modified event: %v", err)
	}
	if input.ParentGenreID != 0 {
		parentGenre := &domains.Genre{
			GenreID: int64(input.ParentGenreID),
		}
		if err := g.db.NewSelect().
			Model(parentGenre).
			Where("category_id = ?", input.CategoryID).
			Where("genre_id = ?", input.ParentGenreID).
			Limit(1).
			Scan(ctx); err != nil {
			return fmt.Errorf("failed to select parent genre: %v", err)
		}
		if _, err := g.db.NewInsert().
			Model(&domains.ParentGenreModifiedEvent{
				ParentGenreID: parentGenre.ID,
				ChildGenreID:  genre.ID,
			}).
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert parent genre modified event: %v", err)
		}
	}
	if input.Name == "B43取込" {
		if _, err := g.db.NewInsert().
			Model(&domains.B43Genre{
				GenreID: genre.ID,
			}).
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert b43 genre: %v", err)
		}
	}
	return nil
}

func (g GenreService) FindByGenreID(ctx context.Context, input domains.FindByGenreIDGenreInput) (*domains.GenreModifiedEvent, error) {
	genre := new(domains.GenreModifiedEvent)
	err := g.db.NewSelect().
		Model(genre).
		Relation("Genre").
		Relation("Category").
		Relation("User").
		Where("category_id = ?", input.UserID).
		Where("genre.genre_id = ?", input.GenreID).
		Order("modified DESC").
		Limit(1).
		Scan(ctx)
	return genre, err
}

func NewGenreService(db bun.Tx, logger echo.Logger) domains.GenreRepository {
	return GenreService{baseService: newBaseService(db, logger)}
}

var _ domains.GenreRepository = (*GenreService)(nil)
