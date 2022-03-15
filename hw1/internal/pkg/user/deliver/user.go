package deliver

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/xfiendx4life/gb_back_2_hw/internal/pkg/user"
)

type EchoDeliver struct {
	u user.User
}

func New(u user.User) user.UserDeliver {
	return &EchoDeliver{
		u: u,
	}
}

func (e *EchoDeliver) Create(ectx echo.Context) error {
	u, err := e.u.Create(ectx.Request().Context(), "Name")
	if err != nil {
		return fmt.Errorf("error while creating element %s", err)
	}
	return ectx.String(200, fmt.Sprintf("%v", u))
}

func (e *EchoDeliver) GetByName(ectx echo.Context) error {
	u, err := e.u.GetByName(ectx.Request().Context(), "Name")
	if err != nil {
		return fmt.Errorf("error while getting element %s", err)
	}
	return ectx.String(200, fmt.Sprintf("%v", u))
}
