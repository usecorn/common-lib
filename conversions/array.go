package conversions

func MapValsToArr[T comparable, U any](m map[T]U) []U {
	out := make([]U, 0, len(m))
	for _, v := range m {
		out = append(out, v)
	}
	return out
}
