/*
@Time : 2019/4/19 16:24 
@Author : kenny zhu
@File : options.go
@Software: GoLand
@Others:
*/
package web

import (
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/selector"
)

type Options struct {
	Registry registry.Registry
	Selector selector.Strategy
}

type Option func(*Options)

func WithRegistry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

func WithSelector(s selector.Strategy) Option {
	return func(o *Options) {
		o.Selector = s
	}
}