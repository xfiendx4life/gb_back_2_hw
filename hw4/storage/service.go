package storage

import (
	"database/sql"
	"fmt"

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
	return &store{
		db,
		mtrcs,
	}, nil
}

func (st *store) AddStudent(name, lastname, faculty string) (*models.Student, error) { //* with exec
	// stmt, err := st.Prepare("INSERT INTO students(name, lastname, faculty) values (?, ?, ?)")
	// if err != nil {
	// 	return nil, fmt.Errorf("can't prepare statement: %s", err)
	// }
	res, err := st.MesurableExec(st.Exec)("INSERT INTO student(name, lastname, faculty) values (?, ?, ?)",
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
	return nil, nil
}
func (st *store) GetAllStudentsForFaculty(faculty string) ([]models.Student, error) { //* query
	return nil, nil
}
