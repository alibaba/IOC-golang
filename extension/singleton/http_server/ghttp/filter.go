package ghttp

type HandleFunc func(controller *GRegisterController) (err error)

type Filter func(controller *GRegisterController, f HandleFunc) (err error)

func NoopFilter(controller *GRegisterController, f HandleFunc) (err error) {
	return f(controller)
}

type Chain []Filter

func (fc *Chain) Handle(controller *GRegisterController, f HandleFunc) (err error) {

	n := len(*fc)

	if n == 0 {
		return NoopFilter(controller, f)
	}

	if n == 1 {
		return (*fc)[0](controller, f)
	}

	lastI := n - 1
	return func(controller *GRegisterController, f HandleFunc) error {

		var (
			chainFunc HandleFunc
			curI      int
		)
		chainFunc = func(controller *GRegisterController) error {
			if curI == lastI {
				return f(controller)
			}
			curI++
			err := (*fc)[curI](controller, chainFunc)
			curI--
			return err
		}
		return (*fc)[0](controller, chainFunc)
	}(controller, f)
}

func (fc *Chain) AddFilter(f []Filter) {
	if f == nil {
		return
	}
	*fc = append(*fc, f...)
}
