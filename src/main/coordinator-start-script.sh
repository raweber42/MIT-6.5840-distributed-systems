#/bin/bash
go build -buildmode=plugin ../mrapps/wc.go
rm mr-pg-testfile*
rm mr-out*
go run mrcoordinator.go pg-testfile.txt
# go run mrcoordinator.go pg-*.txt
