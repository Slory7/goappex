package urls

import (
	"business/constants"
)

type QueryObject struct {
	Start       int
	Limit       int
	OrderBy     string
	IsDecending bool
	FilterBy    string
	Op          constants.Operator
	FilterValue string
}
