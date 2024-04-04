package server

// Reqs key: requestID value: result
var Reqs = make(map[string]chan string)
