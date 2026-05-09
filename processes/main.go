package main

import "fmt"

func main() {

	cmd, done := StartProcess()
	StopProcess(cmd)

	<-done

	fmt.Println("Parent exiting..")
}