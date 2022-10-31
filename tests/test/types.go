package test

import "net/http"

type Test interface {
	Name() string
	Run(host string) []error
	Handler() http.HandlerFunc
}

type Factory func() []Test
