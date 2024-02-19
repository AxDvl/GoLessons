package auxilaries

import (
	"encoding/json"
	"fmt"
	"io"
)

func GetStringFromBody(r io.ReadCloser) (string, error) {
	buf := make([]byte, 100)
	s := ""
	for {
		n, err := r.Read(buf)
		s += string(buf[:n])

		if err == io.EOF {
			break
		}

		if err != nil {
			return s, err
		}
	}
	return s, nil

}

func GetBodyAsJson(r io.ReadCloser, obj any) error {
	bodyText, err := GetStringFromBody(r)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(bodyText), obj)
	if err != nil {
		return err
	}
	return nil
}

func PrintTokenList(tokens []Token) {
	for _, token := range tokens {
		PrintToken(token)
	}
	fmt.Println()
}

func PrintToken(token Token) {
	switch token.GetType() {
	case TTValue:
		if token.(ValueToken).IsInverse() {
			fmt.Print("1/")
		}
		fmt.Print(token.(ValueToken).Value())
	case TTOperator:
		if token.(OperatorToken).Operator() == OPPlus {
			fmt.Print(" + ")
		} else {
			fmt.Print(" * ")
		}
	case TTExpression:
		expr := token.(ExpressionToken)
		fmt.Print("[")
		if expr.LeftOperand().GetType() == TTValue {
			fmt.Print(expr.LeftOperand().Value())
		}
		if expr.LeftOperand().GetType() == TTExpression {
			fmt.Print("Expr")
		}
		if expr.Operator() == OPMul {
			fmt.Print("*")
		} else {
			fmt.Print("+")
		}

		if expr.RightOperand().GetType() == TTValue {
			fmt.Print(expr.RightOperand().Value())
		}
		if expr.RightOperand().GetType() == TTExpression {
			fmt.Print("Expr")
		}
		fmt.Print("]")
	}
}
