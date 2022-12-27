package main

import "net/http"

func main() {
	httpQ := &HTTPQ{}
	http.ListenAndServeTLS(":24744", "localhost.crt", "localhost.key", httpQ.Handler())
}
