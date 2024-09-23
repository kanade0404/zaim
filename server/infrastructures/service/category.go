package service

import (
	"context"
	"fmt"
	"github.com/kanade0404/zaim/server/domains"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

type CategoryService struct {
	baseService
}

func (c CategoryService) FindByCategoryID(ctx context.Context, input domains.FindByCategoryIDInput) (*domains.Category, error) {
	c.logger.Infof("FindByCategoryID category start. categoryID: %d", input.CategoryID)
	defer c.logger.Infof("FindByCategoryID category end. categoryID: %d", input.CategoryID)
	category := &domains.Category{
		CategoryID: int64(input.CategoryID),
		UserID:     input.UserID,
	}
	if err := c.db.NewSelect().
		Model(category).
		Where("user_id = ?", input.UserID).
		Where("category_id = ?", input.CategoryID).
		Limit(1).
		Scan(ctx); err != nil {
		return nil, fmt.Errorf("failed to select category: %v", err)
	}
	return category, nil
}

func (c CategoryService) Save(ctx context.Context, input domains.SaveCategoryInput) error {
	c.logger.Infof("Save category start. input: %+v", input)
	defer c.logger.Infof("Save category end. input: %+v", input)
	// 1. 新規ならカテゴリ作成
	category := &domains.Category{
		CategoryID: int64(input.ID),
		UserID:     input.UserID,
	}
	categoryExists, err := c.db.NewSelect().
		Model(category).
		Where("user_id = ?", input.UserID).
		Where("category_id = ?", input.ID).
		Exists(ctx)
	if err != nil {
		return fmt.Errorf("failed to select category: %v", err)
	}
	if !categoryExists {
		if _, err := c.db.NewInsert().
			Model(category).
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert category: %v", err)
		}
	}
	if err := c.db.NewSelect().
		Model(category).
		Where("user_id = ?", input.UserID).
		Where("category_id = ?", input.ID).
		Limit(1).
		Scan(ctx); err != nil {
		return fmt.Errorf("failed to select category: %v", err)
	}
	// 2. カテゴリ変更イベント作成
	categoryModifiedEvent := &domains.CategoryModifiedEvent{
		CategoryID: category.ID,
		ModeID:     input.ModeID,
		Sort:       input.Sort,
		Name:       input.Name,
	}
	if _, err := c.db.NewInsert().
		Model(categoryModifiedEvent).
		Exec(ctx); err != nil {
		return fmt.Errorf("failed to insert category modified event: %v", err)
	}
	// 3. アクティブならアクティブカテゴリを登録、非アクティブなら非アクティブカテゴリを登録
	if input.Active == 1 {
		activeCategory := &domains.ActiveCategory{
			CategoryID: category.ID,
		}
		if _, err := c.db.NewInsert().
			Model(activeCategory).
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert active category: %v", err)
		}
	} else if input.Active == -1 {
		inActiveCategory := &domains.InActiveCategory{
			CategoryID: category.ID,
		}
		if _, err := c.db.NewInsert().
			Model(inActiveCategory).
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert inactive category: %v", err)
		}
	}
	return nil
}

func NewCategoryService(db bun.Tx, logger echo.Logger) domains.CategoryRepository {
	return CategoryService{baseService: newBaseService(db, logger)}
}
