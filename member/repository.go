package member

import (
	"context"
	"github.com/mascot27/cleanApiRepoPattern/models"
)

type MemberRepository interface {
	GetByID(ctx context.Context, id int64) (*models.Member, error)
	Store(ctx context.Context, m *models.Member) (int64, error)
	Delete(ctx context.Context, id int64) (bool, error)
}
