package init

import (
	appPing "github.com/siyoga/rollstory/internal/app/ping"
	"github.com/siyoga/rollstory/internal/rpc"
	rpcPing "github.com/siyoga/rollstory/internal/rpc/ping"
	"github.com/siyoga/rollstory/pkg/container"
)

// provideRPC registers all RPC layer (network layer) handlers
func provideRPC(c *container.DigContainer) {
	// Register shared error handler
	c.Provide(api.NewErrorHandler)

	// Register ping RPC handler with interface bindings
	c.Provide(rpcPing.NewHandler).
		// Bind app layer handler to RPC layer interface
		Bind(new(appPing.Handler), new(rpcPing.PingHandler)).
		// Bind error handler to RPC layer interface
		Bind(new(api.ErrorHandler), new(rpcPing.ErrorHandler))

	// Future RPC handlers will be registered here with their bindings:
	// c.Provide(rpcUser.NewHandler).
	//   Bind(new(appUser.Handler), new(rpcUser.UserHandler)).
	//   Bind(new(rpc.ErrorHandler), new(rpcUser.ErrorHandler))
}
