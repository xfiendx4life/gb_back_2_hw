package deliver

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/xfiendx4life/gb_back_2_hw/internal/pkg/env"
)

type EchoDeliver struct {
	env env.Env
}

func New(e env.Env) env.EnvDeliver {
	return &EchoDeliver{
		env: e,
	}
}

func (e *EchoDeliver) Create(ectx echo.Context) error {
	env, err := e.env.Create(ectx.Request().Context(), "some env")
	if err != nil {
		return fmt.Errorf("can't create env: %s", err)
	}
	return ectx.String(200, fmt.Sprintf("Created env: %v", env))
}

func (e *EchoDeliver) Get(ectx echo.Context) error {
	env, err := e.env.Get(ectx.Request().Context(), "some env")
	if err != nil {
		return fmt.Errorf("can't get env: %s", err)
	}
	return ectx.String(200, fmt.Sprintf("Got env: %v", env))
}
