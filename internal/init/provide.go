package init

import "github.com/siyoga/rollstory/pkg/container"

func NewContainer() *container.DigContainer {
	di := container.New()

	provideLogger(di)

	provideInf(di)

	// Register app layer (service layer)
	provideApp(di)

	// Register RPC layer (network layer)
	provideRPC(di)

	return di
}
