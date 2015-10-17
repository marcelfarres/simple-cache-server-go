package main

import (
	"fmt"
	"net/http"
	// "reflect"
)

///////////////
// Collector //
///////////////

// A buffered channel that we can send work requests on.
var WorkQueue = make(chan WorkRequest, 512)

func Collector(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Now, we retrieve the key from the request.
	key := r.FormValue("key")

	//fmt.Println(reflect.TypeOf(key))
	// Just do a quick bit of sanity checking to make sure the client actually provided us with a key.
	if key == "" {
		http.Error(w, "You must specify a key.", http.StatusBadRequest)
		return
	}

	dch := make(chan bool)
	// Now, we take key and make a WorkRequest.
	work := WorkRequest{Key: key, ConnResponse: w, done: dch}
	// fmt.Println(work)

	// Push the work onto the queue.
	WorkQueue <- work
	fmt.Println("Work request queued")

	// And let the user know their work request was created.
	w.WriteHeader(http.StatusCreated)
	succ := <-dch
	fmt.Printf("Done signal recived. Work request finished is %v \n\n", succ)
	return
}
