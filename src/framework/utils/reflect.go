package utils

import (
	"reflect"
	"strings"
)

func IsField(o interface{}, fieldname string) bool {
	t := reflect.ValueOf(o).Elem().Type()
	for i := 0; i < t.NumField(); i++ {
		sName := t.Field(i).Name
		if strings.EqualFold(sName, fieldname) {
			return true
		}
	}
	return false
}

func GetInterfaceName(i interface{}) string {
	return reflect.TypeOf(i).Elem().Name()
}
