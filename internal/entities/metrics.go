package entities

type Gauge float64
type Counter int64

type Metrics struct {
	Type  string
	Name  string
	Value float64
}
