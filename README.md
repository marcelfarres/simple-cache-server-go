# simple-cache-server-go

Simple Go server with cache using worker queues. 

Code inspired by this post from `nesv` --> http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html

Installation
-----------
1. Clone or download the repo.
2. Install freecache using `go get github.com/coocood/freecache`.

Build and Run 
-------------

- build --> `go build -o cacheserver *.go` (respecting GOPATH --> go/src/github.org/{user}/simple-cache-server-go) 

- run --> `./cacheserver -n 20000` (This number is the workers you want)

In order to perform the query test I use this 2 bash scripts: 

- SendR.command 
- SendRmass.command