package model

type RawJson struct {
	State   string     `json:"stat"`
	Data8   [][]string `json:"data8"`
	Data9   [][]string `json:"data9"`
	Fields9 []string   `json:"fields9"`
}
