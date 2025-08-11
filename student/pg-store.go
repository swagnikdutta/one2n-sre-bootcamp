package student

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDataStore struct {
	Pool *pgxpool.Pool
}

func NewPostgresDataStore() *PostgresDataStore {
	if os.Getenv(databaseUrl) == "" {
		log.Fatalf("missing env variable: %q", databaseUrl)
	}

	pool, err := pgxpool.New(context.Background(), os.Getenv(databaseUrl))
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}

	store := &PostgresDataStore{Pool: pool}
	err = store.init()
	if err != nil {
		log.Fatalf("Error initializing database. Error: %v", err)
		// TODO: When should I panic? vs log.fatal()
	}

	return store
}

func (p *PostgresDataStore) init() error {
	createQuery := `create table if not exists students (
		id serial primary key,
		name text not null,
		age integer 
	)`

	_, err := p.Pool.Exec(context.Background(), createQuery)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDataStore) CreateStudent(s Student) error {
	query := `INSERT INTO students (name, age) values ($1, $2)`
	_, err := p.Pool.Exec(context.Background(), query, s.Name, s.Age)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresDataStore) GetStudent(studentId int) (*Student, error) {
	query := `SELECT * FROM students where id = $1`
	row := p.Pool.QueryRow(context.Background(), query, studentId)

	var id, age int
	var name string
	if err := row.Scan(&id, &name, &age); err != nil {
		return nil, err
	}

	student := &Student{id, name, age}
	return student, nil
}

func (p *PostgresDataStore) UpdateStudent(id int, s Student) error {
	query := `UPDATE students set name = $1, age = $2 WHERE id = $3`
	cTag, err := p.Pool.Exec(context.Background(), query, s.Name, s.Age, id)
	if err != nil {
		return err
	}

	if cTag.RowsAffected() == 0 {
		return errors.New(errStudentNotFound)
	}
	return nil
}

func (p *PostgresDataStore) DeleteStudent(id int) error {
	query := `DELETE from students where id = $1`
	cTag, err := p.Pool.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if cTag.RowsAffected() == 0 {
		return errors.New(errStudentNotFound)
	}
	return nil
}

func (p *PostgresDataStore) ListStudents() ([]Student, error) {
	query := `SELECT * FROM students`
	rows, err := p.Pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// 	TODO: find out why do I need to close rows?

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
