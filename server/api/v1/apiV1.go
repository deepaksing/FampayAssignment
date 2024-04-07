package v1

import (
	"fmt"

	"github.com/deepaksing/FampayAssignment/store"
	"github.com/labstack/echo"
)

type ApiV1Service struct {
	store *store.Store
}

func NewApiv1Service(store *store.Store) *ApiV1Service {
	return &ApiV1Service{
		store: store,
	}
}

func (a *ApiV1Service) Register(rootGroup *echo.Group) {
	apiv1Group := rootGroup.Group("/api/v1")
	fmt.Println(apiv1Group)
}
