package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"

	_ "github.com/mascot27/cleanApiRepoPattern/member"
	models "github.com/mascot27/cleanApiRepoPattern/models"
)

type mysqlMemberRepository struct {
	Conn *sql.DB
}

func (m *mysqlMemberRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Member, error) {

	rows, err := m.Conn.QueryContext(ctx, query, args...)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer rows.Close()
	result := make([]*models.Member, 0)
	for rows.Next() {
		t := new(models.Member)
		err = rows.Scan(
			&t.ID,
			&t.Name,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func NewMysqlMemberRepository(conn *sql.DB) *mysqlMemberRepository {
	return &mysqlMemberRepository{Conn: conn}
}

func (m *mysqlMemberRepository) Fetch(ctx context.Context, cursor string, num int64) ([]*models.Member, error) {
	query := `SELECT id,name FROM member WHERE ID > ? LIMIT ?`

	return m.fetch(ctx, query, cursor, num)
}

func (m *mysqlMemberRepository) GetByID(ctx context.Context, id int64) (*models.Member, error) {
	query := `SELECT id, name FROM member WHERE ID = ?`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}

	res := &models.Member{}
	if len(list) > 0 {
		res = list[0]
	} else {
		return nil, models.NOT_FOUND_ERROR
	}

	return res, nil
}

func (m *mysqlMemberRepository) Store(ctx context.Context, newMember *models.Member) (int64, error) {
	query := `INSERT  member SET ID=? , Name=? `
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {

		return 0, err
	}

	res, err := stmt.ExecContext(ctx, newMember.ID, newMember.Name)
	if err != nil {

		return 0, err
	}
	return res.LastInsertId()
}

func (m *mysqlMemberRepository) Delete(ctx context.Context, id int64) (bool, error) {
	query := "DELETE FROM member WHERE id = ?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return false, err
	}
	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return false, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", rowsAffected)
		logrus.Error(err)
		return false, err
	}

	return true, nil
}
