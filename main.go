package main

import "net/http"

func main() {
	httpQ := &HTTPQ{}
	http.ListenAndServe(":24744", httpQ.Handler())
}
