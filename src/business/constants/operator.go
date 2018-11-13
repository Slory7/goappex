package constants

type Operator string

const (
	GreatThan      Operator = ">"
	GreatThanEqual Operator = ">="
	LessThan       Operator = "<"
	LessThanEqual  Operator = "<="
	Equal          Operator = "="
	NotEqual       Operator = "<>"
	Like           Operator = "like"
	Empty          Operator = ""
)

func (op Operator) String() string {
	return string(op)
}

var _hObjects = map[string]Operator{
	"GreatThan":      GreatThan,
	"GreatThanEqual": GreatThanEqual,
	"LessThan":       LessThan,
	"LessThanEqual":  LessThanEqual,
	"Equal":          Equal,
	"NotEqual":       NotEqual,
	"Like":           Like,
	"":               Empty,
}

func GetOperator(op string) (bool, Operator) {
	val, ok := _hObjects[op]
	return ok, val
}
