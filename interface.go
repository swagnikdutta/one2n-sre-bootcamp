package main

type StudentStore interface {
	CreateStudent(s Student) (int, error)
	GetStudent(studentId int) (*Student, error)
	UpdateStudent(id int, s Student) error
	DeleteStudent(id int) error
	ListStudents() ([]Student, error)
}
