package main

import (
	"fmt"
	"log"
	"os/exec"
	"syscall"
	"time"
)


func main() {

	cmd := exec.Command("sleep", "30")

	if err := cmd.Start(); err != nil {
		log.Fatal("Error starting process:", err)
	}

	fmt.Println("Child process started with PID: ", cmd.Process.Pid)

	fmt.Println("Waiting 1 second before sending kill signal")
	time.Sleep(1 * time.Second)

	err := cmd.Process.Signal(syscall.SIGTERM)
	if err != nil {
		log.Fatal("Error sending signal:", err)
	}
	
	// wait for child process to exit
	if err = cmd.Wait(); err != nil {
		fmt.Println("Process exited successfully")
	} else {
		fmt.Println("some err: ", err)
	}

}