package toolkit

func Ref[T any](v T) *T {
	return &v
}
