package main

import (
	"github.com/yaphper/WebCamera/app"
	"log"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Fatal(app.StartServer())
}
