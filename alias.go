package gocon

func Alias[T any](alias string) *Definition {
	return &Definition{
		Key:         alias,
		Type:        typeOf[T](),
		Tags:        nil,
		resolveFunc: resolve[T],
	}
}
