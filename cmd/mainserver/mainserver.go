package main

import (
	"github.com/wg815737157/paper-work/internal/mainserver"
)

func main() {
	mainServer := mainserver.DefaultServer()
	mainServer.Init().Run()
}
