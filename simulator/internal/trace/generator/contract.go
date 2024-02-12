package generator

import "github.com/maypok86/benchmarks/simulator/internal/event"

type parser interface {
	Parse(send func(event event.AccessEvent)) (bool, error)
}
