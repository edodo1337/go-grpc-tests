gRPC k-v storage with LRU-cache server & client + load test.

* Specify LRU-cache capacity: server/cmd/app.go -> StorageCapacity
* go build ./server/cmd/* && ./app

* Specify load test consts in load_test.go: CLIENTS_COUNT, KEY_RANGE, LOOPS_COUNT
* Run TestLoad


For educational purposes.

