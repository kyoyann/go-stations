package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
//エラーの場合に*model.TODOにnilを設定するとpanicが発生する。UpdateTODOResponseのTODOにメモリを格納しており、ぬるぽが起きる
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	result, err := s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		log.Println(err)
		return &model.TODO{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return &model.TODO{}, err
	}
	todo := model.TODO{}

	if err = s.db.QueryRowContext(ctx, confirm, id).Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
		log.Println(err)
		return &model.TODO{}, err
	}

	return &todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)
	todos := []*model.TODO{}
	var rows *sql.Rows
	var err error
	if prevID == 0 {
		rows, err = s.db.QueryContext(ctx, read, size)
	} else {
		rows, err = s.db.QueryContext(ctx, readWithID, prevID, size)
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		todo := model.TODO{}
		if err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			log.Fatalf("getRows rows.Scan error err:%v", err)
		}
		todos = append(todos, &todo)
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return todos, nil
}

// UpdateTODO updates the TODO on DB.
//エラーの場合に*model.TODOにnilを設定するとpanicが発生する。UpdateTODOResponseのTODOにメモリを格納しており、ぬるぽが起きる
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	result, err := s.db.ExecContext(ctx, update, subject, description, id)
	if err != nil {
		log.Println(err)
		return &model.TODO{}, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return &model.TODO{}, err
	}
	if rows == 0 {
		return &model.TODO{}, &model.ErrNotFound{}
	}

	todo := model.TODO{}

	if err = s.db.QueryRowContext(ctx, confirm, id).Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
		log.Println(err)
		return &model.TODO{}, err
	}

	return &todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	deleteFmt := fmt.Sprintf(`DELETE FROM todos WHERE id IN (?%s)`, strings.Repeat(", ?", len(ids)-1))
	var arg []interface{}
	for _, v := range ids {
		arg = append(arg, v)
	}
	result, err := s.db.ExecContext(ctx, deleteFmt, arg...)
	if err != nil {
		log.Println(err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return &model.ErrNotFound{}
	}
	return nil
}
