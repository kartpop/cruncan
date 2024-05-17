package retry

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-redsync/redsync/v4"
)

type Job struct {
	logger *slog.Logger
	rs     *redsync.Redsync
	ticker *time.Ticker
}

func NewJob(logger *slog.Logger, rs *redsync.Redsync) *Job {
	return &Job{
		logger: logger,
		rs:     rs,
	}
}

func (j *Job) Start(ctx context.Context) {
	if j.ticker != nil {
		j.logger.WarnContext(ctx, "trying to start a CDR partner transaction retry job that is already running")
		return
	}

	interval := 5 * time.Minute
	j.ticker = time.NewTicker(interval)

	go func() {
		for t := range j.ticker.C {
			j.logger.DebugContext(ctx, fmt.Sprintf("starting transaction retry job at %v", t))
			failedPartnerTxns := GetFailedTransactions()

			for _, failedPartnerTxn := range failedPartnerTxns {
				j.ProcessSingleTxn(ctx, failedPartnerTxn)
			}

			j.logger.InfoContext(ctx, "finished CDR partner transaction retry job")
		}
	}()
}

func (j *Job) ProcessSingleTxn(ctx context.Context, failedPartnerTxn Transaction) {
	mutex := j.rs.NewMutex(failedPartnerTxn.ID, redsync.WithExpiry(60*time.Second))
	if err := mutex.Lock(); err != nil {
		j.logger.ErrorContext(ctx, fmt.Sprintf("failed to obtain redsync mutex for transacation %s: %v", failedPartnerTxn.ID, err.Error()))
		return
	}

	defer func() {
		if ok, err := mutex.Unlock(); !ok || err != nil {
			j.logger.ErrorContext(ctx, fmt.Sprintf("failed to release redsync mutex for transaction %s: %v", failedPartnerTxn.ID, err.Error()))
		}
	}()

	// lock acquired, recheck transaction state
	if IsTransactionFailed(failedPartnerTxn) {
		j.logger.DebugContext(ctx, fmt.Sprintf("retrying transaction %s", failedPartnerTxn.ID))

		// business logic for retry
		// ...
	} else {
		j.logger.DebugContext(ctx, fmt.Sprintf("transaction %s already processed", failedPartnerTxn.ID))
	}
}

func (j *Job) Stop(ctx context.Context) error {
	if j.ticker != nil {
		j.ticker.Stop()
	}

	j.ticker = nil
	return nil
}
