package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"syscall"
)

const selfBinaryPath string = "/proc/self/exe" 
const profsPath = "/tmp/proc"
const shPath = "/bin/sh"
// reexec env
const reexecEnv = "REEXEC"

func main() {
	
	if os.Getenv(reexecEnv) == "1" {
		reexecProcess()
		return
	}

	process()
}

func process() {
	
	fmt.Println("[Parent] PID: ", os.Getpid())

	parentSock, childSock := newSocket()

	state := new(State)

	state.Pid = -1
	state.Status = containerCreating
	state.HostName = "host-container"

	saveState(state)
	fmt.Println("creating state.json")

	cmd := createReExecProcess(childSock)

	childSock.Close()

	fmt.Println("Sending namespace config to child...")

	nsConf := new(NsConf)
	nsConf.HostName = "urunc-demo"
	
	// create encoder to directly send json data to reexec process
	encoder := json.NewEncoder(parentSock)
	
	if err := encoder.Encode(nsConf); err != nil {
		panic(err)
	}

	buf := make([]byte, 100)

	n, err := parentSock.Read(buf)
	if err != nil {
		panic(err)
	}

	fmt.Println("[Parent] Recieved:", string(buf[:n]))

	//update state.json
	state.Pid = cmd.Process.Pid
	state.Status = containerRunning

	saveState(state)

	// wait for child to finish
	cmd.Wait()
	fmt.Println("[Parent] exited with code(0)")
}


func reexecProcess() {

	fmt.Println("[Reexec] PID:", os.Getpid())

	// open fd
	socket := os.NewFile(uintptr(3), "socket")

	conn, err := net.FileConn(socket)
	if err != nil {
		panic(err)
	}

	var nsConf NsConf

	decoder := json.NewDecoder(conn)

	if err := decoder.Decode(&nsConf); err != nil {
		panic(err)
	}

	fmt.Println("[Reexec] Received hostname: ", nsConf.HostName)

	// setting host name
	if err := syscall.Sethostname([]byte(nsConf.HostName)); err != nil {
		panic(err)
	}

	// mount procfs
	os.MkdirAll(profsPath, 0755)

	if err := syscall.Mount("proc", profsPath, "proc", 0, ""); err != nil {
		panic(err)
	}

	fmt.Println("Mounted procfs inside mount namespace")

	// send ack
	conn.Write([]byte("OK"))

	ps := exec.Command("ps", "aux")

	ps.Stdout = os.Stdout
	ps.Stderr = os.Stderr

	ps.Run()

	fmt.Println("Execve into /bin/sh...")

	if err := syscall.Exec(shPath, []string{shPath}, os.Environ()); err != nil {
		panic(err)
	}
}

func createReExecProcess(sock *os.File) *exec.Cmd {

	cmd := exec.Command(selfBinaryPath)

	cmd.Env = append(os.Environ(), reexecEnv+"=1")

	// Pass socket FD to child as FD 3
	cmd.ExtraFiles = []*os.File{sock}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	createNs(cmd)

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	fmt.Println("[Parent] Child created with PID:", cmd.Process.Pid)

	return cmd
}

func createNs(cmd *exec.Cmd) {
	
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWNS,
	}

	fmt.Println("namespace created")
}

func saveState(state *State) {

	// create or rewrite the state.json file
	file, err := os.Create("state.json")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(state); err != nil {
		panic(err)
	}
}