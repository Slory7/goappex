package controllers

import (
	"business/contracts/urls"
	m "datamodels"
	"errors"
	"framework/globals"
	"services/users"
	"strconv"
	"strings"
	v "viewmodels"

	//. "github.com/ahmetb/go-linq"
	"github.com/jinzhu/copier"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

type UserController struct {
	Service users.IUserService
}

//start=0&limit=1&orderby=username&filterby=username&op=GreatThan&filtervalue=a
func (c *UserController) Get(ctx iris.Context) (userDtos []v.UserDto, err error) {
	obj, err := urls.GetQueryObject(ctx, new(v.UserDto))

	if err == nil {
		users, err := c.Service.ListByCondition(obj.Start, obj.Limit, obj.OrderBy, obj.IsDecending, obj.FilterBy, obj.Op, obj.FilterValue)
		if err != nil {
			return userDtos, err
		}
		// From(users).OrderByT(func(c m.User) string {
		// 	return c.UserName
		// }).ToSlice(&users)

		copier.Copy(&userDtos, &users)

	}
	return
}

func (c *UserController) GetDetails(ctx iris.Context) (userDtos []v.UserWithDetailDto, err error) {
	obj, err := urls.GetQueryObject(ctx, new(v.UserDto))

	if err == nil {
		users, err := c.Service.ListWithDetailByCondition(obj.Start, obj.Limit, obj.OrderBy, obj.IsDecending, obj.FilterBy, obj.Op, obj.FilterValue)
		if err != nil {
			return userDtos, err
		}

		copier.Copy(&userDtos, &users)
	}
	return
}

func (c *UserController) GetBy(name string) (userDto v.UserDto, err error, found bool) {
	user, found, err := c.Service.GetByName(name)
	if err == nil && found {
		copier.Copy(&userDto, &user)
	}
	return
}

func (c *UserController) GetDetail(name string, ctx iris.Context) (userAllDto v.UserAllDto, err error, found bool) {
	user, found, err := c.Service.GetAllByName(name)
	if err == nil && found {
		copier.Copy(&userAllDto, &user)
	}
	return
}

func (c *UserController) GetByEmail(email string, ctx iris.Context) (userDto v.UserDto, err error, found bool) {
	user, found, err := c.Service.GetByEmail(email)
	if err == nil && found {
		copier.Copy(&userDto, &user)
	}
	return
}

func (c *UserController) GetSupers(ctx iris.Context) (userDtos []v.UserDto, err error) {
	users, err := c.Service.GetAllSuperUsers()
	if err == nil {
		copier.Copy(&userDtos, &users)
	}
	return
}

func (c *UserController) Post(ctx iris.Context) (userDto v.UserDto, err error) {
	userInfoDto := v.UserInfoDto{}
	err = ctx.ReadJSON(&userInfoDto)
	if err == nil {
		if err = globals.Validator.Struct(userInfoDto); err == nil {
			user := m.User{}
			copier.Copy(&user, &userInfoDto)
			user.IsSuper = false
			affect, err := c.Service.AddUser(&user)
			if err != nil {
				return userDto, err
			}
			if affect != 1 {
				err = errors.New("Add failed")
				return userDto, err
			}
			copier.Copy(&userDto, &user)

		} else {
			err = globals.Validator.GetTranslatedError(err, ctx.GetHeader("Accept-Language"))
		}
	}
	return
}

func (c *UserController) Put(ctx iris.Context) (userDto v.UserDto, err error) {
	err = ctx.ReadJSON(&userDto)
	if err == nil {
		if err = globals.Validator.Struct(userDto); err == nil {
			user := m.User{}
			copier.Copy(&user, &userDto)
			affect, err := c.Service.UpdateUser(&user)
			if err != nil {
				return userDto, err
			}
			if affect != 1 {
				err = errors.New("Update failed: User not exists")
				return userDto, err
			}
			copier.Copy(&userDto, &user)
		} else {
			err = globals.Validator.GetTranslatedError(err, ctx.GetHeader("Accept-Language"))
		}
	}
	return userDto, err
}

func (c *UserController) PutMemo(ctx iris.Context) (userDto v.UserMemoDto, err error) {
	err = ctx.ReadJSON(&userDto)
	if err == nil {
		if err = globals.Validator.Struct(userDto); err == nil {
			user := m.UserWithMemo{}
			copier.Copy(&user, &userDto)
			affect, err := c.Service.UpdateUserMemo(&user)
			if err != nil {
				return userDto, err
			}
			if affect != 1 {
				err = errors.New("Update failed: User not exists")
				return userDto, err
			}
			copier.Copy(&userDto, &user)

		} else {
			err = globals.Validator.GetTranslatedError(err, ctx.GetHeader("Accept-Language"))
		}
	}
	return
}

func (c *UserController) DeleteBy(ids string, ctx iris.Context) interface{} {
	if len(ids) == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
	}
	idsArr := strings.Split(ids, ",")
	idsInts := make([]int64, len(idsArr))
	for idx, sid := range idsArr {
		n, er := strconv.ParseInt(sid, 0, 64)
		if er != nil {
			return iris.StatusBadRequest
		}
		idsInts[idx] = n
	}
	_, err := c.Service.DeleteUsers(&idsInts)
	if err == nil {
		return iris.Map{"deleted": ids}
	}
	return err
}

func (c *UserController) PostUserLogin(name string) (userAllDto v.UserAllDto, err error, found bool) {
	user, found, err := c.Service.GetByName(name)
	if found {
		okLast, lastLogin, err := c.Service.UserLoginTransaction(user.Id)
		if err != nil {
			return userAllDto, err, found
		}
		userAll := m.UserAll{User: user}
		if okLast {
			userAll.LastLoginTime = lastLogin.LoginTime
			userAll.UserLogins = []m.UserLogins{lastLogin}
		}
		copier.Copy(&userAllDto, &userAll)
	}
	return
}

func (c *UserController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/email/{email:string}", "GetByEmail")
	b.Handle("GET", "/detail/{name:string}", "GetDetail")
	b.Handle("POST", "/login/{name:string}", "PostUserLogin")
}
