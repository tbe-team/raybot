package controller

type OptionFunc func(*controller)

func WithIDGenerator(fn func() string) OptionFunc {
	return func(c *controller) {
		c.genIDFunc = fn
	}
}
