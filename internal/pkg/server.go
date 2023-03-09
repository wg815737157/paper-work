package pkg

type DefaultServer interface {
	Init() DefaultServer
	Run()
}
