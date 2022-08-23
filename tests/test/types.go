package test

import "net/http"

type Test interface {
	Run(host string) []error
	Handler() http.HandlerFunc
}

type Factory func() []Test
