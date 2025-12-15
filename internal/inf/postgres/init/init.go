package init

import (
	"fmt"
	psql "github.com/siyoga/rollstory/internal/inf/postgres/public"
	"github.com/siyoga/rollstory/pkg/container"
	"log"
)

func Provide(di *container.DigContainer) {
	di.Provide(func() *psql.Connection {
		conn, err := psql.New()
		if err != nil {
			log.Fatal(fmt.Errorf("can't initialize postgres connection: %w", err))
		}

		return conn
	})

	di.Bind(new(psql.Connection), new(psql.DB))
}
