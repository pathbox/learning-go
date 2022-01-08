package filter

import (
	"context"
	"fmt"
	"strings"
)

type Context struct {
	context.Context
	key string
}

type HandlerFunc func(ctx Context, params ...interface{}) error

type HandlersChain []HandlerFunc

type Filter interface {
	DoFilter() error
}

type FilterImpl struct {
	Ctx   Context
	Chain HandlersChain
}

func (myFilter *FilterImpl) DoFilter() error {
	for _, handler := range myFilter.Chain {
		err := handler(myFilter.Ctx)
		if err != nil {
			fmt.Println(fmt.Sprintf("filter"))
			return err
		}
	}
	return nil
}
