package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"

	"cinderwolf.net/raycaster/engine"
	"github.com/CAFxX/gcnotifier"
)

// main test
func main() {
	runtime.SetCPUProfileRate(10000)

	gcn := gcnotifier.New()
	go func() {
		for range gcn.AfterGC() {
			fmt.Println("gc")
		}
	}()
	defer gcn.Close()

	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	engine := engine.NewEngine()
	engine.Start()
}
