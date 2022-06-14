package storage_test

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/xfiendx4life/gb_back_2_hw/hw9/pkg/models"
	"github.com/xfiendx4life/gb_back_2_hw/hw9/pkg/storage"
)

type testStore struct {
	path string
}

func (ts *testStore) prepare() {
	os.Mkdir(ts.path, 0777)
}

func (ts *testStore) clean() {
	os.RemoveAll(ts.path)
}

func (ts *testStore) writeData(lst *models.List) {
	file, err := os.Create(filepath.Join(ts.path, lst.ID.String()))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	err = json.NewEncoder(file).Encode(lst)
	if err != nil {
		log.Fatal(err)
	}
}

func TestCreate(t *testing.T) {
	ts := &testStore{
		path: "./testpath",
	}
	// defer ts.clean()
	ts.prepare()
	st, err := storage.New(ts.path, "")
	require.NoError(t, err)
	lst := models.List{
		ID: uuid.New(),
		Items: []*models.Item{
			{Name: "test1",
				Price: 10},
			{Name: "test2",
				Price: 20},
		},
	}
	err = st.Create(context.Background(), lst)
	require.NoError(t, err)
}

func TestCreateNotExists(t *testing.T) {
	ts := &testStore{
		path: "./testpath",
	}
	st, err := storage.New(ts.path, "n")
	require.NoError(t, err)
	lst := models.List{
		ID: uuid.New(),
		Items: []*models.Item{
			{Name: "test1",
				Price: 10},
			{Name: "test2",
				Price: 20},
		},
	}
	err = st.Create(context.Background(), lst)
	require.NoError(t, err)
}

func TestRead(t *testing.T) {
	ts := &testStore{
		path: "./testpath",
	}
	ts.prepare()
	defer ts.clean()
	lst := models.List{
		ID: uuid.New(),
		Items: []*models.Item{
			{Name: "test1",
				Price: 10},
			{Name: "test2",
				Price: 20},
		},
	}
	ts.writeData(&lst)
	st, err := storage.New(ts.path, "")
	require.NoError(t, err)
	toCompare, err := st.Read(context.Background(), lst.ID)
	require.NoError(t, err)
	require.EqualValues(t, lst.Items, toCompare.Items)

}

func TestUpdate(t *testing.T) {
	ts := &testStore{
		path: "./testpath",
	}
	ts.prepare()
	defer ts.clean()
	lst := models.List{
		ID: uuid.New(),
		Items: []*models.Item{
			{Name: "test1",
				Price: 10},
			{Name: "test2",
				Price: 20},
		},
	}
	ts.writeData(&lst)
	st, err := storage.New(ts.path, "")
	require.NoError(t, err)
	err = st.Update(context.Background(), lst.ID, []*models.Item{{Name: "test1",
		Price: 10}})
	require.NoError(t, err)
	toCheck, err := os.ReadFile(filepath.Join(ts.path, lst.ID.String()))
	require.NoError(t, err)
	require.NotEqualValues(t, lst, toCheck)
}

func TestUpdateError(t *testing.T) {
	ts := &testStore{
		path: "./testpath",
	}
	ts.prepare()
	defer ts.clean()
	lst := models.List{
		ID: uuid.New(),
		Items: []*models.Item{
			{Name: "test1",
				Price: 10},
			{Name: "test2",
				Price: 20},
		},
	}
	// ts.writeData(&lst)
	st, err := storage.New(ts.path, "")
	require.NoError(t, err)
	err = st.Update(context.Background(), lst.ID, []*models.Item{{Name: "test1",
		Price: 10}})
	require.Error(t, err)
}

func TestDelete(t *testing.T) {
	ts := &testStore{
		path: "./testpath",
	}
	ts.prepare()
	defer ts.clean()
	lst := models.List{
		ID: uuid.New(),
		Items: []*models.Item{
			{Name: "test1",
				Price: 10},
			{Name: "test2",
				Price: 20},
		},
	}
	ts.writeData(&lst)
	st, err := storage.New(ts.path, "")
	require.NoError(t, err)
	err = st.Delete(context.Background(), lst.ID)
	require.NoError(t, err)
}

func TestDeleteError(t *testing.T) {
	ts := &testStore{
		path: "./testpath",
	}
	ts.prepare()
	defer ts.clean()
	lst := models.List{
		ID: uuid.New(),
		Items: []*models.Item{
			{Name: "test1",
				Price: 10},
			{Name: "test2",
				Price: 20},
		},
	}
	st, err := storage.New(ts.path, "")
	require.NoError(t, err)
	err = st.Delete(context.Background(), lst.ID)
	require.Error(t, err)
}
