package appstart

import (
	"data/repositories"
	"framework/globals"
	"framework/utils"
	"services/users"

	"github.com/jwells131313/dargo/ioc"
)

func RegisterIoC(
	repo repositories.IRepository,
	repoReadOnly repositories.IRepositoryReadOnly) {

	locator, err := ioc.CreateAndBind("app", func(binder ioc.Binder) error {

		binder.BindConstant(utils.GetInterfaceName((*repositories.IRepository)(nil)), repo)
		binder.BindConstant(utils.GetInterfaceName((*repositories.IRepositoryReadOnly)(nil)), repoReadOnly)

		binder.Bind(utils.GetInterfaceName((*users.IUserDetailService)(nil)), users.UserDetailService{})
		binder.Bind(utils.GetInterfaceName((*users.IUserLoginService)(nil)), users.UserLoginService{})
		binder.Bind(utils.GetInterfaceName((*users.IUserService)(nil)), users.UserService{})

		//binder.BindWithCreator(LoggerServiceName, newLogger).InScope(ioc.PerLookup)

		return nil
	})
	if err == nil {
		globals.ServiceLocator = locator
	}
}
