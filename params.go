package ringo

import "github.com/julienschmidt/httprouter"

type Params httprouter.Params

func (params *Params) ByName(name string) string {
	return httprouter.Params(*params).ByName(name)
}
