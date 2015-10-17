package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/coocood/freecache"
)

type HTTPResponse struct {
	Status int
	Key    string
	Value  string
}

type JSONData struct {
	Key   string
	Value string
}

////////////
// Worker //
////////////
// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.
func NewWorker(id int, workerQueue chan chan WorkRequest, cache *freecache.Cache) Worker {
	// Create, and return the worker.
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
		Cache:       cache}

	return worker
}

type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
	Cache       *freecache.Cache
}

// This function "starts" the worker by starting a goroutine, that is
// an infinite "for-select" loop.
func (w Worker) Start() {
	go func() {
		var Value string
		var R HTTPResponse

		for {
			// Add ourselves into the worker queue.
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				// fmt.Println(work)
				t0 := time.Now()

				fmt.Printf("Work being dispatcher by worker number:%v\n", w.ID)
				valuecache, err := w.Cache.Get([]byte(work.Key))

				// Not in the cache
				if err == freecache.ErrNotFound {
					fmt.Printf("Key:%v request is NOT in the cache\n", work.Key)
					var d []JSONData

					file, err := ioutil.ReadFile("./database.json")
					if err != nil {
						fmt.Printf("File error: %v\n", err)
						work.done <- false
						continue
					}

					err = json.Unmarshal(file, &d)
					if err != nil {
						fmt.Printf("Error while decoding JSON: %v\n", err)
						work.done <- false
						continue
					}

					found := false
					for _, v := range d {
						if bytes.Equal([]byte(v.Key), []byte(work.Key)) {
							fmt.Printf("Key:%v request is in the database.\n", work.Key)
							found = true
							Value = v.Value
							w.Cache.Set([]byte(work.Key), []byte(Value), 0)
						}
					}

					if found {
						R = HTTPResponse{Status: 200, Key: work.Key, Value: Value}
					} else {
						fmt.Printf("Key:%v request is NOT in the database.\n", work.Key)
						R = HTTPResponse{Status: 204, Key: work.Key, Value: "Null"}
					}

				} else { // In the cache
					fmt.Printf("Key:%v request is in the cache\n", work.Key)
					// n := bytes.IndexByte(valuecache, 0)
					s := string(valuecache[:])
					R = HTTPResponse{Status: 207, Key: work.Key, Value: s}
				}

				if testMode {
					time.Sleep(5 * time.Millisecond)
				}

				js, err := json.Marshal(R)

				// fmt.Println(R)
				// fmt.Println(js)
				if err != nil {
					http.Error(work.ConnResponse, err.Error(), http.StatusInternalServerError)
					fmt.Println("Error during JSON Formating")
					return
				}

				work.ConnResponse.Header().Set("Content-Type", "application/json")
				work.ConnResponse.Write(js)
				work.done <- true
				fmt.Printf("Work done!(from worker %v). \n", w.ID)
				t1 := time.Now()
				fmt.Printf("Processing time:%v\n", t1.Sub(t0))
			case <-w.QuitChan:
				// We have been asked to stop.
				fmt.Printf("worker%d stopping\n", w.ID)
				return
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
// Note that the worker will only stop *after* it has finished its work.
func (w Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}
