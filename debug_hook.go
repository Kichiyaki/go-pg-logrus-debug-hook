package gopglogrusdebughook

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"
	"strings"
)

type DebugHook struct {
	Entry          *logrus.Entry
	MaxQueryLength int
}

var _ pg.QueryHook = (*DebugHook)(nil)

func (logger DebugHook) BeforeQuery(ctx context.Context, evt *pg.QueryEvent) (context.Context, error) {
	q, err := evt.FormattedQuery()
	if err != nil {
		return nil, err
	}

	if evt.Err != nil {
		logger.Entry.Errorf("%s executing a query:\n%s\n", evt.Err, q)
	} else {
		prepared := strings.TrimSpace(string(q))
		if logger.MaxQueryLength > 0 && len(prepared) > logger.MaxQueryLength {
			prepared = prepared[0:logger.MaxQueryLength]
		}
		logger.Entry.Info(prepared)
	}

	return ctx, nil
}

func (DebugHook) AfterQuery(context.Context, *pg.QueryEvent) error {
	return nil
}
