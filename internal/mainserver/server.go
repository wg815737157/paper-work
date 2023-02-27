package mainserver

import (
	"deps/log4go"
	"github.com/gin-gonic/gin"
	internalpkg "github.com/wg815737157/paper-work/internal/pkg"
	"log"
)

func Run() {
	log4go.Info("d")
	tree := &internalpkg.Tree{}
	node := &internalpkg.Node{}
	tree.Nodes = append(tree.Nodes, node)
	ginEngine := gin.New()
	ginEngine.Handler()
	err := ginEngine.Run("localhost:3747")
	if err != nil {
		log.Fatalln(err)
	}
}
