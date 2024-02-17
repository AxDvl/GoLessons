package auxilaries

const (
	TTValue int = iota
	TTOperator
	TTExpression
)

type Token interface {
	GetType() int
}

type ValueToken interface {
	Token
	Value() float32
	IsInverse() bool
	Resolved() bool
}

const (
	OPPlus int = iota
	OPMul
)

type OperatorToken interface {
	Token
	Operator() int
}

type ExpressionToken interface {
	ValueToken
	LeftOperand() ValueToken
	RightOperand() ValueToken
	Operator() int
}

type tokenStruct struct {
	tokenType    int
	tokenValue   float32
	isInverse    bool
	resolved     bool
	operator     int
	leftOperand  ValueToken
	rightOperand ValueToken
}

func (t tokenStruct) GetType() int {
	return t.tokenType
}

func (t tokenStruct) Value() float32 {
	return t.tokenValue
}

func (t tokenStruct) IsInverse() bool {
	return t.isInverse
}

func (t tokenStruct) Resolved() bool {
	return t.resolved
}

func (t tokenStruct) LeftOperand() ValueToken {
	return t.leftOperand
}

func (t tokenStruct) RightOperand() ValueToken {
	return t.rightOperand
}

func (t tokenStruct) Operator() int {
	return t.operator
}

func NewValueToken(value float32, inverse bool) ValueToken {
	return tokenStruct{tokenType: TTValue, isInverse: inverse, resolved: true, tokenValue: value}
}

func NewOperatorToken(operator int) OperatorToken {
	return tokenStruct{tokenType: TTOperator, operator: operator}
}

func NewExpressionToken(operator int, leftOperand ValueToken, rightOperand ValueToken) ExpressionToken {
	return tokenStruct{
		tokenType:    TTExpression,
		operator:     operator,
		leftOperand:  leftOperand,
		rightOperand: rightOperand,
		resolved:     false,
		isInverse:    leftOperand.IsInverse() && rightOperand.IsInverse(),
	}
}
