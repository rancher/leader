package election

import (
	"errors"
	"os"
	"os/exec"
	"syscall"

	"github.com/rancher/go-rancher-metadata/metadata"
)

type Watcher struct {
	leader  metadata.Container
	command []string
	port    int
	client  *metadata.Client
	forward *TcpProxy
}

func New(client *metadata.Client, port int, command []string) *Watcher {
	return &Watcher{
		command: command,
		port:    port,
		client:  client,
	}
}

func (w *Watcher) getLeader() (metadata.Container, bool, error) {
	selfContainer, err := w.client.GetSelfContainer()
	if err != nil {
		return metadata.Container{}, false, err
	}

	index := selfContainer.CreateIndex
	leader := selfContainer

	containers, err := w.client.GetServiceContainers(
		selfContainer.ServiceName,
		selfContainer.StackName,
	)
	if err != nil {
		return metadata.Container{}, false, err
	}

	for _, container := range containers {
		if container.CreateIndex < index {
			index = container.CreateIndex
			leader = container
		}
	}

	w.leader = leader
	return leader, leader.PrimaryIp == selfContainer.PrimaryIp, nil
}

func (w *Watcher) Watch() error {
	w.forward = NewTcpProxy(w.port, func() string {
		return w.leader.PrimaryIp
	})

	go w.client.OnChange(2, func(version string) {
		if w.IsLeader() {
			w.forward.Close()
		}
	})

	if w.port > 0 {
		if err := w.forward.Forward(); err != nil {
			return err
		}
	}

	if w.IsLeader() {
		if len(w.command) == 0 {
			return errors.New("No command")
		}

		prog, err := exec.LookPath(w.command[0])
		if err != nil {
			return err
		}
		return syscall.Exec(prog, w.command, os.Environ())
	}

	return errors.New("Unexpected loop termination")
}

func (w *Watcher) IsLeader() bool {
	_, leader, err := w.getLeader()
	return leader && err == nil
}
