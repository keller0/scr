package docker

import (
	"context"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/keller0/scr/internal/env"
	log "github.com/sirupsen/logrus"
	"io"
	"strconv"
	"time"
)

var (
	MaxOutInBytes    int64 = 2 * 1024 * 1024
	ErrTooMuchOutPut       = errors.New("too much out put")
	ErrWorkerTimeOut       = errors.New("time out")
	memLimit               = env.Get("CONTAINER_MEM_LIMIT", "50")
	diskLimit              = env.Get("CONTAINER_DISK_LIMIT", "5")
	pidLimit               = env.Get("CONTAINER_DISK_LIMIT", "50")
)

var (
	GccWorker  chan string
	GoWorker   chan string
	QuitSignal chan int
)

type Job struct {
	Image   string // images name
	Payload io.Reader
}

func init() {
	log.Info("manager init")
	GccWorker = make(chan string, 2)
	GoWorker = make(chan string, 2)
	QuitSignal = make(chan int)
}

func (jb *Job) Do() (string, string, error) {

	work := new(Worker)
	containerID, err := getContainerByName(jb.Image)
	if err != nil {
		return "", "", err
	}
	log.Info("got container: ", containerID)
	work.containerID = containerID
	work.cli, err = client.NewEnvClient()
	if err != nil {
		return "", "", err
	}
	work.ricIn = jb.Payload
	work.ctx = context.Background()
	defer removeContainer(containerID)

	return work.Run()

}

func StartManagers() {

	log.Info("starting manager")
	go startWorkers(GccWorker, "yximages/gcc:8.3", QuitSignal)
	go startWorkers(GoWorker, "yximages/golang:1.12", QuitSignal)

}

func startWorkers(ws chan string, image string, q chan int) {
	for {
		select {
		case <-q:
			close(ws)
			log.Info("stopping workers")
			return

		default:
			log.Debug("starting a ", image)
			cId, err := CreateContainer(image)
			if err != nil {
				log.Error("create container failed")
			}
			ws <- cId
		}
	}

}

func removeContainer(cid string) {
	log.Info("container ", cid, " removing...")

	cli, err := client.NewEnvClient()
	if err != nil {
		log.Error(err)
	}
	err = cli.ContainerRemove(context.Background(), cid, types.ContainerRemoveOptions{})
	if err != nil {
		log.Error(err)
	}
	log.Info("container ", cid, " removed")

}

// JobStop stop all jobs
func JobStop() {
	log.Info("start to stop all jobs")
	close(QuitSignal)
	time.Sleep(2 * time.Second)

	log.Info("start remove all containers")
	for len(GccWorker) > 0 {
		removeContainer(<-GccWorker)
	}
	for len(GoWorker) > 0 {
		removeContainer(<-GoWorker)
	}
	log.Info("all job stopped")
}

func getContainerByName(image string) (string, error) {
	log.Info("try get a container of:", image)
	switch image {
	case "yximages/gcc:10":
		return <-GccWorker, nil
	case "yximages/golang:1.14":
		return <-GoWorker, nil
	default:
		return CreateContainer(image)
	}
}

func CreateContainer(image string) (string, error) {

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return "", err
	}

	config := &container.Config{
		Image:        image,
		Cmd:          []string{"/home/ric/run"},
		AttachStdin:  true, // Attach the standard input, makes possible user interaction
		AttachStdout: true, // Attach the standard output
		AttachStderr: true,
		Tty:          false,
		OpenStdin:    true,
		StdinOnce:    true,
	}
	var diskLimit64 int64
	ml, _ := strconv.Atoi(memLimit)
	dl, _ := strconv.Atoi(diskLimit)
	pl, _ := strconv.Atoi(pidLimit)
	diskLimit64 = int64(dl)
	hostConfig := &container.HostConfig{
		Resources: container.Resources{
			CPUPeriod: 100000,
			CPUQuota:  50000,
			Memory:    int64(ml) * 1024 * 1024,
			PidsLimit: int64(pl),
			// TODO: put this resources to config
			// advanced kernel-level features
			// CPURealtimePeriod : 1000000,
			// CPURealtimeRuntime: 950000,

			DiskQuota: diskLimit64 * 1024 * 1024,
		},
		Privileged: false,
		LogConfig: container.LogConfig{
			Type: "json-file",
		},
	}
	var tmpId string
	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, "")
	if err != nil {
		return "", err
	}
	tmpId = resp.ID
	_, err = cli.ContainerInspect(ctx, tmpId)
	if err != nil {
		return "", err
	}

	return tmpId, nil
}
