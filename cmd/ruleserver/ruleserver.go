package main

import "github.com/wg815737157/paper-work/internal/ruleserver"

func main() {
	ruleServer := ruleserver.DefaultServer()
	ruleServer.Init().Run()
}
