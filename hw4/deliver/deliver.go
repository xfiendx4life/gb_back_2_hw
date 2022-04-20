package deliver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/xfiendx4life/gb_back_2_hw/hw4/models"
	"github.com/xfiendx4life/gb_back_2_hw/hw4/storage"
)

type Del struct {
	st storage.Storage
}

func New(st storage.Storage) *Del {
	return &Del{
		st: st,
	}
}

func (d *Del) SetStudent(e echo.Context) error {
	stdnt := &models.Student{}
	err := json.NewDecoder(e.Request().Body).Decode(stdnt)
	if err != nil {
		log.Errorf("can't decode body")
		return echo.ErrBadRequest
	}
	stdnt, err = d.st.AddStudent(stdnt.Name, stdnt.Lastname, stdnt.Faculty)
	if err != nil {
		log.Errorf("can't decode body")
		return echo.ErrInternalServerError
	}
	return e.String(http.StatusCreated, strconv.FormatInt(stdnt.ID, 10))
}

func (d *Del) GetByLastName(e echo.Context) error {
	lName := e.QueryParam("lastname")
	res, err := d.st.GetStudentByLastname(lName)
	if err != nil {
		log.Errorf("can't decode body")
		return echo.ErrInternalServerError
	}
	return e.JSON(http.StatusOK, res)
}
