package auxilaries

import (
	"errors"
	"strconv"
)

func Parse(text string) ([]Token, error) {
	var res []Token
	plus := rune('+')
	minus := rune('-')
	mul := rune('*')
	divide := rune('/')
	inverse := false
	negative := false

	var buf []rune
	for i, r := range text {
		if r == plus || r == minus || r == mul || r == divide {
			if i == 0 {
				if r == minus {
					buf = append(buf, r)
					continue
				}
				return res, errors.New("Выражение не может начинаться с оператора")
			}

			operator := OPPlus
			switch r {
			case plus:
				operator = OPPlus
			case minus:
				operator = OPPlus
			case mul:
				operator = OPMul
			case divide:
				operator = OPMul
			}

			if len(buf) > 0 {
				val, err := strconv.ParseFloat(string(buf), 32)
				if err != nil {
					return res, err
				}
				if negative {
					val = -val
				}
				res = append(res, NewValueToken(float32(val), inverse))
				buf = buf[:0]
			}
			res = append(res, NewOperatorToken(operator))
			negative = r == minus
			inverse = r == divide
		} else {
			buf = append(buf, r)
		}
	}
	if len(buf) > 0 {
		val, err := strconv.ParseFloat(string(buf), 32)
		if err != nil {
			return res, err
		}
		if negative {
			val = -val
		}
		res = append(res, NewValueToken(float32(val), inverse))
	}
	return res, nil
}

func BuildGraph(tokens []Token) ([]Token, error) {
	var res []Token

	const (
		exLeftOperand  = 0
		exOperator     = 1
		exRightOperand = 2
	)

	tokensInExpression := make([]Token, 3)
	currentTokenType := 0

	for i, token := range tokens {
		//Операторы могут стоять только на нечетных позициях (индексы начинаются с 0). Ситуацию когда выражение начинается с "-" уже обработали при парсинге
		if i%2 == 0 && token.GetType() == TTOperator {
			return res, errors.New("Два оператора подряд недопустимы")
		}
		//Проверим на всякий случай что в нечетной позиции оператор (хотя после парсинга это условие никогда не будет выполнятся)
		if i%2 == 1 && token.GetType() != TTOperator {
			return res, errors.New("Пропущен оператор")
		}
		if currentTokenType > 2 {
			if token.(OperatorToken).Operator() > tokensInExpression[TTOperator].(OperatorToken).Operator() {
				res = append(res, tokensInExpression[:2]...)
				tokensInExpression[exLeftOperand] = tokensInExpression[exRightOperand]
				tokensInExpression[exOperator] = token
				currentTokenType = exRightOperand
				continue
			}
			expression := NewExpressionToken(
				tokensInExpression[exOperator].(OperatorToken).Operator(),
				tokensInExpression[exLeftOperand].(ValueToken),
				tokensInExpression[exRightOperand].(ValueToken),
			)
			res = append(res, expression)
			res = append(res, token)
			currentTokenType = exLeftOperand
			continue
		}

		if currentTokenType == exOperator && len(res) > 0 && res[len(res)-1].(OperatorToken).Operator() > token.(OperatorToken).Operator() {
			res = append(res, tokensInExpression[exLeftOperand])
			res = append(res, token)
			currentTokenType = exLeftOperand
			continue
		}

		tokensInExpression[currentTokenType] = token
		currentTokenType++

	}

	switch currentTokenType - 1 {
	case exLeftOperand:
		res = append(res, tokensInExpression[exLeftOperand])
	case exOperator:
		return res, errors.New("Незаконченное выражение")
	case exRightOperand:
		expression := NewExpressionToken(
			tokensInExpression[exOperator].(OperatorToken).Operator(),
			tokensInExpression[exLeftOperand].(ValueToken),
			tokensInExpression[exRightOperand].(ValueToken),
		)
		res = append(res, expression)
	}

	return res, nil
}
