package domains

import (
	"context"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Category struct {
	bun.BaseModel          `bun:"table:category,alias:c"`
	ID                     uuid.UUID                `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	CategoryID             int64                    `bun:"category_id,notnull,unique:unq_category"`
	UserID                 int64                    `bun:"user_id,notnull,unique:unq_category"`
	User                   *User                    `bun:"rel:belongs-to,join:user_id=id"`
	Genres                 []*Genre                 `bun:"rel:has-many,join:id=category_id"`
	CategoryModifiedEvents []*CategoryModifiedEvent `bun:"rel:has-many,join:id=category_id"`
	ActiveCategories       []*ActiveCategory        `bun:"rel:has-many,join:id=category_id"`
	InActiveCategories     []*InActiveCategory      `bun:"rel:has-many,join:id=category_id"`
}

var _ Domain = (*Category)(nil)

func (c *Category) Indexes() []*Index {
	return []*Index{
		createIndex("category", []string{"category_id", "user_id"}),
	}
}

type CategoryModifiedEvent struct {
	bun.BaseModel `bun:"table:category_modified_event,alias:cme"`
	ID            uuid.UUID `bun:"id,type:uuid,pk,default:gen_random_uuid()"`
	CategoryID    uuid.UUID `bun:"category_id,type:uuid,notnull,unique:unq_category_modified_event"`
	Category      *Category `bun:"rel:belongs-to,join:category_id=id"`
	Name          string    `bun:"name,notnull,unique:unq_category_modified_event"`
	ModeID        string    `bun:"mode_id,notnull,unique:unq_category_modified_event"`
	Mode          *Mode     `bun:"rel:belongs-to,join:mode_id=id"`
	Sort          int       `bun:"sort,notnull,unique:unq_category_modified_event"`
	Modified      time.Time `bun:"modified,notnull,default:current_timestamp,unique:unq_category_modified_event"`
}

var _ Domain = (*CategoryModifiedEvent)(nil)

func (c *CategoryModifiedEvent) Indexes() []*Index {
	return []*Index{
		createIndex("category_modified_event", []string{"category_id"}),
		createIndex("category_modified_event", []string{"name"}),
		createIndex("category_modified_event", []string{"mode_id"}),
		createIndex("category_modified_event", []string{"modified"}),
	}
}

type ActiveCategory struct {
	bun.BaseModel `bun:"table:active_category,alias:ac"`
	ID            uuid.UUID `bun:"id,type:uuid,pk,default:gen_random_uuid()"`
	CategoryID    uuid.UUID `bun:"category_id,type:uuid,notnull,unique:unq_active_category"`
	Category      *Category `bun:"rel:has-one,join:category_id=id"`
	Activated     time.Time `bun:"activated,notnull,default:current_timestamp,unique:unq_active_category"`
}

var _ Domain = (*ActiveCategory)(nil)

func (a *ActiveCategory) Indexes() []*Index {
	return []*Index{
		createIndex("active_category", []string{"category_id", "activated"}),
	}
}

type InActiveCategory struct {
	bun.BaseModel `bun:"table:inactive_category,alias:iac"`
	ID            uuid.UUID `bun:"id,type:uuid,pk,default:gen_random_uuid()"`
	CategoryID    uuid.UUID `bun:"category_id,type:uuid,notnull,unique:unq_inactive_category"`
	Category      *Category `bun:"rel:has-one,join:category_id=id"`
	InActivated   time.Time `bun:"inactivated_at,notnull,default:current_timestamp,unique:unq_inactive_category"`
}

var _ Domain = (*InActiveCategory)(nil)

func (i *InActiveCategory) Indexes() []*Index {
	return []*Index{
		createIndex("inactive_category", []string{"category_id", "inactivated_at"}),
	}
}

type FindByCategoryIDInput struct {
	CategoryID int
	UserID     int64
}
type SaveCategoryInput struct {
	ID     int
	Name   string
	ModeID string
	Sort   int
	Active int
	UserID int64
}

type CategoryRepository interface {
	FindByCategoryID(ctx context.Context, input FindByCategoryIDInput) (*Category, error)
	Save(ctx context.Context, input SaveCategoryInput) error
}
