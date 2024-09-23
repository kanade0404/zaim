package domains

import (
	"context"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

const PaymentMode = "payment"
const IncomeMode = "income"
const TransferMode = "transfer"

type Mode struct {
	ID                     string                   `bun:"id,pk"`
	CategoryModifiedEvents []*CategoryModifiedEvent `bun:"rel:has-many,join:id=mode_id"`
}

func (m *Mode) Indexes() []*Index {
	return nil
}

func (m *Mode) BeforeAppendModel(_ context.Context, _ schema.Query) error {
	switch m.ID {
	case PaymentMode, IncomeMode, TransferMode:
		return nil
	default:
		return fmt.Errorf("invalid mode: %s", m.ID)
	}
}

var _ bun.BeforeAppendModelHook = (*Mode)(nil)
