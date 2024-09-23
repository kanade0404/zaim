package service

import (
	"context"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/kanade0404/zaim/server/domains"
	"github.com/kanade0404/zaim/server/tests"
	"github.com/uptrace/bun"
	"testing"
	"time"
)

type testSeedEvent struct {
	Name     string
	Sort     int
	Modified time.Time
}
type testSeedB43 struct {
	Modified time.Time
}
type testSeed struct {
	AccountID   int64
	UserID      int64
	Events      []testSeedEvent
	Activated   []time.Time
	Inactivated []time.Time
	B43         *testSeedB43
}

func insertTestSeeds(ctx context.Context, db bun.IDB, testRecords []testSeed) error {
	for _, record := range testRecords {
		a := &domains.Account{
			AccountID: record.AccountID,
			UserID:    record.UserID,
		}
		if _, err := db.NewInsert().Model(a).Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert account: %w\naccount:%+v", err, a)
		}
		var account domains.Account
		if err := db.NewSelect().Model(&account).
			Where("account_id = ?", record.AccountID).
			Where("user_id = ?", record.UserID).
			Scan(ctx); err != nil {
			return fmt.Errorf("failed to select account: %w", err)
		}
		for _, event := range record.Events {
			ame := &domains.AccountModifiedEvent{
				AccountID: account.ID,
				Sort:      event.Sort,
				Name:      event.Name,
				Modified:  event.Modified,
			}
			if _, err := db.NewInsert().Model(ame).Exec(ctx); err != nil {
				return fmt.Errorf("failed to insert account modified event: %w,\naccount_modified_event:%+v", err, ame)
			}
		}
		for _, a := range record.Activated {
			if _, err := db.NewInsert().Model(&domains.ActiveAccount{
				AccountID: account.ID,
				Activated: a,
			}).Exec(ctx); err != nil {
				return fmt.Errorf("failed to insert active account: %w", err)
			}
		}
		for _, i := range record.Inactivated {
			if _, err := db.NewInsert().Model(&domains.InActiveAccount{
				AccountID:   account.ID,
				InActivated: i,
			}).Exec(ctx); err != nil {
				return fmt.Errorf("failed to insert active account: %w", err)
			}
		}
		if record.B43 != nil {
			if _, err := db.NewInsert().Model(&domains.B43Account{AccountID: account.ID, Modified: record.B43.Modified}).Exec(ctx); err != nil {
				return fmt.Errorf("failed to insert b43 account: %w", err)
			}
		}
	}

	return nil
}

func TestAccountService_Save(t *testing.T) {
	t.Parallel()
	type want struct {
		accountCount              int
		accountModifiedEventCount int
		activeAccountCount        int
		inactiveAccountCount      int
	}
	ctx := context.TODO()
	container, db, err := tests.SetUp(ctx, t)
	if err != nil {
		t.Errorf("failed to setup: %v", err)
		return
	}
	user := &domains.User{ID: 1, Name: "test"}
	if _, err := db.NewInsert().Model(user).Exec(ctx); err != nil {
		t.Errorf("failed to insert user: %v", err)
		return
	}
	testCases := []struct {
		name  string
		input domains.SaveAccountInput
		want
		wantErr bool
	}{
		{
			name: "success",
			input: domains.SaveAccountInput{
				ID:     1,
				Name:   "test",
				Sort:   1,
				Active: 1,
			},
			want: want{
				accountCount:              1,
				accountModifiedEventCount: 1,
				activeAccountCount:        1,
				inactiveAccountCount:      0,
			},
		},
		{
			name: "存在しないユーザーなので登録失敗",
			input: domains.SaveAccountInput{
				ID:     2,
				UserID: 2,
				Name:   "unknown",
				Sort:   1,
				Active: 1,
			},
			want: want{
				accountCount:              0,
				accountModifiedEventCount: 0,
				activeAccountCount:        0,
				inactiveAccountCount:      0,
			},
			wantErr: true,
		},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				t.Logf("cleanup:%s", tt.name)
				if err := container.Restore(ctx); err != nil {
					t.Errorf("failed to restore: %v", err)
					return
				}
			})
			a := AccountService{
				baseService: newTestBaseService(db),
			}
			tt.input.UserID = user.ID
			if err := a.Save(ctx, tt.input); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Run("account count", func(t *testing.T) {
				accountCount, err := db.NewSelect().Model((*domains.Account)(nil)).Where("account_id = ?", tt.input.ID).Where("user_id = ?", tt.input.UserID).Count(ctx)
				if err != nil {
					t.Errorf("failed to count account: %v", err)
				} else {
					if diff := cmp.Diff(tt.accountCount, accountCount); diff != "" {
						t.Errorf("(-want, +got)\n%s", diff)
					}
				}
			})
			t.Run("account modified event count", func(t *testing.T) {
				accountModifiedEventCount, err := db.NewSelect().Model((*domains.AccountModifiedEvent)(nil)).Where("name = ?", tt.input.Name).Where("sort = ?", tt.input.Sort).Count(ctx)
				if err != nil {
					t.Errorf("failed to count account modified event: %v", err)
				} else {
					if diff := cmp.Diff(tt.want.accountModifiedEventCount, accountModifiedEventCount); diff != "" {
						t.Errorf("(-want, +got)\n%s", diff)
					}
				}
			})
			t.Run("active account count", func(t *testing.T) {
				activeAccountCount, err := db.NewSelect().Model((*domains.ActiveAccount)(nil)).Relation("Account", func(query *bun.SelectQuery) *bun.SelectQuery {
					return query.Where("account_id = ?", tt.input.ID).Where("account.user_id = ?", tt.input.UserID)
				}).Count(ctx)
				if err != nil {
					t.Errorf("failed to count active account: %v", err)
				} else {
					if diff := cmp.Diff(tt.want.activeAccountCount, activeAccountCount); diff != "" {
						t.Errorf("(-want, +got)\n%s", diff)
					}
				}
			})
			t.Run("inactive account count", func(t *testing.T) {
				inactiveAccount, err := db.NewSelect().Model((*domains.InActiveAccount)(nil)).Relation("Account", func(query *bun.SelectQuery) *bun.SelectQuery {
					return query.Where("account_id = ?", tt.input.ID).Where("user_id = ?", tt.input.UserID)
				}).Count(ctx)
				if err != nil {
					t.Errorf("failed to count inactive account: %v", err)
				} else {
					if diff := cmp.Diff(tt.want.inactiveAccountCount, inactiveAccount); diff != "" {
						t.Errorf("(-want, +got)\n%s", diff)
					}
				}
			})
		})
	}
}

