package mr

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"

	"6.5840/mr/logger"
)

type Coordinator struct {
	// Your definitions here.
	files map[string]bool
	nReduce int
}

// Your code here -- RPC handlers for the worker to call.

// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	logger.Log.Info("COORDINATOR: Called example coordinator function")
	reply.Y = args.X + 1
	return nil
}

func (c *Coordinator) GiveWork(args *AskForWorkArgs, reply *AskForWorkReply) error {
	logger.Log.Infof("COORDINATOR: Worker #%v called 'AskForWork'", args.WorkerId)

	reply.WorkType = "map" // TODO: use something more sophisticated here
	reply.Filename = "pg-testfile.txt"
	reply.ReduceTasks = c.nReduce

	return nil
}

func (c *Coordinator) DoneWork(args *DoneWorkArgs, reply *DoneWorkReply) error {
	logger.Log.Infof("COORDINATOR: Worker #%v called 'DoneWork'", args.WorkerId)

	// Mark file as DONE
	c.files[args.Filename] = true
	logger.Log.Infof("COORDINATOR: File %s called MAP successfully", args.Filename)

	return nil
}

// start a thread that listens for RPCs from worker.go
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	logger.Log.Info("COORDINATOR: Coordinator starts listening on socket:", sockname)
	go http.Serve(l, nil)
}

// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
func (c *Coordinator) Done() bool {
	for f, _ := range c.files {
		if !c.files[f] {
			logger.Log.Info("COORDINATOR: Done() has been called to check whether the coordinator is done. Right now the value is:", false)
			return false
		}
	}
	// Your code here.

	logger.Log.Info("COORDINATOR: Done() has been called to check whether the coordinator is done. Right now the value is:", true)
	return true
}

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	fileMap := make(map[string]bool, len(files))
	for _, f := range files {
		fileMap[f] = false
	}

	c := Coordinator{fileMap,nReduce}

	logger.Log.Info("COORDINATOR: Created coordinator...")

	// Your code here.
	// TODO: Divide intermediate keys in nReduce buckets -> Each mapper should create nReduce intermediate files for consumption by the reduce tasks.
	logger.Log.Info("Files are: ", files)

	c.server()
	return &c
}
