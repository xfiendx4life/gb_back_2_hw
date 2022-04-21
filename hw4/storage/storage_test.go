package storage_test

import (
	"database/sql"
	"log"
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

func fillData(data []models.Student) {
	db, _ := sql.Open("sqlite3", "test")
	for _, std := range data {
		stmt, _ := db.Prepare("INSERT INTO student(name, lastname, faculty) values (?,?,?)")
		stmt.Exec(std.Name, std.Lastname, std.Faculty)
	}
	db.Close()
}

var testTarget = models.Student{
	ID:       1,
	Name:     "testname",
	Lastname: "testlast",
	Faculty:  "iu",
}

var testSlice = []models.Student{
	{
		ID:       1,
		Name:     "testname1",
		Lastname: "testlast1",
		Faculty:  "iu",
	},
	{
		ID:       2,
		Name:     "testname2",
		Lastname: "testlast2",
		Faculty:  "iu",
	},
	{
		ID:       3,
		Name:     "testname3",
		Lastname: "testlast3",
		Faculty:  "mt",
	},
}

var st storage.Storage

func init() {
	var err error
	st, err = storage.New("test", metrics.New(true))
	if err != nil {
		log.Fatal(err)
	}
}

func TestAddStudent(t *testing.T) {
	initTestTable()
	stud, err := st.AddStudent("testname", "testlast", "iu")
	os.Remove("test")
	assert.NoError(t, err)
	assert.Equal(t, testTarget, *stud)

}

func TestAddStudentError(t *testing.T) {
	_, err := st.AddStudent("testname", "testlast", "iu")
	os.Remove("test")
	assert.Error(t, err)
}

func TestGetStudentByLastName(t *testing.T) {
	initTestTable()
	fillData([]models.Student{testTarget})
	res, err := st.GetStudentByLastname(testTarget.Lastname)
	os.Remove("test")
	assert.NoError(t, err)
	assert.Equal(t, testTarget, *res)
}

func TestGetStudentByLastNameError(t *testing.T) {
	initTestTable()
	fillData([]models.Student{testTarget})
	res, err := st.GetStudentByLastname("wrong lastname")
	os.Remove("test")
	assert.Error(t, err)
	assert.Nil(t, res)
}

// ! Some unexplaining things are going on here
func TestGetAllStudentsForFaculty(t *testing.T) {
	initTestTable()
	fillData(testSlice)
	res, err := st.GetAllStudentsForFaculty(testTarget.Faculty)
	os.Remove("test")
	assert.NoError(t, err)
	assert.Equal(t, testSlice[:2], res)
}

func TestGetAllStudentsForFacultyEmpty(t *testing.T) {
	initTestTable()
	fillData(testSlice)
	res, _ := st.GetAllStudentsForFaculty("wrong faculty")
	os.Remove("test")
	assert.Empty(t, res)
}
