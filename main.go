package main

import "magmar/http"

func main() {
	e := http.EchoHandler()
	e.Logger.Fatal(e.Start(":33333"))
}
