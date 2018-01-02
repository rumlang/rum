package runtime

import "errors"

var (
	errWrongNumberPar = errors.New("Wrong number of parameters")
	errParameterType  = errors.New("Error in parameter type")
)

//Adapter is an function type can be used to change/check params to registred golang functions
type Adapter func(values ...interface{}) ([]interface{}, error)

//CheckArity return an function to check if the arity of function call is n
func CheckArity(n int) Adapter {
	return func(values ...interface{}) ([]interface{}, error) {
		if len(values) != n {
			return []interface{}{}, errWrongNumberPar
		}
		return values, nil
	}
}

//ParamToFloat64 return an Adapter to convert the p nth param to float64 type
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
			return []interface{}{}, errParameterType
		}
	}
}

//ParamToInt64 return an Adapter to convert the p nth param to int64 type
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
			return []interface{}{}, errParameterType
		}
	}
}
