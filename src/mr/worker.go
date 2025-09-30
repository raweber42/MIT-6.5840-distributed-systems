package mr

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"

	"6.5840/mr/logger"
)

// Map functions return a slice of KeyValue.
type KeyValue struct {
	Key   string
	Value string
}

// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

// main/mrworker.go calls this function.
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	// Your worker implementation here.
	logger.Log.Info("WORKER: Created worker...")
	// TODO: Add worker number here

	// TODO: Use this later:
	// fmt.Fprintf(ofile, "%v %v\n", intermediate[i].Key, output)

	// uncomment to send the Example RPC to the coordinator.
	reply, err := AskForWork()
	if err != nil {
		log.Fatalf(err.Error())
	}
	if reply.WorkType == "map" {
		callMapFunction(mapf, reply)
	} else { // Run reduce
		callReduceFunction(reducef, reply)
	}
}

func AskForWork() (AskForWorkReply, error) {
	args := AskForWorkArgs{} // TODO: Add arguments
	workerId := 99999
	args.WorkerId = workerId // TODO: Replace with actual ID
	reply := AskForWorkReply{}

	ok := call("Coordinator.GiveWork", &args, &reply)
	if !ok {
		logger.Log.Info("Coordinator.GiveWork call failed, shutting down!\n")
		return AskForWorkReply{}, errors.New("coordinator.GiveWork call failed") // TODO: Is there a prettier way to return nil?
	}

	return reply, nil
}

func callMapFunction(mapf func(string, string) []KeyValue, reply AskForWorkReply) {
	// TODO: LOCK THIS!!
	file, err := os.Open(reply.Filename)
	if err != nil {
		log.Fatalf("cannot open %v", reply.Filename)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("cannot read %v", reply.Filename)
	}
	file.Close()
	kva := mapf(reply.Filename, string(content))

	// Store intermediate file and tell coordinator, that you're done
	// Each mapper should create nReduce intermediate files for consumption by the reduce tasks.
	for i := 1; i <= reply.ReduceTasks; i++ {
		intermediateFileName := fmt.Sprintf("mr-%v-%v", reply.Filename, ihash(reply.Filename)/i)
		file, err = os.Create(intermediateFileName)
		if err != nil {
			log.Fatalf("cannot open %v", intermediateFileName)
		}
		enc := json.NewEncoder(file)
		for j := i; j < len(kva); j += reply.ReduceTasks { // Evenly distribute among reducetask files
			err := enc.Encode(&kva[j])
			if err != nil {
				log.Fatalf("cannot encode %v to json", kva[j])
			}
		}
		file.Close() // TODO: Is this neccessary?
	}

	// TODO: Is this call neccessary?
	doneWorkReply := DoneWorkReply{}
	ok := call("Coordinator.DoneWork", &DoneWorkArgs{123456789, reply.Filename}, &doneWorkReply)
	if !ok {
		logger.Log.Info("Coordinator.DoneWork call failed!")
	}
}

func callReduceFunction(reducef func(string, []string) string, reply AskForWorkReply) {
	// TODO: Implement me!!!
}

// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()
	logger.Log.Info("WORKER: Calling from worker via RPC...")
	err = c.Call(rpcname, args, reply)
	if err == nil {
		logger.Log.Info("WORKER: Received response via RPC in worker...")
		return true
	}

	logger.Log.Info("WORKER: RPC call from worker failed with error:", err)
	return false
}
