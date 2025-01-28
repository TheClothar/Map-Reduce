package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//rpc.go

import (
	"os"
	"strconv"
)

// example to show how to declare the arguments
// and reply for an RPC.
type WorkerArgs struct {
	WorkerID int
}

type RequestReply struct {
	Nmap      int
	FileName  string
	NReduce   int
	TaskType  TaskType
	TaskID    int
	TaskStat  TaskStatus
	Terminate bool
	NMap      int
	NumMap    int
}

type ReportDoneArgs struct {
	WorkerID int
	TaskID   int
	TaskType TaskType
}

type DoneReply struct {
	Terminate bool
}

// Add your RPC definitions here.

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
