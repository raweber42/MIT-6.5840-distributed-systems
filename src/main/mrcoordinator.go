package main

//
// start the coordinator process, which is implemented
// in ../mr/coordinator.go
//
// go run mrcoordinator.go pg*.txt
//
// Please do not change this file.
//

import (
	"fmt"
	"os"
	"time"

	"6.5840/mr"
	"6.5840/mr/logger"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: mrcoordinator inputfiles...\n")
		os.Exit(1)
	}

	logger.Log.Info("$$$$$$$$$$$$$$$$$$ COORDINATOR STARTED $$$$$$$$$$$$$$$$$$")
	m := mr.MakeCoordinator(os.Args[1:], 10)
	for m.Done() == false {
		time.Sleep(time.Second)
	}

	logger.Log.Info("$$$$$$$$$$$$$$$$$$ COORDINATOR FINISHED $$$$$$$$$$$$$$$$$$")

	time.Sleep(time.Second)
}
