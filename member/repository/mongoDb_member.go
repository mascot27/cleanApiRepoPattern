package repository

import (
	"context"
	"github.com/mascot27/cleanApiRepoPattern/models"
	"gopkg.in/mgo.v2"
)

type mongoDbMemberRepository struct {
	Conn *mgo.Database
}

func NewMongoDbMemberRepository(conn *mgo.Database) *mongoDbMemberRepository {

	return &mongoDbMemberRepository{Conn: conn}
}

func (mr *mongoDbMemberRepository) Fetch(ctx context.Context, cursor string, num int64) ([]*models.Member, error) {
	panic("implement me")
}

func (mr *mongoDbMemberRepository) GetByID(ctx context.Context, id int64) (*models.Member, error) {
	panic("implement me")
}

func (mr *mongoDbMemberRepository) Store(ctx context.Context, m *models.Member) (int64, error) {
	panic("implement me")
}

func (mr *mongoDbMemberRepository) Delete(ctx context.Context, id int64) (bool, error) {
	panic("implement me")
}
