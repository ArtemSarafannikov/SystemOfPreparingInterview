package python

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CodefriendOrg/daenerys/internal/pb/github.com/CodefriendOrg/daenerys/pkg/daenerys"
	"github.com/CodefriendOrg/daenerys/internal/pb/github.com/CodefriendOrg/tirion/pkg/tirion"
	"github.com/CodefriendOrg/daenerys/internal/pkg/docker"
	"github.com/CodefriendOrg/daenerys/internal/pkg/logger"
	"github.com/CodefriendOrg/daenerys/internal/pkg/store"
	"github.com/riverqueue/river"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

// Worker .
type Worker struct {
	dockerClient           *docker.Client
	storage                *store.Storage
	tirionClient           tirion.TirionClient
	sandboxMemoryOverheadB int64
	sandboxSem             *semaphore.Weighted

	river.WorkerDefaults[Args]
}

// NewWorker .
func NewWorker(
	dockerClient *docker.Client,
	storage *store.Storage,
	tirionClient tirion.TirionClient,
	sandboxMemoryOverheadB int64,
	sandboxSem *semaphore.Weighted,
) *Worker {
	return &Worker{
		dockerClient:           dockerClient,
		storage:                storage,
		tirionClient:           tirionClient,
		sandboxMemoryOverheadB: sandboxMemoryOverheadB,
		sandboxSem:             sandboxSem,
	}
}

// Work .
func (w *Worker) Work(ctx context.Context, job *river.Job[Args]) error {
	status := daenerys.SubmissionStatus_STATUS_INTERNAL_ERROR
	// TODO: есть риск, что задача останется висеть в NEW, если проблемы с бд
	defer func() {
		// Обновляем статус для INTERNAL_ERROR только на последней попытке
		if status == daenerys.SubmissionStatus_STATUS_INTERNAL_ERROR {
			logger.Errorf(ctx, fmt.Sprintf("FAILED JUDGE SUBMISSION ID = %s", job.Args.SubmissionID))
			if job.Attempt < job.MaxAttempts {
				return
			}
		}

		_, errUpd := w.storage.UpdateSubmissionStatus(ctx, job.Args.SubmissionID, status)
		if errUpd != nil {
			logger.Errorf(ctx, fmt.Sprintf("storage.UpdateSubmissionStatus: %v", errUpd))
		}
	}()

	_, err := w.storage.UpdateSubmissionStatus(ctx, job.Args.SubmissionID, daenerys.SubmissionStatus_STATUS_JUDGING)
	if err != nil {
		return fmt.Errorf("storage.UpdateSubmissionStatus: %w", err)
	}

	resp, err := w.tirionClient.GetProblem(ctx, &tirion.GetProblemRequest{
		Id:        job.Args.ProblemID.String(),
		WithTests: true,
	})
	if err != nil {
		return fmt.Errorf("tirionClient.GetProblem: %w", err)
	}

	if err := w.sandboxSem.Acquire(ctx, 1); err != nil {
		return fmt.Errorf("sandboxSem.Acquire: %w", err)
	}
	defer w.sandboxSem.Release(1)

	imageName := "python:" + job.Args.PythonVersion
	if err := w.dockerClient.PullImage(ctx, imageName); err != nil {
		return fmt.Errorf("dockerClient.PullImage: %w", err)
	}

	problemMemBytes := resp.GetProblem().GetMemoryLimitKb() * 1024
	containerMemBytes := problemMemBytes + w.sandboxMemoryOverheadB
	containerID, err := w.dockerClient.CreateContainer(ctx, imageName, containerMemBytes)
	if err != nil {
		return fmt.Errorf("dockerClient.CreateContainer: %w", err)
	}
	defer func() {
		errRemove := w.dockerClient.RemoveContainer(ctx, containerID)
		if errRemove != nil {
			logger.Errorf(ctx, "dockerClient.RemoveContainer", zap.Error(errRemove))
		}
	}()

	err = w.dockerClient.CopyCodeToContainer(ctx, containerID, "/app", "main.py", job.Args.Code)
	if err != nil {
		return fmt.Errorf("dockerClient.CopyCodeToContainer: %w", err)
	}

	status, err = w.runTests(ctx, containerID, resp.GetTests(), resp.GetProblem())
	if err != nil {
		return fmt.Errorf("runTests: %w", err)
	}
	return nil
}

func (w *Worker) runTests(ctx context.Context, containerID string, tests []*tirion.Test, problem *tirion.Problem) (daenerys.SubmissionStatus, error) {
	for _, test := range tests {
		maxExecTime := time.Duration(problem.GetTimeLimitMs()) * time.Millisecond
		if maxExecTime == 0 {
			maxExecTime = 10 * time.Second
		}
		output, errRun := w.dockerClient.RunTest(ctx, containerID, []string{"python3", "/app/main.py"}, test.InputData, maxExecTime)
		switch {
		case errRun == nil:
			if output.Stdout != test.OutputData {
				return daenerys.SubmissionStatus_STATUS_WRONG_ANSWER, nil
			}
		case errors.Is(errRun, docker.RuntimeError):
			return daenerys.SubmissionStatus_STATUS_RUNTIME_ERROR, nil
		case errors.Is(errRun, docker.TimeLimitExceeded):
			return daenerys.SubmissionStatus_STATUS_TIME_LIMIT_EXCEEDED, nil
		case errors.Is(errRun, docker.MemoryLimitExceeded):
			return daenerys.SubmissionStatus_STATUS_MEMORY_LIMIT_EXCEEDED, nil
		default:
			return daenerys.SubmissionStatus_STATUS_INTERNAL_ERROR, fmt.Errorf("docker.RunTest: %w", errRun)
		}
	}

	return daenerys.SubmissionStatus_STATUS_OK, nil
}
