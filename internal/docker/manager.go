package docker

import (
	"context"
	"errors"
	"fmt"
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
	ErrTooMuchOutPut       = errors.New("Too much out put")
	ErrWorkerTimeOut       = errors.New("Time out")
	memLimit               = env.Get("CONTAINER_MEM_LIMIT", "50")
	diskLimit              = env.Get("CONTAINER_DISK_LIMIT", "5")
)

var (
	doneWorkers chan string
	gccWorker   chan string
	goWorker    chan string
)

type Job struct {
	Image   string // images name
	Payload io.Reader
}

func (jb *Job) Do() (string, string, error) {

	work := new(Worker)

	// TODO get a container instead create one
	containerID, err := getWorkerByName(jb.Image)
	if err != nil {
		return "", "", err
	}
	work.containerID = containerID
	work.cli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", "", err
	}
	work.ricIn = jb.Payload
	work.ctx = context.Background()
	defer func() {
		doneWorkers <- containerID
	}()

	return work.Run()

}

func StartManagers() {

	log.Info("starting manager")
	go startWorkers(gccWorker, "yximages/gcc:8.3", 10)
	go startWorkers(goWorker, "yximages/golang:1.12", 10)
	go WorkersGo(doneWorkers)

}

func startWorkers(ws chan string, image string, num int) {
	ws = make(chan string, num)
	for {
		log.Debug("starting a ", image)
		cId, err := CreateContainer(image)
		if err != nil {
			log.Error("create contianer failed")
		}
		ws <- cId
	}

}

// WorkersGo remove containers
func WorkersGo(wsg chan string) {
	for {
		ws := <-wsg
		fmt.Println(ws, "removing...")
		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			panic(err)
		}
		err = cli.ContainerRemove(context.Background(), ws, types.ContainerRemoveOptions{})
		if err != nil {
			panic(err)
		}
		fmt.Println("removed container:", ws)

		time.Sleep(3 * time.Second)
		fmt.Println(ws, "removed")
	}
}

func getWorkerByName(image string) (string, error) {
	switch image {

	case "yximages/gcc:8.3":
		return <-gccWorker, nil
	case "yximages/golang:1.12":
		return <-goWorker, nil
	default:
		return CreateContainer(image)
	}
}

func CreateContainer(image string) (string, error) {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
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
	ml, _ := strconv.Atoi(memLimit)
	dl, _ := strconv.Atoi(diskLimit)
	hostConfig := &container.HostConfig{
		Resources: container.Resources{
			CPUPeriod: 100000,
			CPUQuota:  50000,
			Memory:    int64(ml) * 1024 * 1024,
			PidsLimit: 50,
			// advanced kernel-level features
			// CPURealtimePeriod : 1000000,
			// CPURealtimeRuntime: 950000,

			DiskQuota: int64(dl) * 1024 * 1024,
		},
		Privileged: false,
		LogConfig: container.LogConfig{
			Type: "json-file",
		},
	}
	var tmpId string
	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, "")
	if err != nil {
		if resp.ID != "" {
			tmpId = resp.ID
		}
		return "", err
	}
	tmpId = resp.ID
	_, err = cli.ContainerInspect(ctx, tmpId)
	if err != nil {
		return "", err
	}

	return tmpId, nil
}