func TestAccountService_FindByName(t *testing.T) {
	ctx := context.TODO()
	_, db, err := tests.SetUp(ctx, t)
	if err != nil {
		t.Fatalf("failed to setup: %v", err)
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
	testRecords := []testSeed{
		{
			// 1-1. アクティブなアカウント
			AccountID: 1,
			UserID:    user.ID,
			Events: []testSeedEvent{
				{
					Name:     "test",
					Sort:     1,
					Modified: now.Add(time.Second * -1),
				},
			},
			Activated:   []time.Time{now},
			Inactivated: []time.Time{now.Add(time.Second * -1)},
		},
		{
			// 1-2. アクティブなアカウント
			AccountID:   2,
			UserID:      user.ID,
			Events:      nil,
			Activated:   []time.Time{now.Add(time.Second * -2), now},
			Inactivated: []time.Time{now.Add(time.Second * -1)},
		},
		{
			// 1-3. アクティブなアカウント
			AccountID: 3,
			UserID:    user.ID,
			Events: []testSeedEvent{
				{
					Name:     "test2-1",
					Sort:     2,
					Modified: now.Add(time.Second * -1),
				},
				{
					Name:     "test2-2",
					Sort:     3,
					Modified: now,
				},
			},
			Activated:   []time.Time{now},
			Inactivated: []time.Time{now.Add(time.Second * -1)},
		},
		{
			// 2. アクティブだが違うユーザーのアカウント
			AccountID: 4,
			UserID:    anotherUser.ID,
			Events: []testSeedEvent{
				{
					Name:     "another",
					Sort:     1,
					Modified: now,
				},
			},
			Activated: []time.Time{now},
		},
		{
			// 3. 非アクティブなアカウント
			AccountID: 5,
			UserID:    user.ID,
			Events: []testSeedEvent{
				{
					Name:     "inactive_1",
					Sort:     1,
					Modified: now,
				},
			},
			Inactivated: []time.Time{now},
		},
		{
			// 3. 非アクティブなアカウント
			AccountID: 6,
			UserID:    user.ID,
			Events: []testSeedEvent{
				{
					Name:     "inactive_2",
					Sort:     1,
					Modified: now,
				},
			},
			Activated:   []time.Time{now.Add(time.Second * -1)},
			Inactivated: []time.Time{now.Add(time.Second * -2), now},
		},
	}
	if err := insertTestSeeds(ctx, db, testRecords); err != nil {
		t.Fatalf("failed to insert test seeds: %v", err)
	}
	testCases := []struct {
		name    string
		input   domains.FindByNameAccountInput
		want    *domains.Account
		wantErr bool
	}{
		{
			name: "IsActiveOnlyがtrueの場合はアクティブなアカウントのみ取得できる",
			input: domains.FindByNameAccountInput{
				UserID:       user.ID,
				Name:         "test",
				IsActiveOnly: true,
			},
			want: &domains.Account{
				AccountID: 1,
				UserID:    user.ID,
			},
		},
		{
			name: "IsActiveOnlyがfalseでもアクティブなアカウントが取得できる",
			input: domains.FindByNameAccountInput{
				UserID: user.ID,
				Name:   "test2-2",
			},
			want: &domains.Account{
				AccountID: 3,
				UserID:    user.ID,
			},
		},
		{
			name: "違うユーザーのアカウントは取得できない",
			input: domains.FindByNameAccountInput{
				UserID:       user.ID,
				Name:         "another",
				IsActiveOnly: true,
			},
			wantErr: true,
		},
		{
			name: "非アクティブなアカウントは取得できない",
			input: domains.FindByNameAccountInput{
				UserID:       user.ID,
				Name:         "inactive_2",
				IsActiveOnly: true,
			},
			wantErr: true,
		},
		{
			name: "非アクティブだがIsActiveOnlyがfalseなので取得できる",
			input: domains.FindByNameAccountInput{
				UserID: user.ID,
				Name:   "inactive_2",
			},
			want: &domains.Account{
				AccountID: 6,
				UserID:    user.ID,
			},
		},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a := AccountService{
				baseService: newTestBaseService(db),
			}
			got, err := a.FindByName(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(domains.Account{}, "ID")); diff != "" {
				t.Errorf("(-want+get)\n%s", diff)
			}
		})
	}
}

