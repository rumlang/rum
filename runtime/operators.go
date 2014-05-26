package runtime

import "fmt"

// OpAdd implements the '+' function. It tries to determine automatically the
// type based on the first argument.
func OpAdd(values ...interface{}) interface{} {
	if len(values) < 1 {
		panic("Function '+' should take at least one argument")
	}

	switch values[0].(type) {
	case int64:
		var total int64
		for _, v := range values {
			total += v.(int64)
		}
		return total
	case float64:
		var total float64
		for _, v := range values {
			total += v.(float64)
		}
		return total
	default:
		panic(fmt.Sprintf("Unable to add values of type %T", values[0]))
	}
}

// OpAddInt64 implements '+int64' function.
func OpAddInt64(values ...int64) int64 {
	var total int64
	for _, v := range values {
		total += v
	}
	return total
}

// OpAddFloat64 implements '+float64' function.
func OpAddFloat64(values ...float64) float64 {
	var total float64
	for _, v := range values {
		total += v
	}
	return total
}

// OpMul implements the '*' function. It tries to determine automatically the
// type based on the first argument.
func OpMul(values ...interface{}) interface{} {
	if len(values) < 1 {
		panic("Function '*' should take at least one argument")
	}

	switch values[0].(type) {
	case int64:
		var total int64 = 1
		for _, v := range values {
			total *= v.(int64)
		}
		return total
	case float64:
		var total float64 = 1.0
		for _, v := range values {
			total *= v.(float64)
		}
		return total
	default:
		panic(fmt.Sprintf("Unable to add values of type %T", values[0]))
	}
}

// OpMulInt64 implements '*int64' function.
func OpMulInt64(values ...int64) int64 {
	var total int64 = 1.0
	for _, v := range values {
		total *= v
	}
	return total
}

// OpMulFloat64 implements '*float64' function.
func OpMulFloat64(values ...float64) float64 {
	var total float64 = 1.0
	for _, v := range values {
		total *= v
	}
	return total
}
