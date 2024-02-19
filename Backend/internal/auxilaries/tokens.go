package auxilaries

import "fmt"

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
}

const (
	OPPlus int = iota
	OPMul
	OPDivide //Операция деления при парсинге не будет использоваться, она будет применяться уже при отправке агентам
)

type OperatorToken interface {
	Token
	Operator() int
}

type ExpressionToken interface {
	ValueToken
	LeftOperand() ValueToken
	RightOperand() ValueToken
	CanBeResolved() bool
	GetID() string
	Operator() int
	SetLeftOperand(value ValueToken)
	SetRightOperand(value ValueToken)
}

type tokenStruct struct {
	tokenType    int
	tokenValue   float32
	isInverse    bool
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

func (t tokenStruct) CanBeResolved() bool {
	_, leftIsExpression := t.leftOperand.(ExpressionToken)
	_, rightIsExpression := t.rightOperand.(ExpressionToken)
	return !leftIsExpression && !rightIsExpression
}

func (t tokenStruct) GetID() string {
	if !t.CanBeResolved() {
		return ""
	}
	operatorText := "+"
	leftPrefix := ""
	rightPrefix := ""

	if t.operator == OPMul {
		operatorText = "*"
		if t.leftOperand.IsInverse() && !t.rightOperand.IsInverse() {
			leftPrefix = "!"
		} else if !t.leftOperand.IsInverse() && t.rightOperand.IsInverse() {
			rightPrefix = "!"
		}
	}

	return fmt.Sprintf("%s%f%s%s%f", leftPrefix, t.leftOperand.Value(), operatorText, rightPrefix, t.rightOperand.Value())
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

func (t *tokenStruct) SetLeftOperand(value ValueToken) {
	t.leftOperand = value
}
func (t *tokenStruct) SetRightOperand(value ValueToken) {
	t.rightOperand = value
}

func NewValueToken(value float32, inverse bool) ValueToken {
	return tokenStruct{tokenType: TTValue, isInverse: inverse, tokenValue: value}
}

func NewOperatorToken(operator int) OperatorToken {
	return tokenStruct{tokenType: TTOperator, operator: operator}
}

func NewExpressionToken(operator int, leftOperand ValueToken, rightOperand ValueToken) ExpressionToken {
	return &tokenStruct{
		tokenType:    TTExpression,
		operator:     operator,
		leftOperand:  leftOperand,
		rightOperand: rightOperand,
		isInverse:    leftOperand.IsInverse() && rightOperand.IsInverse(),
	}
}
