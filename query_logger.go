package gopglogrusquerylogger

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"
	"strings"
)

type QueryLogger struct {
	Entry          *logrus.Entry
	MaxQueryLength int
}

var _ pg.QueryHook = (*QueryLogger)(nil)

func (logger QueryLogger) BeforeQuery(ctx context.Context, evt *pg.QueryEvent) (context.Context, error) {
	q, err := evt.FormattedQuery()
	if err != nil {
		return nil, err
	}
	prepared := strings.TrimSpace(string(q))
	if logger.MaxQueryLength > 0 && len(prepared) > logger.MaxQueryLength {
		prepared = prepared[0:logger.MaxQueryLength]
	}

	if evt.Err != nil {
		logger.Entry.Errorf("%s executing a query:\n%s\n", evt.Err, prepared)
	} else {
		logger.Entry.Info(prepared)
	}

	return ctx, nil
}

func (QueryLogger) AfterQuery(context.Context, *pg.QueryEvent) error {
	return nil
}