func TestAccountService_FindB43(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	_, db, err := tests.SetUp(ctx, t)
	if err != nil {
		t.Fatalf("failed to setup: %v", err)
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
	testSeeds := []testSeed{
		{
			// 1-1. アクティブなアカウント
			AccountID: 1,
			UserID:    user.ID,
			Events: []testSeedEvent{
				{
					Name:     "test",
					Sort:     1,
					Modified: now.Add(time.Second * -1),
				},
			},
			Activated:   []time.Time{now},
			Inactivated: []time.Time{now.Add(time.Second * -1)},
			B43:         &testSeedB43{Modified: now.Add(time.Second * -1)},
		},
		{
			// 1-2. アクティブなアカウント
			AccountID:   2,
			UserID:      user.ID,
			Events:      nil,
			Activated:   []time.Time{now.Add(time.Second * -2), now},
			Inactivated: []time.Time{now.Add(time.Second * -1)},
		},
		{
			// 1-3. アクティブなアカウント
			AccountID: 3,
			UserID:    user.ID,
			Events: []testSeedEvent{
				{
					Name:     "test2-1",
					Sort:     2,
					Modified: now.Add(time.Second * -1),
				},
				{
					Name:     "test2-2",
					Sort:     3,
					Modified: now,
				},
			},
			Activated:   []time.Time{now},
			Inactivated: []time.Time{now.Add(time.Second * -1)},
			B43:         &testSeedB43{Modified: now},
		},
		{
			// 2. アクティブだが違うユーザーのアカウント
			AccountID: 4,
			UserID:    anotherUser.ID,
			Events: []testSeedEvent{
				{
					Name:     "another",
					Sort:     1,
					Modified: now,
				},
			},
			Activated: []time.Time{now},
			B43:       &testSeedB43{Modified: now.Add(time.Second * 2)},
		},
		{
			// 3. 非アクティブなアカウント
			AccountID: 5,
			UserID:    user.ID,
			Events: []testSeedEvent{
				{
					Name:     "inactive_1",
					Sort:     1,
					Modified: now,
				},
			},
			Inactivated: []time.Time{now},
		},
		{
			// 3. 非アクティブなアカウント
			AccountID: 6,
			UserID:    user.ID,
			Events: []testSeedEvent{
				{
					Name:     "inactive_2",
					Sort:     1,
					Modified: now,
				},
			},
			Activated:   []time.Time{now.Add(time.Second * -1)},
			Inactivated: []time.Time{now.Add(time.Second * -2), now},
			B43:         &testSeedB43{Modified: now.Add(time.Second)},
		},
	}
	if err := insertTestSeeds(ctx, db, testSeeds); err != nil {
		t.Fatalf("failed to insert test seeds: %v", err)
	}
	testCases := []struct {
		name    string
		input   domains.FindB43Input
		want    *domains.Account
		wantErr bool
	}{
		{
			name: "アクティブなアカウントのみ取得できる",
			input: domains.FindB43Input{
				UserID: user.ID,
			},
			want: &domains.Account{
				AccountID: 3,
				UserID:    user.ID,
			},
		},
		{
			name: "非アクティブなアカウントも取得できる",
			input: domains.FindB43Input{
				UserID:          user.ID,
				IncludeInActive: true,
			},
			want: &domains.Account{
				AccountID: 6,
				UserID:    user.ID,
			},
		},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a := AccountService{
				baseService: newTestBaseService(db),
			}
			got, err := a.FindB43(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindB43() error = %v, wantErr %v, got %+v", err, tt.wantErr, got)
				return
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(domains.Account{}, "ID")); diff != "" {
				t.Errorf("(-want+get)\n%s", diff)
			}
		})
	}
}
