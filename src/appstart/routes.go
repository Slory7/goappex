package appstart

import (
	"controllers"
	"framework/globals"
	"framework/utils"
	"services/users"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

func ConfigureRoutes(app *iris.Application) {
	mvc.Configure(app.Party("/user"), userPart)
}

//user controller
func userPart(app *mvc.Application) {

	serviceName := utils.GetInterfaceName((*users.IUserService)(nil))
	userService, _ := globals.ServiceLocator.GetDService(serviceName)
	app.Register(userService)

	app.Handle(new(controllers.UserController))
}
