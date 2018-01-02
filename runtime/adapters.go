package runtime

import (
	"errors"
	"fmt"
)

var (
	errWrongNumberPar = errors.New("Wrong number of parameters")
	errParameterType  = errors.New("Error in parameter type")
)

type Adapter func(values ...interface{}) ([]interface{}, error)

func CheckArity(n int) Adapter {
	return func(values ...interface{}) ([]interface{}, error) {
		if len(values) != n {
			return []interface{}{}, errWrongNumberPar
		}
		return values, nil
	}
}

func ParamToFloat64(p int) Adapter {
	return func(values ...interface{}) ([]interface{}, error) {
		switch values[p].(type) {
		case int64:
			v := values[p].(int64)
			values[p] = float64(v)
			return values, nil
		case float64:
			return values, nil
		default:
			panic(fmt.Sprintf("Unable to compare values of type %T", values[0]))
			return []interface{}{}, errWrongNumberPar
		}
	}
}

func ParamToInt64(p int) Adapter {
	return func(values ...interface{}) ([]interface{}, error) {
		switch values[p].(type) {
		case int64:
			return values, nil
		case float64:
			v := values[p].(float64)
			values[p] = int64(v)
			return values, nil
		default:
			panic(fmt.Sprintf("Unable to compare values of type %T", values[0]))
			return []interface{}{}, errWrongNumberPar
		}
	}
}
