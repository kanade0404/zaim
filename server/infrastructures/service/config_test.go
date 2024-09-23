package service

import (
	"context"
	"github.com/dghubble/oauth1"
	"github.com/google/go-cmp/cmp"
	"github.com/kanade0404/zaim/server/domains"
	"github.com/kanade0404/zaim/server/tests"
	"testing"
	"time"
)

func TestConfigService_FindByUserName(t *testing.T) {
	ctx := context.TODO()
	_, db, err := tests.SetUp(ctx, t)
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}
	user := &domains.User{ID: 1, Name: "test"}
	if _, err := db.NewInsert().Model(user).Exec(ctx); err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}
	anotherUser := &domains.User{ID: 2, Name: "another"}
	if _, err := db.NewInsert().Model(anotherUser).Exec(ctx); err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}
	now := time.Now()
	testRecords := []struct {
		UserID  int64
		zaimApp *domains.ZaimApplication
	}{
		{
			UserID: user.ID,
			zaimApp: &domains.ZaimApplication{
				UserID:         user.ID,
				ConsumerKey:    "consumer_key",
				ConsumerSecret: "consumer_secret",
				EnableZaimApplicationEvents: []*domains.EnableZaimApplicationEvent{
					{
						Enabled: now,
					},
				},
				ZaimOauth: []*domains.ZaimOAuth{
					{
						Token:  "token",
						Secret: "secret",
						EnableZaimOAuthEvent: []*domains.EnableZaimOAuthEvent{
							{
								Enabled: now,
							},
						},
					},
					{
						Token:  "old_token",
						Secret: "old_secret",
						EnableZaimOAuthEvent: []*domains.EnableZaimOAuthEvent{
							{
								Enabled: now.Add(-1 * time.Hour),
							},
						},
					},
				},
			},
		},
		{
			UserID: user.ID,
			zaimApp: &domains.ZaimApplication{
				UserID:         user.ID,
				ConsumerKey:    "consumer_key2",
				ConsumerSecret: "consumer_secret2",
				EnableZaimApplicationEvents: []*domains.EnableZaimApplicationEvent{
					{
						Enabled: now.Add(-1 * time.Hour),
					},
				},
				ZaimOauth: []*domains.ZaimOAuth{
					{
						Token:  "token2",
						Secret: "secret2",
						EnableZaimOAuthEvent: []*domains.EnableZaimOAuthEvent{
							{
								Enabled: now.Add(-1 * time.Hour),
							},
						},
					},
				},
			},
		},
	}
	for _, testRecord := range testRecords {
		if _, err := db.NewInsert().Model(testRecord.zaimApp).Exec(ctx); err != nil {
			t.Fatalf("failed to insert zaim application: %v", err)
		}
		for _, event := range testRecord.zaimApp.EnableZaimApplicationEvents {
			e := &domains.EnableZaimApplicationEvent{
				ZaimAppID: testRecord.zaimApp.ID,
				Enabled:   event.Enabled,
			}
			if _, err := db.NewInsert().Model(e).Exec(ctx); err != nil {
				t.Fatalf("failed to insert enable zaim app event: %v", err)
			}
		}
		for _, oauth := range testRecord.zaimApp.ZaimOauth {
			o := &domains.ZaimOAuth{
				Token:     oauth.Token,
				Secret:    oauth.Secret,
				ZaimAppID: testRecord.zaimApp.ID,
			}
			if _, err := db.NewInsert().Model(o).Exec(ctx); err != nil {
				t.Fatalf("failed to insert zaim oauth: %v", err)
			}
			for _, event := range oauth.EnableZaimOAuthEvent {
				e := &domains.EnableZaimOAuthEvent{
					ZaimOAuthID: o.ID,
					Enabled:     event.Enabled,
				}
				if _, err := db.NewInsert().Model(e).Exec(ctx); err != nil {
					t.Fatalf("failed to insert enable zaim oauth event: %v", err)
				}
			}
		}
	}
	testCases := []struct {
		name    string
		userID  int64
		want    *domains.Config
		wantErr bool
	}{
		{
			name:   "success",
			userID: user.ID,
			want: &domains.Config{
				OAuthConfig: &oauth1.Config{
					ConsumerKey:    "consumer_key",
					ConsumerSecret: "consumer_secret",
				},
				OAuthToken: domains.OAuthToken{
					Token:  "token",
					Secret: "secret",
				},
			},
		},
		{
			name:    "not found",
			userID:  anotherUser.ID,
			wantErr: true,
		},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := ConfigService{
				baseService: newTestBaseService(db),
			}
			got, err := c.FindByUserID(ctx, tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("(-want+got)\n%s", diff)
			}
		})
	}
}
