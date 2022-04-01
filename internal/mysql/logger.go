package mysql

import (
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"
)

type QueryExecContext interface {
	SelectContext(context.Context, interface{}, string, ...interface{}) error
	GetContext(context.Context, interface{}, string, ...interface{}) error
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
}

type queryLogger struct {
	queryExecr QueryExecContext
	logger     *logrus.Logger
}

func NewQueryLogger(execer QueryExecContext, logger *logrus.Logger) *queryLogger {
	return &queryLogger{
		queryExecr: execer,
		logger:     logger,
	}
}

func (s *queryLogger) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	s.logger.WithField("service", "mysql").WithField("args", args).Debug(query)
	return s.queryExecr.SelectContext(ctx, dest, query, args...)
}

func (s *queryLogger) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	s.logger.WithField("service", "mysql").Debug(query)
	return s.queryExecr.GetContext(ctx, dest, query, args...)
}

func (s *queryLogger) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	s.logger.WithField("service", "mysql").Debug(query)
	return s.queryExecr.ExecContext(ctx, query, args...)
}
