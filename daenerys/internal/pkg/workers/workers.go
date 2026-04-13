package workers

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/CodefriendOrg/daenerys/internal/config"
	"github.com/CodefriendOrg/daenerys/internal/pb/github.com/CodefriendOrg/tirion/pkg/tirion"
	"github.com/CodefriendOrg/daenerys/internal/pkg/docker"
	"github.com/CodefriendOrg/daenerys/internal/pkg/logger"
	"github.com/CodefriendOrg/daenerys/internal/pkg/store"
	"github.com/CodefriendOrg/daenerys/internal/pkg/workers/python"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivertype"
	"go.uber.org/zap/exp/zapslog"
	"golang.org/x/sync/semaphore"
)

type retryPolicy struct{}

// NextRetry .
func (r retryPolicy) NextRetry(job *rivertype.JobRow) time.Time {
	return time.Now().UTC().Add(retryTimeoutByAttempt(job.Attempt))
}

func retryTimeoutByAttempt(attempt int) time.Duration {
	return time.Duration(math.Pow(2, float64(attempt))) * time.Second
}

// RegisterRiverClient .
func RegisterRiverClient(
	ctx context.Context,
	dbPool *pgxpool.Pool,
	dockerClient *docker.Client,
	storage *store.Storage,
	tirionClient tirion.TirionClient,
	judge config.JudgeConfig,
) (*river.Client[pgx.Tx], error) {
	workers := river.NewWorkers()
	overheadBytes := judge.SandboxMemoryOverheadMB * 1024 * 1024
	sandboxSem := semaphore.NewWeighted(int64(judge.MaxConcurrentSandboxes))
	river.AddWorker[python.Args](workers, python.NewWorker(dockerClient, storage, tirionClient, overheadBytes, sandboxSem))

	client, err := river.NewClient[pgx.Tx](riverpgxv5.New(dbPool), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {
				MaxWorkers: judge.MaxRiverWorkers,
			},
		},
		RetryPolicy: retryPolicy{},
		MaxAttempts: 5,
		Workers:     workers,
		Logger:      slog.New(zapslog.NewHandler(logger.Logger.Core())),
	})
	if err != nil {
		return nil, fmt.Errorf("river.NewClient: %w", err)
	}

	err = client.Start(ctx)
	if err != nil {
		return nil, fmt.Errorf("river.NewClient")
	}

	return client, nil
}
