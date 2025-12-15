package init

import (
	"github.com/siyoga/rollstory/pkg/container"
	"github.com/siyoga/rollstory/pkg/logger"
	"log"
)

func provideLogger(di *container.DigContainer) {
	di.Provide(
		func() *logger.Logger {
			customLogger, err := logger.New()
			if err != nil {
				log.Fatalf("инициализация логгера: %w", err)
			}

			return customLogger
		},
	)
}
