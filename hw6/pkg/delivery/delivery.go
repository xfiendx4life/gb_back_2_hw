package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	confStore "github.com/xfiendx4life/gb_back_2_hw/hw6/pkg/confirmation/storage"
	conCase "github.com/xfiendx4life/gb_back_2_hw/hw6/pkg/confirmation/usecase"
	"github.com/xfiendx4life/gb_back_2_hw/hw6/pkg/models"
	ustore "github.com/xfiendx4life/gb_back_2_hw/hw6/pkg/user/storage"
	usercase "github.com/xfiendx4life/gb_back_2_hw/hw6/pkg/user/usecase"
)

type Deliver struct {
	conf conCase.Confirmation
	user usercase.User
}

func New(host, port string, ttl time.Duration) (*Deliver, error) {

	st, err := confStore.NewConfirmationStorage(host,
		port, time.Duration(ttl)*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("can't initialize storage %s", err)
	}
	ust, err := ustore.NewUserStorage(host,
		port, time.Duration(ttl)*time.Hour)
	return &Deliver{
		conf: conCase.New(st),
		user: usercase.NewUserCase(ust),
	}, nil
}

func (d *Deliver) CreateUser(e echo.Context) error {
	u := &models.User{}
	err := json.NewDecoder(e.Request().Body).Decode(u)
	if err != nil {
		log.Errorf("can't decode body")
		return echo.ErrBadRequest
	}
	err = d.user.Register(e.Request().Context(), u.Name, u.Password)
	if err != nil {
		return fmt.Errorf("can't register user: %s", err)
	}
	code, err := d.conf.Create(e.Request().Context(), u.Name)
	if err != nil {
		return fmt.Errorf("casn't create user: %s", err)
	}
	return e.String(http.StatusOK, code)
}

func (d *Deliver) Confirm(e echo.Context) (err error) {
	res := make(map[string]string)
	code := e.FormValue("code")
	name := e.Param("user")
	log.Infof("got code %s from user %s", code, name)
	ok, err := d.conf.Confirm(e.Request().Context(), name, code)
	if err != nil {
		log.Errorf("can't get confirmation: %s", err)
		return echo.ErrBadRequest
	}
	err = d.user.Confirm(e.Request().Context(), name)
	if err != nil {
		log.Errorf("can't change status of user %s", err)
		return echo.ErrInternalServerError
	}
	if !ok {
		res["status"] = "incorrect"
		return e.JSON(http.StatusBadRequest, res)
	}
	res["status"] = "accepted"
	return e.JSON(http.StatusAccepted, res)

}
