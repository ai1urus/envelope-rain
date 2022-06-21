package util

func Min(a, b interface{}) interface{} {
	var _a, _b int64

	switch a.(type) {
	case int:
		_a = int64(a.(int))
	case int32:
		_a = int64(a.(int32))
	case int64:
		_a = int64(a.(int64))
	}

	switch b.(type) {
	case int:
		_b = int64(b.(int))
	case int32:
		_b = int64(b.(int32))
	case int64:
		_b = int64(b.(int64))
	}

	if _a < _b {
		return a
	}
	return b
}
