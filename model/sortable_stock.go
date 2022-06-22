package model

import (
	"main/model/compare"
	"time"
)

type SortableStock struct {
	Stokes  []Stock
	date    time.Time
	Compare compare.Compare
}

func NewSortableStock(stokes []Stock, date time.Time, compare compare.Compare) SortableStock {
	return SortableStock{
		Stokes:  stokes,
		date:    date,
		Compare: compare,
	}
}

func (s *SortableStock) Len() int {
	return len(s.Stokes)
}

func (s *SortableStock) Less(i, j int) bool {
	prev := s.Stokes[i].Deals[s.date]
	post := s.Stokes[j].Deals[s.date]
	result := false
	switch s.Compare {
	case compare.Volume:
		result = prev.Volume > post.Volume
	case compare.VolumeMoney:
		result = prev.VolumeMoney > post.VolumeMoney
	case compare.Start:
		result = prev.Start > post.Start
	case compare.Max:
		result = prev.Max > post.Max
	case compare.Min:
		result = prev.Min > post.Min
	case compare.End:
		result = prev.End > post.End
	case compare.Spread:
		result = prev.Spread > post.Spread
	case compare.Per:
		result = prev.Per < post.Per
	}
	return result
}

func (s *SortableStock) Swap(i, j int) {
	temp := s.Stokes[i]
	s.Stokes[i] = s.Stokes[j]
	s.Stokes[j] = temp
}
