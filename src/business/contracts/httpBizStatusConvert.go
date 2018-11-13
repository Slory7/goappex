package contracts

import (
	"net/http"
)

func HttpToBizStatus(httpStatus int) (result ResultStatus) {
	result = Success
	switch {
	case httpStatus == 404:
		result = NotFound
	case httpStatus == 409:
		result = Conflict
	case httpStatus == 403:
		result = Forbidden
	case httpStatus == 422:
		result = BadLogic
	case httpStatus == 400:
		result = BadData
	case httpStatus == 401:
		result = Unauthorized
	case httpStatus >= 300:
		result = Error
	}
	return
}

func BizStatusToHttp(status ResultStatus) (result int) {
	result = http.StatusOK
	switch {
	case status == NotFound:
		result = http.StatusNotFound
	case status == Conflict:
		result = http.StatusConflict
	case status == Forbidden:
		result = http.StatusForbidden
	case status == BadLogic:
		result = http.StatusUnprocessableEntity
	case status == BadData:
		result = http.StatusBadRequest
	case status == Unauthorized:
		result = http.StatusUnauthorized
	case status == Error:
		result = http.StatusInternalServerError
	}
	return
}
