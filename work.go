package main

import (
	"net/http"
)

type WorkRequest struct {
	Key          string
	ConnResponse http.ResponseWriter
	done         chan bool
}
