package api

import "strconv"

type TaskToken struct {
	Value      int  //Значение: для операторов 0 - "+", 1 - "-", 2 - "*", 3 - "/". Для операндов - собственно значение
	IsOperator bool //Если true, то это оператор (+,-,*,/), иначе это число
}

func Parse(text string) ([]TaskToken, error) {
	var res []TaskToken
	plus := rune('+')
	minus := rune('-')
	mul := rune('*')
	divide := rune('/')

	var buf []rune
	for _, r := range text {
		if r == plus || r == minus || r == mul || r == divide {
			if len(buf) > 0 {
				val, err := strconv.Atoi(string(buf))
				if err != nil {
					return res, err
				}
				res = append(res, TaskToken{Value: val, IsOperator: false})
				buf = buf[:0]
			}
			operator := 0
			switch r {
			case plus:
				operator = 0
			case minus:
				operator = 1
			case mul:
				operator = 2
			case divide:
				operator = 3
			}
			res = append(res, TaskToken{Value: operator, IsOperator: true})
		} else {
			buf = append(buf, r)
		}
	}
	if len(buf) > 0 {
		val, err := strconv.Atoi(string(buf))
		if err != nil {
			return res, err
		}
		res = append(res, TaskToken{Value: val, IsOperator: false})
		buf = buf[:0]
	}
	return res, nil
}
