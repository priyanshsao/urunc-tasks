package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	cmd, err := StartServer()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("QEMU PID:", cmd.Process.Pid)

	time.Sleep(2 * time.Second)

	ctx, cancel := context.WithCancel(context.Background())

	qemu, err := Connect()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan interface{})
	go qemu.Read(ctx, done)

	time.Sleep(2 * time.Second)
	fmt.Println("getting capacilities....")

	qemu.Execute(map[string]interface{}{
		"execute": "qmp_capabilities",
	})

	time.Sleep(2 * time.Second)
	fmt.Println("querying status....")

	qemu.Execute(map[string]interface{}{
		"execute": "query-status",
	})

	time.Sleep(time.Second)
	fmt.Println("closing read goroutine....")
	cancel()

	time.Sleep(2 * time.Second)
	fmt.Println("closing socket....")

	qemu.Execute(map[string]interface{}{
		"execute": "quit",
	})

	<-done // wait to close completely
	
	err = cmd.Wait()
	if err != nil {
		log.Println("QEMU exited:", err)
	}

	fmt.Println("successfully exited vmm")
}