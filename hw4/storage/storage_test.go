package storage_test

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/xfiendx4life/gb_back_2_hw/hw4/metrics"
	"github.com/xfiendx4life/gb_back_2_hw/hw4/models"
	"github.com/xfiendx4life/gb_back_2_hw/hw4/storage"
)

func initTestTable() {
	db, _ := sql.Open("sqlite3", "test")
	db.Exec(`CREATE TABLE student (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(64),
		lastname VARCHAR(64),
		faculty VARCHAR(5)
	)`)
	db.Close()
}

func TestAddStudent(t *testing.T) {
	initTestTable()
	st, err := storage.New("test", metrics.New(true))
	assert.NoError(t, err)
	stud, err := st.AddStudent("testname", "testlast", "iu")
	os.Remove("test")
	assert.NoError(t, err)
	assert.Equal(t, models.Student{
		ID:       1,
		Name:     "testname",
		Lastname: "testlast",
		Faculty:  "iu",
	}, *stud)

}
