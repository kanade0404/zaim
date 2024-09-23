package domains

import (
	"context"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Genre struct {
	bun.BaseModel             `bun:"table:genre,alias:g"`
	ID                        uuid.UUID                   `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	GenreID                   int64                       `bun:"genre_id,notnull,unique:unq_genre"`
	CategoryID                uuid.UUID                   `bun:"category_id,type:uuid,notnull,unique:unq_genre"`
	Category                  *Category                   `bun:"rel:belongs-to,join:category_id=id"`
	GenreModifiedEvents       []*GenreModifiedEvent       `bun:"rel:has-many,join:id=genre_id"`
	ActiveGenres              []*ActiveGenre              `bun:"rel:has-many,join:id=genre_id"`
	InActiveGenres            []*InActiveGenre            `bun:"rel:has-many,join:id=genre_id"`
	B43Genres                 []*B43Genre                 `bun:"rel:has-many,join:id=genre_id"`
	ParentGenreModifiedEvents []*ParentGenreModifiedEvent `bun:"rel:has-many,join:id=parent_genre_id"`
	ChildGenreModifiedEvents  []*ParentGenreModifiedEvent `bun:"rel:has-many,join:id=child_genre_id"`
}

var _ Domain = (*Genre)(nil)

func (g *Genre) Indexes() []*Index {
	return []*Index{
		createIndex("genre", []string{"category_id"}),
	}
}

type GenreModifiedEvent struct {
	bun.BaseModel `bun:"table:genre_modified_event,alias:gme"`
	ID            uuid.UUID `bun:"id,type:uuid,pk,default:gen_random_uuid()"`
	Name          string    `bun:"name,notnull,unique:unq_genre_modified_event"`
	Sort          int       `bun:"sort,notnull,unique:unq_genre_modified_event"`
	GenreID       uuid.UUID `bun:"genre_id,type:uuid,notnull,unique:unq_genre_modified_event"`
	Genre         *Genre    `bun:"rel:belongs-to,join:genre_id=id"`
	Modified      time.Time `bun:"modified,notnull,default:current_timestamp,unique:unq_genre_modified_event"`
}

var _ Domain = (*GenreModifiedEvent)(nil)

func (g *GenreModifiedEvent) Indexes() []*Index {
	return []*Index{
		createIndex("genre_modified_event", []string{"genre_id"}),
		createIndex("genre_modified_event", []string{"name"}),
		createIndex("genre_modified_event", []string{"modified"}),
	}
}

type ParentGenreModifiedEvent struct {
	bun.BaseModel `bun:"table:parent_genre_modified_event,alias:pgme"`
	ID            uuid.UUID `bun:"id,type:uuid,pk,default:gen_random_uuid()"`
	ChildGenreID  uuid.UUID `bun:"child_genre_id,type:uuid,notnull,unique:unq_parent_genre_modified_event"`
	ChildGenre    *Genre    `bun:"rel:belongs-to,join:child_genre_id=id"`
	ParentGenreID uuid.UUID `bun:"parent_genre_id,type:uuid,notnull,unique:unq_parent_genre_modified_event"`
	ParentGenre   *Genre    `bun:"rel:belongs-to,join:parent_genre_id=id"`
	Modified      time.Time `bun:"modified,notnull,default:current_timestamp,unique:unq_parent_genre_modified_event"`
}

func (p *ParentGenreModifiedEvent) Indexes() []*Index {
	return []*Index{
		createIndex("parent_genre_modified_event", []string{"child_genre_id"}),
		createIndex("parent_genre_modified_event", []string{"parent_genre_id"}),
		createIndex("parent_genre_modified_event", []string{"modified"}),
	}
}

var _ Domain = (*ParentGenreModifiedEvent)(nil)

type ActiveGenre struct {
	bun.BaseModel `bun:"table:active_genre,alias:ag"`
	ID            uuid.UUID `bun:"id,type:uuid,pk,default:gen_random_uuid()"`
	GenreID       uuid.UUID `bun:"genre_id,type:uuid,notnull,unique:unq_active_genre"`
	Genre         *Genre    `bun:"rel:belongs-to,join:genre_id=id"`
	Activated     time.Time `bun:"activated,notnull,default:current_timestamp,unique:unq_active_genre"`
}

var _ Domain = (*ActiveGenre)(nil)

func (a *ActiveGenre) Indexes() []*Index {
	return nil
}

type B43Genre struct {
	bun.BaseModel `bun:"table:b43_genre,alias:b43g"`
	ID            uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	GenreID       uuid.UUID `bun:"genre_id,type:uuid,notnull,unique:unq_b43_genre"`
	Genre         *Genre    `bun:"rel:belongs-to,join:genre_id=id"`
	Modified      time.Time `bun:"modified,notnull,default:current_timestamp,unique:unq_b43_genre"`
}

func (b *B43Genre) Indexes() []*Index {
	return nil
}

var _ Domain = (*B43Genre)(nil)

type InActiveGenre struct {
	bun.BaseModel `bun:"table:inactive_genre,alias:iag"`
	ID            uuid.UUID `bun:"id,type:uuid,pk,default:gen_random_uuid()"`
	GenreID       uuid.UUID `bun:"genre_id,type:uuid,notnull,unique:unq_inactive_genre"`
	Genre         *Genre    `bun:"rel:belongs-to,join:genre_id=id"`
	InActivated   time.Time `bun:"inactivated,notnull,default:current_timestamp,unique:unq_inactive_genre"`
}

var _ Domain = (*InActiveGenre)(nil)

func (i *InActiveGenre) Indexes() []*Index {
	return nil
}

type FindByGenreIDGenreInput struct {
	UserID  int
	GenreID int
}

type FindB43Genre struct {
	UserID int64
}

type SaveGenreInput struct {
	ID            int
	Name          string
	CategoryID    uuid.UUID
	Sort          int
	Active        int
	ParentGenreID int
	UserID        int64
}
type GenreRepository interface {
	FindByGenreID(ctx context.Context, input FindByGenreIDGenreInput) (*GenreModifiedEvent, error)
	FindB43Genre(ctx context.Context, input FindB43Genre) (*Genre, error)
	Save(ctx context.Context, input SaveGenreInput) error
}
