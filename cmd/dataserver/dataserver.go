package main

import (
	"github.com/wg815737157/paper-work/internal/dataserver"
)

func main() {
	dataServer := dataserver.DefaultServer()
	dataServer.Init().Run()
}
