package runner

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/skyterra/elastic-embedding-searcher/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

type ModelXManager struct {
	cmd     *exec.Cmd
	cmdlock sync.Mutex

	client  pb.ModelxClient
	clilock sync.RWMutex

	stop chan struct{}

	workers   int
	cmdPath   string
	modelPath string
}

// fork starts the ModelX process and waits for readiness.
func (m *ModelXManager) fork() error {
	python := "python"

	// fork modelx process. it is too slow...
	cmd := exec.Command(python, m.cmdPath, "--max_workers", strconv.Itoa(m.workers), "--uds", "True", "--model_path", m.modelPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	// wait for modelx process ready.
	sigReady := make(chan os.Signal, 1)
	signal.Notify(sigReady, syscall.SIGUSR1)
	defer signal.Stop(sigReady)

	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	select {
	case <-sigReady:
		m.setcmd(cmd)
		return nil

	case <-ticker.C:
		return errors.New("fork ModelX timeout")
	}
}

// kill stops the ModelX process if it's running.
func (m *ModelXManager) kill() error {
	m.cmdlock.Lock()
	defer m.cmdlock.Unlock()

	if m.cmd == nil || m.cmd.Process == nil {
		return nil
	}

	// check to see if process is running.
	// err is NOT nil means that process NOT exist.
	if err := m.cmd.Process.Signal(syscall.Signal(0)); err != nil {
		return nil
	}

	// kill it. (kill -9 process_id)
	if err := m.cmd.Process.Kill(); err != nil {
		return err
	}

	// clear sub-process <defunct>
	m.cmd.Process.Wait()

	m.cmd = nil
	return nil
}

// dial connects to the ModelX service.
func (m *ModelXManager) dial() error {
	conn, err := grpc.NewClient(uds,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(10<<20),
			grpc.MaxCallSendMsgSize(10<<20),
		),
	)
	if err != nil {
		return err
	}

	c := pb.NewModelxClient(conn)
	if err = m.ping(c); err != nil {
		return err
	}

	// maybe the following code looks too stupid,
	// BUT it is required, the modelx process maybe restart at any time.
	m.clilock.Lock()
	m.client = c
	m.clilock.Unlock()
	return nil
}

// monitor periodically checks the ModelX service health.
func (m *ModelXManager) monitor() error {
	ticker := time.NewTicker(time.Second)
	retry := 0

	go func() {
		for {
			select {
			case <-ticker.C:
				m.check(&retry)
			case <-m.stop:
				return
			}
		}
	}()

	return nil
}

// check performs a health check and restarts ModelX if needed.
func (m *ModelXManager) check(retry *int) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic: err:%v\n", err)
		}
	}()

	// perform a health check (ping) on the ModelX client.
	if err := m.ping(m.client); err == nil {
		return
	}

	*retry++

	// If the retry count is not a multiple of 3, return without restarting.
	// This means restart will only be attempted every third failure.
	if *retry%3 != 0 {
		return
	}

	// restart the ModelX service.
	if err := m.restart(); err != nil {
		log.Printf("restart ModelX failed. err:%s\n", err.Error())
		return
	}

	log.Println("restart ModelX succeed.")
}

// restart stops and restarts ModelX and reconnects.
func (m *ModelXManager) restart() error {
	if err := m.kill(); err != nil {
		return err
	}

	if err := m.fork(); err != nil {
		return err
	}

	return m.dial()
}

// setcmd updates the ModelX process command.
func (m *ModelXManager) setcmd(cmd *exec.Cmd) {
	m.cmdlock.Lock()
	defer m.cmdlock.Unlock()

	m.cmd = cmd
}

// ping nothing to say.
func (m *ModelXManager) ping(client pb.ModelxClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	reply, err := client.Ping(ctx, &pb.PingRequest{})
	if err != nil {
		return err
	}

	if reply.Code != 0 {
		return errors.New(fmt.Sprintf("ping failed. code:%d", reply.Code))
	}

	return nil
}
