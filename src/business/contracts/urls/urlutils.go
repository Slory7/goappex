package urls

import (
	"business/constants"
	"errors"
	"framework/utils"
	"strings"

	"github.com/kataras/iris"
)

func GetQueryObject(ctx iris.Context, toQueryDto interface{}) (ret QueryObject, err error) {
	start, _ := ctx.URLParamInt("start")
	limit, _ := ctx.URLParamInt("limit")
	if limit > 50 || limit < 1 {
		err = errors.New("param error: limit should be 1-50")
		return
	}
	orderby := ctx.URLParam("orderby")
	if len(orderby) > 0 && !utils.IsField(toQueryDto, orderby) {
		err = errors.New("param error: orderby value not support")
		return
	}
	isdecending, _ := ctx.URLParamBool("isdecending")
	filterby := ctx.URLParam("filterby")
	if len(filterby) > 0 && !utils.IsField(toQueryDto, filterby) {
		err = errors.New("param error: filterby value not support")
		return
	}
	ok, op := constants.GetOperator(ctx.URLParam("op"))
	if !ok {
		err = errors.New("param error: op value not support")
		return
	}
	filterValue := ctx.URLParam("filtervalue")
	ret = QueryObject{
		Start:       start,
		Limit:       limit,
		OrderBy:     strings.ToLower(orderby),
		IsDecending: isdecending,
		FilterBy:    strings.ToLower(filterby),
		Op:          op,
		FilterValue: filterValue,
	}
	return
}
