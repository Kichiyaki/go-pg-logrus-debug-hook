package querylogger

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"
	"strings"
)

type Logger struct {
	Log            logrus.FieldLogger
	MaxQueryLength int
}

var _ pg.QueryHook = (*Logger)(nil)

func (logger Logger) BeforeQuery(ctx context.Context, evt *pg.QueryEvent) (context.Context, error) {
	if logger.Log == nil {
		return ctx, nil
	}

	q, err := evt.FormattedQuery()
	if err != nil {
		return nil, err
	}
	prepared := strings.TrimSpace(string(q))
	if logger.MaxQueryLength > 0 && len(prepared) > logger.MaxQueryLength {
		prepared = prepared[0:logger.MaxQueryLength]
	}

	if evt.Err != nil {
		logger.Log.Errorf("%s executing a query:\n%s\n", evt.Err, prepared)
	} else {
		logger.Log.Info(prepared)
	}

	return ctx, nil
}

func (Logger) AfterQuery(context.Context, *pg.QueryEvent) error {
	return nil
}
