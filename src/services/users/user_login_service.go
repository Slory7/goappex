package users

import (
	"data/repositories"
	m "datamodels"
)

type IUserLoginService interface {
	New(repo repositories.IRepository, repoReadOnly repositories.IRepositoryReadOnly) IUserLoginService
	ListByUserID(userid int64) (logins []m.UserLogins, err error)
	GetLastLogin(userid int64) (login m.UserLogins, ok bool, err error)
	AddLogin(login m.UserLogins) (affect int64, err error)
	DeleteLogin(id int64) (affect int64, err error)
}

type UserLoginService struct {
	Repository         repositories.IRepository         `inject:"IRepository"`
	RepositoryReadOnly repositories.IRepositoryReadOnly `inject:"IRepositoryReadOnly"`
}

var _ IUserLoginService = (*UserLoginService)(nil)

func (s *UserLoginService) New(repo repositories.IRepository, repoReadOnly repositories.IRepositoryReadOnly) IUserLoginService {
	return &UserLoginService{repo, repoReadOnly}
}

func (s *UserLoginService) ListByUserID(userid int64) (logins []m.UserLogins, err error) {
	err = s.RepositoryReadOnly.List(&logins, "", m.UserLogins{UserID: userid})
	return
}

func (s *UserLoginService) GetLastLogin(userid int64) (login m.UserLogins, ok bool, err error) {
	login = m.UserLogins{UserID: userid}
	ok, err = s.RepositoryReadOnly.GetByOrder(&login, "id", true)
	return
}

func (s *UserLoginService) getOlderThan10Logins(userid int64) (login []m.UserLogins, err error) {
	queryLogin := m.UserLogins{UserID: userid}
	err = s.Repository.DB().ListByObjects(&login, "id", 10, -1, "id", true, nil, queryLogin)
	return
}

func (s *UserLoginService) AddLogin(login m.UserLogins) (affect int64, err error) {
	affect, err = s.Repository.Add(login)
	if err == nil {
		oldLogins, err := s.getOlderThan10Logins(login.UserID)
		if err == nil {
			for i := 0; i < len(oldLogins); i++ {
				s.DeleteLogin(oldLogins[i].Id)
			}
		}
	}
	return
}

func (s *UserLoginService) DeleteLogin(id int64) (affect int64, err error) {
	affect, err = s.Repository.Delete(m.UserLogins{Id: id})
	return
}
