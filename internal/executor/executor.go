package executor

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

func Run(ctx context.Context, command string, timeout time.Duration) (string, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	cmd := exec.CommandContext(ctxTimeout, "bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctxTimeout.Err() == context.DeadlineExceeded {
			return string(output), fmt.Errorf("command timed out (%v) : %w", err, ctxTimeout.Err())
		}
		return string(output), fmt.Errorf("command failed : %w", err)

	}
	return string(output), nil
}
