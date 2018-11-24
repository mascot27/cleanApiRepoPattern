package usecase

import (
	"context"
	"github.com/mascot27/cleanApiRepoPattern/member"
	"github.com/mascot27/cleanApiRepoPattern/models"
	"github.com/segmentio/ksuid"
	"strconv"
	"time"
)

type memberUsecase struct {
	memberRepos    member.MemberRepository
	contextTimeout time.Duration
}

func NewMemberUsecase(memberRepos member.MemberRepository, contextTimeout time.Duration) *memberUsecase {
	return &memberUsecase{memberRepos: memberRepos, contextTimeout: contextTimeout}
}

func (m *memberUsecase) Fetch(c context.Context, cursor string, num int64) ([]*models.Member, string, error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	listMember, err := m.memberRepos.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	nextCursor := ""

	if size := len(listMember); size == int(num) {
		lastId := listMember[num-1].PublicId
		nextCursor = strconv.Itoa(int(lastId))
	}

	return listMember, nextCursor, nil

}

func (m *memberUsecase) GetByID(c context.Context, id int64) (*models.Member, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	res, err := m.memberRepos.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *memberUsecase) Store(c context.Context, newUser *models.Member) (*models.Member, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	guidGenerated := ksuid.New()
	newUser.Guid = guidGenerated.String()

	id, err := m.memberRepos.Store(ctx, newUser)
	if err != nil {
		return nil, err
	}
	newUser.PublicId = id
	return newUser, nil
}

func (m *memberUsecase) Delete(c context.Context, id int64) (bool, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()
	existedArticle, _ := m.memberRepos.GetByID(ctx, id)
	if existedArticle == nil {
		return false, models.NOT_FOUND_ERROR
	}
	return m.memberRepos.Delete(ctx, id)
}
