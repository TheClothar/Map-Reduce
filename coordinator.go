package mr

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"time"
)

type TaskType string

const (
	MapTask    TaskType = "map"
	ReduceTask TaskType = "reduce"
)

type TaskStatus string

const (
	NotStarted TaskStatus = "not started"
	InProgress TaskStatus = "in progress"
	Finished   TaskStatus = "completed"
)

type Task struct {
	ID       int
	File     string
	Status   TaskStatus
	Type     TaskType
	WorkerID int
}

type Coordinator struct {
	Safe     sync.Mutex
	TaskList []Task
	Phase    TaskType
	NReduce  int
	NMap     int
}

// handle worker requests for work
func (c *Coordinator) TaskRequestHandler(args *WorkerArgs, reply *RequestReply) error {
	c.Safe.Lock()
	defer c.Safe.Unlock()

	// Handle phase transition from map to reduce.
	if c.Phase == MapTask && c.allTasksDone(MapTask) {
		c.transitionToReduce()
	}

	// Assign a task based on current phase and task status.
	for i := range c.TaskList {
		task := &c.TaskList[i]
		if task.Status == NotStarted && task.Type == c.Phase {
			c.assignTask(task, reply, args.WorkerID)
			return nil
		}
	}

	// If no task is available and all reduce tasks are done, signal termination.
	if c.Phase == ReduceTask && c.allTasksDone(ReduceTask) {
		reply.Terminate = true
	}

	return nil
}

// Helper function to check if all tasks of a given type are done.
func (c *Coordinator) allTasksDone(taskType TaskType) bool {
	for _, task := range c.TaskList {
		if task.Type == taskType && task.Status != Finished {
			return false
		}
	}
	return true
}

// Transition to reduce phase
func (c *Coordinator) transitionToReduce() {
	c.Phase = ReduceTask
	log.Println("Switching to reduce phase.")
	for i := range c.TaskList {
		c.TaskList[i].Type = ReduceTask
		c.TaskList[i].Status = NotStarted
	}
}

// Assigns a task to a worker and starts a timeout timer.
func (c *Coordinator) assignTask(task *Task, reply *RequestReply, workerID int) {
	task.Status = InProgress
	task.WorkerID = workerID
	reply.FileName = task.File
	reply.TaskType = task.Type
	reply.TaskID = task.ID
	reply.NReduce = c.NReduce
	reply.NumMap = c.NMap
	go c.tenSecTimer(task)
}

// Timeout handler for tasks.
func (c *Coordinator) tenSecTimer(task *Task) {
	time.AfterFunc(time.Second*10, func() {
		c.Safe.Lock()
		defer c.Safe.Unlock()
		if task.Status == InProgress {
			//log.Printf("Task %d timeout, resetting status.\n", task.ID)
			task.Status = NotStarted
			task.WorkerID = -1
		}
	})
}

// Mark task as done on completion.
func (c *Coordinator) DoneHandler(args *ReportDoneArgs, reply *DoneReply) error {
	c.Safe.Lock()
	defer c.Safe.Unlock()
	for i, task := range c.TaskList {
		if task.ID == args.TaskID && task.WorkerID == args.WorkerID && task.Status == InProgress {
			c.TaskList[i].Status = Finished
			break
		}
	}
	return nil
}

// Checks if all tasks are completed to determine if the job is done.
func (c *Coordinator) Done() bool {
	c.Safe.Lock()
	defer c.Safe.Unlock()
	return c.allTasksDone(MapTask) && c.allTasksDone(ReduceTask)
}

// Coordinator initialization.
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{
		Phase:    MapTask,
		TaskList: make([]Task, len(files)),
		NReduce:  nReduce,
		NMap:     len(files),
	}

	for i, file := range files {
		c.TaskList[i] = Task{
			ID:       i,
			File:     file,
			Status:   NotStarted,
			Type:     MapTask,
			WorkerID: -1,
		}
	}

	c.server()
	return &c
}

// RPC server setup.
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}
