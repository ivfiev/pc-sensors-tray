package types

type Result interface {
	Value() float64
	IsClose(result Result) bool

	Icon() string
	Lines() []string
	Colour() string
	Intensity() float64
}
