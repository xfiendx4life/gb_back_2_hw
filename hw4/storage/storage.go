package storage

import "github.com/xfiendx4life/gb_back_2_hw/hw4/models"

type Storage interface {
	AddStudent(name, lastname, faculty string) (*models.Student, error) // with exec
	GetStudentByLastname(lastname string) (*models.Student, error)      // queryRow
	GetAllStudentsForFaculty(faculty string) ([]models.Student, error)  //query
}
