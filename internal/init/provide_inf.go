package init

import (
	infPsql "github.com/siyoga/rollstory/internal/inf/postgres/init"
	"github.com/siyoga/rollstory/pkg/container"
)

// provideInf для провайдинга всех компонентов из пакета internal/inf
func provideInf(di *container.DigContainer) {
	infPsql.Provide(di)
}
