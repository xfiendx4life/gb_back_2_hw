package deliver_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/xfiendx4life/gb_back_2_hw/hw9/pkg/models"
	"github.com/xfiendx4life/gb_back_2_hw/hw9/pkg/rest/prcr"
	"github.com/xfiendx4life/gb_back_2_hw/hw9/pkg/rest/prcr/deliver"
)

type mStorage struct {
	err error
}

const tstId = "311e0467-a754-4286-bf89-47c3a83eeb68"

var (
	testCase1 = models.List{
		Items: []*models.Item{
			{Name: "test1", Price: 10},
			{Name: "test2", Price: 30},
		},
	}
	testCreateJSON = `{"items":[{"name":"test1","price":10},{"name":"test2","price":20}]}`
	testFullJSON   = `{"id":"311e0467-a754-4286-bf89-47c3a83eeb68","items":[{"name":"test1","price":10},{"name":"test2","price":20}]}`
	testPatchJson  = `{"id":"311e0467-a754-4286-bf89-47c3a83eeb68","items":[{"name":"test1","price":100}]}`
)

func (mc *mStorage) Create(ctx context.Context, list models.List) error {
	return mc.err
}

func (mc *mStorage) Read(ctx context.Context, id uuid.UUID) (list *models.List, err error) {
	if id.String() == tstId {
		testCase1.ID, _ = uuid.Parse(tstId)
		return &testCase1, mc.err
	}
	return nil, mc.err
}
func (mc *mStorage) Update(ctx context.Context, id uuid.UUID, items []*models.Item) error {
	return mc.err

}

func (mc *mStorage) Delete(ctx context.Context, id uuid.UUID) error {
	return mc.err
}

func TestCreate(t *testing.T) {
	e := echo.New()
	defer e.Close()
	req := httptest.NewRequest(http.MethodPost, "/list/create",
		strings.NewReader(testCreateJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	d := deliver.New(&mStorage{})
	sw := prcr.ServerInterfaceWrapper{Handler: d}
	require.NoError(t, sw.CreateList(ctx))
	require.Equal(t, http.StatusOK, rec.Code)
	require.NotEqual(t, "", rec.Body.String())
}

func TestCreateError(t *testing.T) {
	e := echo.New()
	defer e.Close()
	req := httptest.NewRequest(http.MethodPost, "/list/create",
		strings.NewReader(testCreateJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	d := deliver.New(&mStorage{err: fmt.Errorf("test error")})
	sw := prcr.ServerInterfaceWrapper{Handler: d}
	require.Error(t, sw.CreateList(ctx))
}

func TestRead(t *testing.T) {
	e := echo.New()
	defer e.Close()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/list/:listid")
	ctx.SetParamNames("listId")
	ctx.SetParamValues(tstId)
	d := deliver.New(&mStorage{})
	sw := prcr.ServerInterfaceWrapper{Handler: d}
	require.NoError(t, sw.ReadList(ctx))
	require.Equal(t, http.StatusOK, rec.Code)
	require.NotEqual(t, testFullJSON, rec.Body.String())
}

func TestReadError(t *testing.T) {
	e := echo.New()
	defer e.Close()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	d := deliver.New(&mStorage{})
	sw := prcr.ServerInterfaceWrapper{Handler: d}
	require.Error(t, sw.ReadList(ctx))
}

func TestUpdate(t *testing.T) {
	e := echo.New()
	defer e.Close()
	req := httptest.NewRequest(http.MethodPatch, "/list/update",
		strings.NewReader(testPatchJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	d := deliver.New(&mStorage{})
	sw := prcr.ServerInterfaceWrapper{Handler: d}
	require.NoError(t, sw.UpdateListObject(ctx))
}

func TestUpdateError(t *testing.T) {
	e := echo.New()
	defer e.Close()
	req := httptest.NewRequest(http.MethodPatch, "/list/update",
		strings.NewReader(testPatchJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	d := deliver.New(&mStorage{err: fmt.Errorf("testError")})
	sw := prcr.ServerInterfaceWrapper{Handler: d}
	require.Error(t, sw.UpdateListObject(ctx))
}

func TestDelete(t *testing.T) {
	e := echo.New()
	defer e.Close()
	req := httptest.NewRequest(http.MethodPatch, "/",
		strings.NewReader(testPatchJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/list/delete/:listid")
	ctx.SetParamNames("listId")
	ctx.SetParamValues(tstId)
	d := deliver.New(&mStorage{})
	sw := prcr.ServerInterfaceWrapper{Handler: d}
	require.NoError(t, sw.DeleteList(ctx))
}

func TestDeleteError(t *testing.T) {
	e := echo.New()
	defer e.Close()
	req := httptest.NewRequest(http.MethodPatch, "/",
		strings.NewReader(testPatchJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/list/delete/:listid")
	ctx.SetParamNames("listId")
	ctx.SetParamValues(tstId)
	d := deliver.New(&mStorage{err: fmt.Errorf("testError")})
	sw := prcr.ServerInterfaceWrapper{Handler: d}
	require.Error(t, sw.DeleteList(ctx))
}
