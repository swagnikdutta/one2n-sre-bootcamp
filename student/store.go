package student

import (
	"database/sql"
	"errors"
	"log"
	"os"
)

type Store interface {
	CreateStudent(s Student) (int, error)
	GetStudent(studentId int) (*Student, error)
	UpdateStudent(id int, s Student) error
	DeleteStudent(id int) error
	ListStudents() ([]Student, error)
}

type SQLiteDataStore struct {
	db *sql.DB
}

func NewSQLiteDataStore() *SQLiteDataStore {
	db, err := sql.Open(sqliteDriverName, os.Getenv(dbPath))
	if err != nil {
		panic(err)
	}
	// TODO: do we need to close connection?
	// defer db.Close()
	store := &SQLiteDataStore{db: db}

	err = store.init()
	if err != nil {
		log.Fatalf("Error initializing database. Error: %v", err)
		// TODO: When should I panic? vs log.fatal()
	}

	return store
}

func (s *SQLiteDataStore) init() error {
	createQuery := `create table if not exists students (
		id integer primary key autoincrement,
		name text not null,
		age integer 
	)`

	_, err := s.db.Exec(createQuery)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteDataStore) CreateStudent(student Student) (int, error) {
	query := `insert into students (name, age) values (?, ?)`
	res, err := s.db.Exec(query, student.Name, student.Age)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil
}

func (s *SQLiteDataStore) GetStudent(studentId int) (*Student, error) {
	query := `select * from students where id = ?`
	row := s.db.QueryRow(query, studentId)

	var id, age int
	var name string
	if err := row.Scan(&id, &name, &age); err != nil {
		return nil, err
	}

	student := &Student{id, name, age}
	return student, nil
}

func (s *SQLiteDataStore) UpdateStudent(studentId int, student Student) error {
	query := `update students set name = ?, age = ? where id = ?`
	res, err := s.db.Exec(query, student.Name, student.Age, studentId)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New(errStudentNotFound)
	}
	return nil
}

func (s *SQLiteDataStore) DeleteStudent(studentId int) error {
	query := `delete from students where id = ?`
	res, err := s.db.Exec(query, studentId)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New(errStudentNotFound)
	}
	return nil
}

func (s *SQLiteDataStore) ListStudents() ([]Student, error) {
	query := `select * from students`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// TODO: Again, why do I need to close rows?

	var students []Student
	for rows.Next() {
		var id, age int
		var name string
		if err := rows.Scan(&id, &name, &age); err != nil {
			return nil, err
		}
		students = append(students, Student{id, name, age})
	}
	return students, nil
}
