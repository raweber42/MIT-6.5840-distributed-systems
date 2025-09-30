## Sequential output of pg-testfile.txt
main [fix-add-logs-for-process-flow]$ go run mrsequential.go wc.so pg-testfile.txt
main [fix-add-logs-for-process-flow]$ cat mr-out-0                                
Hello 1
Test 1
file 1
is 1
my 1
small 1
test 3
this 1


---

### Clean coordinator start: `sh coordinator-start-script.sh`
### Worker start: `go run mrworker.go wc.so`

for rpc see here: https://dev.to/atanda0x/a-beginners-guide-to-rpc-in-golang-understanding-the-basics-4eeb

for pretty printing: https://blog.josejg.com/debugging-pretty/
