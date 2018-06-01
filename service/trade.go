package service

import (
	"log"

	"github.com/rudeigerc/broker-gateway/mapper"
	"github.com/rudeigerc/broker-gateway/model"
)

type Trade struct {
}

func (t Trade) NewTrade(trade model.Trade) {
	m := mapper.NewMapper()
	err := m.Create(&trade)
	if err != nil {
		log.Fatalf("[service.order.NewOrder] [FETAL] %s", err)
	}
}
