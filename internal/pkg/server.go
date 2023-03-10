package pkg

type DefaultServerInterface interface {
	Init() DefaultServerInterface
	Run()
}
