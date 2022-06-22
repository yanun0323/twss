package compare

type Compare int

const (
	Volume Compare = iota + 1
	VolumeMoney
	Start
	Max
	Min
	End
	Spread
	Per
)
