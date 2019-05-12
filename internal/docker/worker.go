package docker

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"io"
	"time"
)

// Worker store all infomations about the run job
type Worker struct {
	Image       string // images name
	containerID string
	cli         *client.Client
	ctx         context.Context
	// ric's stdin stdout stderr
	ricIn  io.Reader
	ricOut bytes.Buffer
	ricErr bytes.Buffer
}

// Run start a worker
func (w *Worker) Run() (string, string, error) {

	err := w.attachContainer()
	if err != nil && w.ricErr.Len() == 0 {
		return "", "", err
	}

	return w.ricOut.String(), w.ricErr.String(), nil
}

func (w *Worker) attachContainer() (err error) {
	options := types.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	}

	log.Info("container ", w.containerID, " Attaching...")
	hijacked, err := w.cli.ContainerAttach(w.ctx, w.containerID, options)
	if err != nil {
		return
	}
	defer hijacked.Close()

	log.Info("container ", w.containerID, " Starting ...")
	err = w.cli.ContainerStart(w.ctx, w.containerID, types.ContainerStartOptions{})
	if err != nil {
		return
	}

	log.Info("container ", w.containerID, " Waiting for attach to finish...")
	attachCh := make(chan error, 2)

	// Copy any output to the build trace
	go func() {
		oc, err := stdcopy.StdCopy(&w.ricOut, &w.ricErr, hijacked.Reader)
		if oc > MaxOutInBytes {
			attachCh <- ErrTooMuchOutPut
		}
		if err != nil {
			attachCh <- err
		}
	}()

	// Write the input to the container and close its STDIN
	go func() {
		_, err := io.Copy(hijacked.Conn, w.ricIn)
		hijacked.CloseWrite()
		if err != nil {
			attachCh <- err
		}
	}()

	waitCh := make(chan error, 1)
	go func() {
		waitCh <- w.waitForContainer()
	}()

	select {
	case <-w.ctx.Done():
		err = w.killContainer(w.containerID, waitCh)
		log.Error(err)
		err = errors.New("Aborted")

	case err = <-attachCh:
		errk := w.killContainer(w.containerID, waitCh)
		if errk != nil {
			log.Error(errk)
		}

		log.Info("container ", w.containerID, " attach finished with ", err)

	case err = <-waitCh:
		log.Info("container ", w.containerID, " wait finished with ", err)

	case <-time.After(10 * time.Second):
		errk := w.killContainer(w.containerID, waitCh)
		if errk != nil {
			log.Error(errk)
		}

		err = ErrWorkerTimeOut
		log.Info("container ", w.containerID, " time out")
	}
	return
}

func (w *Worker) waitForContainer() error {
	log.Info("container ", w.containerID, " Waiting...")

	retries := 0
	// Use active wait
	for {
		container, err := w.cli.ContainerInspect(w.ctx, w.containerID)
		if err != nil {
			log.Info(err.Error())
			if client.IsErrNotFound(err) {
				return err
			}

			if retries > 6 {
				return err
			}

			retries++
			time.Sleep(500 * time.Millisecond)
			continue
		}

		// Reset retry timer
		retries = 0
		if container.State.Running {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		if container.State.ExitCode != 0 {
			return fmt.Errorf("exit code %d", container.State.ExitCode)
		}

		return nil
	}
}

func (w *Worker) killContainer(id string, waitCh chan error) error {
	for {
		log.Info("container ", id, " Killing ...")
		err := w.cli.ContainerKill(w.ctx, id, "SIGKILL")
		if err != nil {
			log.Error(err)
			return err
		}
		// Wait for signal that container were killed
		// or retry after some time
		select {
		case err = <-waitCh:
			return err

		case <-time.After(time.Second):
		}
	}
}
