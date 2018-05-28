package matcher

import (
	"github.com/rudeigerc/broker-gateway/model"
	"github.com/shopspring/decimal"
)

type Level struct {
	Price decimal.Decimal
	Order []*model.Order
}

type LevelHeap []Level

func (h LevelHeap) Len() int           { return len(h) }
func (h LevelHeap) Less(i, j int) bool { return h[i].Price.Cmp(h[j].Price) < 0 }
func (h LevelHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *LevelHeap) Push(x interface{}) {
	*h = append(*h, x.(Level))
}

func (h *LevelHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}