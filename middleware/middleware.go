package middleware

import . "github.com/julienschmidt/httprouter"

type Func func(Handle) Handle
type MidCollect []Func

// Add add the new middleware.
// Newly added middleware will
// to be in the top level chain.
func (c MidCollect) Add(next Func) MidCollect {
	return append(c, next)
}

// Wrap wrap the handler in a chain
// middleware to appropriate this collection.
func (c MidCollect) Wrap(handle Handle) Handle {
	tmp := handle

	for i := len(c) - 1; i >= 0; i-- {
		tmp = c[i](tmp)
	}

	return tmp
}
