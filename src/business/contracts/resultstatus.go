package contracts

type ResultStatus int

const (
	Success ResultStatus = (iota + 1)
	NotFound
	Forbidden
	Unauthorized
	Conflict
	BadData
	BadLogic
	Error
)
