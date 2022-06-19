package deliver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/xfiendx4life/gb_back_2_hw/hw9/pkg/models"
	"github.com/xfiendx4life/gb_back_2_hw/hw9/pkg/rest/prcr"
	"github.com/xfiendx4life/gb_back_2_hw/hw9/pkg/storage"
)

type deliver struct {
	st storage.Storage
}

func New(st storage.Storage) prcr.ServerInterface {
	return &deliver{
		st: st,
	}
}

func (d *deliver) CreateList(ctx echo.Context, params prcr.CreateListParams) error {
	params.List.ID = uuid.New()
	err := d.st.Create(ctx.Request().Context(), *params.List)
	if err != nil {
		log.Printf("error while saving data: %s", err)
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Sprintf("error while saving data: %s", err))
	}
	return ctx.String(http.StatusOK, params.List.ID.String())
}

func (d *deliver) DeleteList(ctx echo.Context, listId string) error {
	id, err := uuid.Parse(listId)
	if err != nil {
		log.Printf("can't parse id: %s\n", err)
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Sprintf("can't parse id: %s\n", err))
	}
	log.Printf("ready to delete object with id %s\n", id)
	err = d.st.Delete(ctx.Request().Context(), id)
	if err != nil {
		log.Printf("can't delete list: %s\n", err)
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Sprintf("can't delete list: %s\n", err))
	}
	log.Printf("deleted\n")
	return ctx.NoContent(http.StatusOK)
}

func (d *deliver) UpdateListObject(ctx echo.Context) error {
	var list models.List
	err := json.NewDecoder(ctx.Request().Body).Decode(&list)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Sprintf("Invalid format for parameter list: %s", err))
	}
	err = d.st.Update(ctx.Request().Context(), list.ID, list.Items)
	if err != nil {
		log.Printf("can't update object: %s\n", err)
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Sprintf("can't update object: %s\n", err))
	}
	return err
}

func (d *deliver) ReadList(ctx echo.Context, listId string) error {
	log.Printf("got id %s", listId)
	id, err := uuid.Parse(listId)
	if err != nil {
		log.Printf("can't parse uuid: %s", err)
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Sprintf("Invalid format for id: %s", err))
	}
	data, err := d.st.Read(ctx.Request().Context(), id)
	if err != nil {
		log.Printf("can't get data from storage: %s", err)
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Sprintf("can't get data from storage: %s", err))
	}
	return ctx.JSON(http.StatusOK, data)
}