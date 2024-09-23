package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kanade0404/zaim/server/commands"
	"github.com/kanade0404/zaim/server/domains"
	"github.com/kanade0404/zaim/server/driver"
	"github.com/kanade0404/zaim/server/entity/zaim"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
	"io"
)

type UpdateGenreUseCase interface {
	UpdateGenres(ctx context.Context, userName string) error
}

type updateGenreUseCase struct {
	logger echo.Logger
	uow    commands.GenreUnitOfWork
}

func (u updateGenreUseCase) UpdateGenres(ctx context.Context, userName string) error {
	u.logger.Infof("UpdateGenres start. userName: %s", userName)
	defer u.logger.Infof("UpdateGenres end. userName: %s", userName)

	if err := u.uow.Do(ctx, func(ctx context.Context, uowRepositoryManager commands.GenreRepositoryManager) error {
		user, err := uowRepositoryManager.UserRepositoryManager().FindByName(ctx, userName)
		if err != nil {
			return fmt.Errorf("failed to find user: %v", err)
		}
		cfg, err := uowRepositoryManager.ConfigRepositoryManager().FindByUserID(ctx, user.ID)
		if err != nil {
			return fmt.Errorf("failed to find config: %w", err)
		}
		zaimDriver, err := driver.NewZaimDriver(cfg.OAuthConfig.ConsumerKey, cfg.OAuthConfig.ConsumerSecret, cfg.OAuthToken.Token, cfg.OAuthToken.Secret)
		if err != nil {
			return fmt.Errorf("failed to new zaim driver: %v", err)
		}
		res, err := zaimDriver.Get("home/genre", nil)
		if err != nil {
			return fmt.Errorf("failed to get genre: %v", err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				u.logger.Errorf("failed to close body: %v", err)
			}
		}(res.Body)
		var r struct {
			Genres    []zaim.Genre `json:"genres"`
			Requested int          `json:"requested"`
		}
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to read body: %v", err)
		}
		if err := json.Unmarshal(b, &r); err != nil {
			return fmt.Errorf("failed to unmarshal: %v", err)
		}
		var (
			notHaveParentGenres []zaim.Genre
			haveParentGenres    []zaim.Genre
		)
		for i := range r.Genres {
			if r.Genres[i].ParentGenreID == 0 {
				notHaveParentGenres = append(notHaveParentGenres, r.Genres[i])
			} else {
				haveParentGenres = append(haveParentGenres, r.Genres[i])
			}
		}
		for _, genre := range notHaveParentGenres {
			category, err := uowRepositoryManager.CategoryRepositoryManager().FindByCategoryID(ctx, domains.FindByCategoryIDInput{CategoryID: genre.CategoryID, UserID: user.ID})
			if err != nil {
				return fmt.Errorf("failed to find category: %v", err)
			}
			if err := uowRepositoryManager.GenreRepositoryManager().Save(ctx, domains.SaveGenreInput{
				ID:            genre.ID,
				Name:          genre.Name,
				Sort:          genre.Sort,
				Active:        genre.Active,
				ParentGenreID: genre.ParentGenreID,
				UserID:        user.ID,
				CategoryID:    category.ID,
			}); err != nil {
				return fmt.Errorf("failed to save genre: %v", err)
			}
		}
		for _, genre := range haveParentGenres {
			category, err := uowRepositoryManager.CategoryRepositoryManager().FindByCategoryID(ctx, domains.FindByCategoryIDInput{CategoryID: genre.CategoryID, UserID: user.ID})
			if err != nil {
				return fmt.Errorf("failed to find category: %v", err)
			}
			if err := uowRepositoryManager.GenreRepositoryManager().Save(ctx, domains.SaveGenreInput{
				ID:            genre.ID,
				Name:          genre.Name,
				Sort:          genre.Sort,
				Active:        genre.Active,
				ParentGenreID: genre.ParentGenreID,
				UserID:        user.ID,
				CategoryID:    category.ID,
			}); err != nil {
				return fmt.Errorf("failed to save genre: %v", err)
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to do: %v", err)
	}

	return nil
}

func NewUpdateGenreUseCase(db *bun.DB, logger echo.Logger) UpdateGenreUseCase {
	return updateGenreUseCase{
		logger: logger,
		uow:    commands.NewGenreUnitOfWork(db, logger),
	}

}
