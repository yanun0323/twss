package model

type StockMap map[string]string

func (m StockMap) List() StockList {
	list := make(StockList, 0, len(m))
	for id, name := range m {
		list = append(list, StockListUnit{
			ID:   id,
			Name: name,
		})
	}
	return list
}
