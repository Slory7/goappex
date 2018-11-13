package users

import (
	"business/constants"
	"data"
	"data/repositories"
	m "datamodels"
	"errors"
	"strconv"
	"time"
)

type IUserService interface {
	New(repo repositories.IRepository, repoReadOnly repositories.IRepositoryReadOnly) IUserService
	GetByName(name string) (m.User, bool, error)
	GetByEmail(email string) (m.User, bool, error)
	GetAllByName(name string) (userAll m.UserAll, ok bool, err error)
	ListByCondition(start int, limit int, orderby string, isdecending bool, filterby string, op constants.Operator, filterValue interface{}) (users *[]m.User, err error)
	ListWithDetailByCondition(start int, limit int, orderby string, isdecending bool, filterby string, op constants.Operator, filterValue interface{}) (users *[]m.UserWithDetail, err error)
	GetAllSuperUsers() (users *[]m.User, err error)
	AddUser(user *m.User) (affect int64, err error)
	UpdateUser(user *m.User) (affect int64, err error)
	UpdateUserMemo(user *m.UserWithMemo) (affect int64, err error)
	DeleteUsers(ids *[]int64) (affect int64, err error)

	UserLogin(userid int64) (okLast bool, lastLogin m.UserLogins, err error)
	UserLoginTransaction(userid int64) (okLast bool, lastLogin m.UserLogins, err error)
}

type UserService struct {
	Repository         repositories.IRepository         `inject:"IRepository"`
	RepositoryReadOnly repositories.IRepositoryReadOnly `inject:"IRepositoryReadOnly"`
	LoginService       IUserLoginService                `inject:"IUserLoginService"`
	DetailService      IUserDetailService               `inject:"IUserDetailService"`
}

var _ IUserService = (*UserService)(nil)

func (s *UserService) New(repo repositories.IRepository, repoReadOnly repositories.IRepositoryReadOnly) IUserService {
	return &UserService{
		repo,
		repoReadOnly,
		s.LoginService.New(repo, repoReadOnly),
		s.DetailService.New(repo, repoReadOnly),
	}
}

func (s *UserService) GetByName(name string) (m.User, bool, error) {
	user := m.User{UserName: name}
	b, err := s.RepositoryReadOnly.Get(&user)
	return user, b, err
}

func (s *UserService) GetByEmail(email string) (m.User, bool, error) {
	user := m.User{Email: email}
	b, err := s.RepositoryReadOnly.Get(&user)
	return user, b, err
}

func (s *UserService) GetAllByName(name string) (userAll m.UserAll, ok bool, err error) {
	user := m.User{UserName: name}
	ok, err = s.RepositoryReadOnly.Get(&user)
	if ok {
		userAll.User = user
		var userDetail m.UserDetail
		ok2, err := s.RepositoryReadOnly.GetByID(user.Id, &userDetail)
		if err != nil {
			return userAll, ok, err
		}
		if ok2 {
			userAll.UserDetail = userDetail
		}
		userLogins, err := s.LoginService.ListByUserID(user.Id)
		if err != nil {
			return userAll, ok, err
		}
		userAll.UserLogins = userLogins
	}
	return
}

func (s *UserService) ListByCondition(start int, limit int, orderby string, isdecending bool, filterby string, op constants.Operator, filterValue interface{}) (users *[]m.User, err error) {
	query := filterby
	if len(query) > 0 {
		query += " " + op.String() + " ?"
	}
	if op == constants.Like {
		filterValue = "%" + filterValue.(string) + "%"
	}
	usersSlice := make([]m.User, 0)
	users = &usersSlice
	err = s.RepositoryReadOnly.DB().ListByCondition(users, "", start, limit, orderby, isdecending, nil, query, filterValue)
	return
}

func (s *UserService) ListWithDetailByCondition(start int, limit int, orderby string, isdecending bool, filterby string, op constants.Operator, filterValue interface{}) (users *[]m.UserWithDetail, err error) {
	query := filterby
	if len(query) > 0 {
		query += " " + op.String() + " ?"
	}
	if op == constants.Like {
		filterValue = "%" + filterValue.(string) + "%"
	}
	usersSlice := make([]m.UserWithDetail, 0)
	users = &usersSlice
	joins := &[]data.Join{data.Join{Operator: "INNER", TableOrName: "user_detail", JoinCondition: "users.id=user_detail.userid"}}
	err = s.RepositoryReadOnly.DB().ListByCondition(users, "", start, limit, orderby, isdecending, joins, query, filterValue)
	return
}

func (s *UserService) GetAllSuperUsers() (users *[]m.User, err error) {
	usersSlice := make([]m.User, 0)
	users = &usersSlice
	err = s.RepositoryReadOnly.List(users, "issuper", m.User{IsSuper: true})
	return nil, err
}

func (s *UserService) GetNoMemoUsers(users *[]m.User) error {
	err := s.RepositoryReadOnly.Query(users, "select * from `user` where `memo` is null or `memo`=''")
	return err
}

func (s *UserService) AddUser(user *m.User) (affect int64, err error) {
	affect, err = s.Repository.Add(user)
	return
}

func (s *UserService) UpdateUser(user *m.User) (affect int64, err error) {
	affect, err = s.Repository.UpdateByID(user.Id, user, false)
	return
}

func (s *UserService) UpdateUserMemo(user *m.UserWithMemo) (affect int64, err error) {
	affect, err = s.Repository.UpdateByID(user.Id, user, true)
	return
}

func (s *UserService) DeleteUsers(ids *[]int64) (affect int64, err error) {
	users := make([]m.User, len(*ids))
	for idx, id := range *ids {
		users[idx].Id = id
	}
	dbNew := s.Repository.DB().NewTransaction()
	repoNew := s.Repository.New(dbNew)
	for _, u := range users {
		n, err1 := repoNew.Delete(&u)
		if err1 != nil {
			err = err1
			break
		}
		if n != 1 {
			err = errors.New("User not exist:" + strconv.FormatInt(u.Id, 10))
		}
		affect += n
	}
	if err == nil {
		dbNew.Commit()
		dbNew.ClearCacheEntity(new(m.User))
	} else {
		dbNew.RollBack()
	}
	return
}

func (s *UserService) UserLogin(userid int64) (okLast bool, lastLogin m.UserLogins, err error) {
	lastLogin, okLast, _ = s.LoginService.GetLastLogin(userid)
	login := m.UserLogins{UserID: userid, LoginTime: time.Now()}
	_, err = s.LoginService.AddLogin(login)
	if err == nil {
		_, err = s.DetailService.AddOrUpdate(&m.UserDetail{UserID: userid, LastLoginTime: login.LoginTime})
	}
	return
}

func (s *UserService) UserLoginTransaction(userid int64) (okLast bool, lastLogin m.UserLogins, err error) {
	lastLogin, okLast, _ = s.LoginService.GetLastLogin(userid)

	dbNew := s.Repository.DB().NewTransaction()
	repoNew := s.Repository.New(dbNew)

	login := m.UserLogins{UserID: userid, LoginTime: time.Now()}

	loginService := s.LoginService.New(repoNew, repoNew)

	_, err = loginService.AddLogin(login)
	if err == nil {
		detailService := s.DetailService.New(repoNew, repoNew)
		_, err = detailService.AddOrUpdate(&m.UserDetail{UserID: userid, LastLoginTime: login.LoginTime})
	}
	if err == nil {
		dbNew.Commit()
		dbNew.ClearCacheEntity(new(m.UserDetail))
	} else {
		dbNew.RollBack()
	}
	return
}
