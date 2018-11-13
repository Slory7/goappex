package contracts

import "fmt"

type BizError struct {
	Message string
	Status  ResultStatus
	//SubCode int
}

func (er *BizError) Error() string {
	return fmt.Sprintf("Business Error: %s(Status:%d)", er.Message, er.Status)
}
