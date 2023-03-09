package forms

type Operation string

var matchOperation = map[Operation]string{
	"update": "upd",
}

type MatchForm struct {
	OperationValue Operation
}
