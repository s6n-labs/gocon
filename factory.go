package gocon

import "reflect"

func Factory[T any](fn func(c Container) (T, error)) *Definition {
	rt := typeOf[T]()
	def := &Definition{
		Key:  keyOf(rt),
		Type: rt,
	}

	def.configureFunc = func(c Container) error {
		v, err := fn(c)
		if err != nil {
			return err
		}

		rv := reflect.ValueOf(v)
		def.Value = &rv

		return nil
	}

	return def
}
