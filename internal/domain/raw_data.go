package domain

type Raw interface {
	GetData() (interface{}, error)
}

type RawData interface {
	IsOK() bool
	Parse() []interface{}
}
