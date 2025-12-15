package init

import (
	"github.com/siyoga/rollstory/internal/app/ping"
	"github.com/siyoga/rollstory/pkg/container"
)

// provideApp registers all app layer (service layer) handlers
func provideApp(c *container.DigContainer) {
	// Register ping service handler
	c.Provide(ping.NewHandler)

	// Future handlers will be registered here:
	// c.Provide(user.NewHandler)
	// c.Provide(auth.NewHandler)
	// etc.
}
