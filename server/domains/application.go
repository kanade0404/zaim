package domains

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type ZaimApplication struct {
	bun.BaseModel               `bun:"table:zaim_app,alias:za"`
	ID                          uuid.UUID                     `bun:"id,type:uuid,pk,default:gen_random_uuid()"`
	ConsumerKey                 string                        `bun:"consumer_key,notnull,unique:unq_zaim_app"`
	ConsumerSecret              string                        `bun:"consumer_secret,notnull,unique:unq_zaim_app"`
	UserID                      int64                         `bun:"user_id,notnull,unique:unq_zaim_app"`
	User                        *User                         `bun:"rel:belongs-to,join:user_id=id"`
	EnableZaimApplicationEvents []*EnableZaimApplicationEvent `bun:"rel:has-many,join:id=zaim_app_id"`
	ZaimOauth                   []*ZaimOAuth                  `bun:"rel:has-many,join:id=zaim_app_id"`
}

func (z *ZaimApplication) Indexes() []*Index {
	return []*Index{
		createIndex("zaim_app", []string{"user_id"}),
	}
}

type EnableZaimApplicationEvent struct {
	bun.BaseModel `bun:"table:enable_zaim_app_event,alias:eza"`
	ID            uuid.UUID `bun:"id,type:uuid,pk,default:gen_random_uuid()"`
	ZaimAppID     uuid.UUID `bun:"zaim_app_id,type:uuid,notnull,unique:unq_enable_zaim_app_event"`
	//ZaimApp       *ZaimApplication `bun:"rel:belongs-to,join:zaim_app_id=id"`
	Enabled time.Time `bun:"enabled,notnull,default:current_timestamp,unique:unq_enable_zaim_app_event"`
}

func (e *EnableZaimApplicationEvent) Indexes() []*Index {
	return nil
}

type ZaimOAuth struct {
	bun.BaseModel        `bun:"table:zaim_oauth,alias:zo"`
	ID                   uuid.UUID               `bun:"id,type:uuid,pk,default:gen_random_uuid()"`
	Token                string                  `bun:"token,notnull,unique:unq_zaim_oauth"`
	Secret               string                  `bun:"secret,notnull,unique:unq_zaim_oauth"`
	ZaimAppID            uuid.UUID               `bun:"zaim_app_id,type:uuid,notnull,unique:unq_zaim_oauth"`
	ZaimApplication      *ZaimApplication        `bun:"rel:belongs-to,join:zaim_app_id=id"`
	EnableZaimOAuthEvent []*EnableZaimOAuthEvent `bun:"rel:has-many,join:id=zaim_oauth_id"`
}

func (z *ZaimOAuth) Indexes() []*Index {
	return []*Index{
		createIndex("zaim_oauth", []string{"zaim_app_id"}),
	}
}

type EnableZaimOAuthEvent struct {
	bun.BaseModel `bun:"table:enable_zaim_oauth_event,alias:ezo"`
	ID            uuid.UUID  `bun:"id,type:uuid,pk,default:gen_random_uuid()"`
	ZaimOAuthID   uuid.UUID  `bun:"zaim_oauth_id,type:uuid,notnull,unique:unq_enable_zaim_oauth_event"`
	ZaimOAuth     *ZaimOAuth `bun:"rel:belongs-to,join:zaim_oauth_id=id"`
	Enabled       time.Time  `bun:"enabled,notnull,default:current_timestamp,unique:unq_enable_zaim_oauth_event"`
}

func (e *EnableZaimOAuthEvent) Indexes() []*Index {
	return nil
}

type ZaimRepository interface {
	FindZaimApplicationByUserID(userID int64) ([]*ZaimApplication, error)
	FindZaimOAuthByUserID(userID int64) ([]*ZaimOAuth, error)
}
