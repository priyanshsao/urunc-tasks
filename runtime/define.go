package main

type containerStatus string

const (
	containerCreating containerStatus = "creating"
	containerRunning containerStatus = "running"
	containerPaused containerStatus = "paused"
)

// container state
type State struct {
	Pid      int             `json:"pid"`
	Status   containerStatus `json:"status"`
	HostName string          `json:"hostname"`
}

type NsConf struct {
	HostName string `json:"hostname"`
}