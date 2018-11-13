package users

import (
	"data/repositories"
	m "datamodels"
	"errors"
)

type IUserDetailService interface {
	New(repo repositories.IRepository, repoReadOnly repositories.IRepositoryReadOnly) IUserDetailService
	GetByID(userid int64) (ok bool, m m.UserDetail, err error)
	AddOrUpdate(detail *m.UserDetail) (affect int64, err error)
}
type UserDetailService struct {
	Repository         repositories.IRepository         `inject:"IRepository"`
	RepositoryReadOnly repositories.IRepositoryReadOnly `inject:"IRepositoryReadOnly"`
}

var _ IUserDetailService = (*UserDetailService)(nil)

func (s *UserDetailService) New(repo repositories.IRepository, repoReadOnly repositories.IRepositoryReadOnly) IUserDetailService {
	return &UserDetailService{repo, repoReadOnly}
}

func (s *UserDetailService) GetByID(userid int64) (ok bool, m m.UserDetail, err error) {
	ok, err = s.RepositoryReadOnly.GetByID(userid, &m)
	return
}

func (s *UserDetailService) AddOrUpdate(detail *m.UserDetail) (affect int64, err error) {
	if detail.UserID <= 0 {
		return 0, errors.New("AddOrUpdate must have key")
	}
	affect, err = s.Repository.AddOrUpdate(detail.UserID, detail)
	return
}
