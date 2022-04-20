package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xfiendx4life/gb_back_2_hw/hw4/metrics"
	"github.com/xfiendx4life/gb_back_2_hw/hw4/models"
)

type store struct {
	*sql.DB
	*metrics.Metr
}

func New(path string, mtrcs *metrics.Metr) (Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %s", err)
	}
	st := &store{
		db,
		mtrcs,
	}
	if err = st.CreateTable(); err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return st, nil
}

func (st *store) CreateTable() error {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occured")
		}
	}()
	_, err := st.MesurableExec(st.Exec, `CREATE TABLE IF NOT EXISTS student (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(64),
		lastname VARCHAR(64),
		faculty VARCHAR(5)
	)`)
	if err != nil {
		return fmt.Errorf("can't perform create: %s", err)
	}
	return nil
}

func (st *store) AddStudent(name, lastname, faculty string) (*models.Student, error) { //* with exec
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occured")
		}
	}()
	res, err := st.MesurableExec(st.Exec, "INSERT INTO student(name, lastname, faculty) values (?, ?, ?)",
		name, lastname, faculty)
	if err != nil {
		return nil, fmt.Errorf("can't perform exec: %s", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("can't get last id: %s", err)
	}
	return &models.Student{
		ID:       id,
		Name:     name,
		Lastname: lastname,
		Faculty:  faculty,
	}, nil
}

func (st *store) GetStudentByLastname(lastname string) (*models.Student, error) { //* queryRow
	stud := models.Student{}
	row := st.MesurableQueryRow(st.QueryRow, `SELECT * FROM student WHERE lastname=? LIMIT 1`, lastname)
	err := row.Scan(&stud.ID, &stud.Name, &stud.Lastname, &stud.Faculty)
	if err != nil {
		return nil, fmt.Errorf("can't read data: %s", err)
	}
	return &stud, nil
}
func (st *store) GetAllStudentsForFaculty(faculty string) ([]models.Student, error) { //* query
	stds := make([]models.Student, 0)
	rows, err := st.MesurableQuery(st.Query, `SELECT * FROM student WHERE faculty=?`, faculty)
	if err != nil {
		return nil, fmt.Errorf("can't get students for faculty %s", err)
	}
	for rows.Next() {
		stud := models.Student{}
		err = rows.Scan(&stud.ID, &stud.Name, &stud.Lastname, &stud.Faculty)
		if err != nil {
			return nil, fmt.Errorf("can't scan results %s", err)
		}
		stds = append(stds, stud)
	}
	return stds, nil
}
