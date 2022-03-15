package deliver

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/xfiendx4life/gb_back_2_hw/internal/pkg/connector"
	"github.com/xfiendx4life/gb_back_2_hw/models"
)

type ConnDel struct {
	uc connector.Connector
}

func New(uc connector.Connector) connector.ConnectorDeliver {
	return &ConnDel{
		uc: uc,
	}
}

func (c *ConnDel) AddToEnv(ectx echo.Context) error {
	return c.uc.AddToEnv(ectx.Request().Context(), models.User{
		Name: "username",
		ID:   uuid.New(),
	}, models.Env{
		ID:   uuid.New(),
		Name: "someenv",
	})
}

func (c *ConnDel) GetByEnv(ectx echo.Context) error {
	mods, err := c.uc.GetByEnv(ectx.Request().Context(), models.Env{
		ID:   uuid.New(),
		Name: "someenv",
	})
	if err != nil {
		return fmt.Errorf("can't get by env %s", err)
	}
	return ectx.JSON(200, mods)
}

func (c *ConnDel) GetByUser(ectx echo.Context) error {
	mods, err := c.uc.GetByUser(ectx.Request().Context(), models.User{
		Name: "username",
		ID:   uuid.New(),
	})
	if err != nil {
		return fmt.Errorf("can't get GetByEnvby env %s", err)
	}
	return ectx.JSON(200, mods)
}
func (c *ConnDel) DeleteUserFromEnv(ectx echo.Context) error {
	return c.uc.DeleteUserFromEnv(ectx.Request().Context(),
		models.User{
			Name: "username",
			ID:   uuid.New(),
		}, models.Env{
			ID:   uuid.New(),
			Name: "someenv",
		})
}
