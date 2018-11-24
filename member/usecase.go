package member

import (
	"context"
	"github.com/mascot27/cleanApiRepoPattern/models"
)

type MemberUsecase interface {
	GetByID(ctx context.Context, id int64) (*models.Member, error)
	Store(context.Context, *models.Member) (*models.Member, error)
	Delete(ctx context.Context, id int64) (bool, error)
}
