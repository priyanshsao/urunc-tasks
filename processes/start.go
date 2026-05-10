package process

import (
	"fmt"
	"log"
	"os/exec"
	"syscall"
	"time"
)

func StartProcess() (*exec.Cmd, chan struct{}) {

	cmd := exec.Command("sleep", "2")

	err := cmd.Start()
	if err != nil {
		log.Fatal("Error starting process:", err)
	}

	fmt.Println("Child process started with PID:", cmd.Process.Pid)

	done := make(chan struct{})

	// Reap child process
	go func() {

		defer close(done)

		err := cmd.Wait()

		if err != nil {
			fmt.Println("Process exited:", err)
		} else {
			fmt.Println("Process exited successfully")
		}
	}()

	return cmd, done
}

func StopProcess(cmd *exec.Cmd) {

	fmt.Println("Waiting 1 second before sending SIGTERM...")
	time.Sleep(1 * time.Second)

	err := cmd.Process.Signal(syscall.SIGTERM)
	if err != nil {
		log.Fatal("Error sending signal:", err)
	}

	fmt.Println("SIGTERM sent")
}
