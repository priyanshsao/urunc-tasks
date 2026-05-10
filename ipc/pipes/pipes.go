package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

const binaryPath = "/proc/self/exe"

func main() {
	Start()
}

func Start() {
	if os.Getenv("REXEC") == "1" {
		reExec()
	} else {
		process()
	}
}

func process() {
	fmt.Println("[Parent] PID: ", os.Getpid())

	rP, wP, err := os.Pipe()
	if err != nil {
		log.Fatal("[Parent] Err generating pipe: ", err)
	}

	cmd := exec.Command(binaryPath)
	cmd.Env = append(os.Environ(), "REXEC=1")

	// pass read pipe to child
	cmd.ExtraFiles = []*os.File{rP}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	//start child process
	if err = cmd.Start(); err != nil {
		panic(err)
	}

	rP.Close()

	message := "ACK_FROM_PARENT"

	fmt.Println("[Parent] sending: ", message)

	if _, err = wP.Write([]byte(message)); err != nil {
		panic(err)
	}

	wP.Close()

	// wait for child to return
	cmd.Wait()
	fmt.Println("[Parent] exited with exit code(0)")
}

func reExec() {

	fmt.Println("[Child] PID:", os.Getpid())

	rP := os.NewFile(uintptr(3), "pipe")

	fmt.Println("[Child] waiting for message...")

	buf := make([]byte, 100)

	n, err := rP.Read(buf)
	if err != nil {
		panic(err)
	}

	fmt.Println("[Child] Recieved: ", string(buf[:n]))
	fmt.Println("[Child] exiting...")
}
