package main

import (
	"flag"
	"log"
	"os"
	"path"
	"time"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/fix50sp2/newordersingle"
	"github.com/quickfixgo/quickfix"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

type Sender struct {
}

// OnCreate implemented as part of Application interface.
func (r Sender) OnCreate(sessionID quickfix.SessionID) { return }

// OnLogon implemented as part of Application interface.
func (r Sender) OnLogon(sessionID quickfix.SessionID) { return }

// OnLogout implemented as part of Application interface.
func (r Sender) OnLogout(sessionID quickfix.SessionID) { return }

// ToAdmin implemented as part of Application interface.
func (r Sender) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) { return }

// ToApp implemented as part of Application interface.
func (r Sender) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) error {
	log.Printf("Sending %s\n", msg)
	return nil
}

// FromAdmin implemented as part of Application interface
func (r Sender) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}

// FromApp implemented as part of Application interface.
func (r Sender) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	log.Printf("FromApp: %s\n", msg.String())
	return nil
}

func main() {
	flag.Parse()

	cfgFileName := path.Join("config", "sender.cfg")
	if flag.NArg() > 0 {
		cfgFileName = flag.Arg(0)
	}

	cfg, err := os.Open(cfgFileName)
	if err != nil {
		log.Printf("Error opening %v, %v\n", cfgFileName, err)
		return
	}

	appSettings, err := quickfix.ParseSettings(cfg)
	if err != nil {
		log.Println("Error reading cfg,", err)
		return
	}

	app := Sender{}
	fileLogFactory, err := quickfix.NewFileLogFactory(appSettings)

	if err != nil {
		log.Println("Error creating file log factory,", err)
		return
	}

	initiator, err := quickfix.NewInitiator(app, quickfix.NewMemoryStoreFactory(), appSettings, fileLogFactory)
	if err != nil {
		log.Printf("Unable to create Initiator: %s\n", err)
		return
	}

	initiator.Start()

	for {
		clOrdID := field.NewClOrdID(uuid.NewV1().String())
		side := field.NewSide(enum.Side_BUY)
		transacttime := field.NewTransactTime(time.Now())
		ordtype := field.NewOrdType(enum.OrdType_MARKET)

		order := newordersingle.New(clOrdID, side, transacttime, ordtype)
		order.SetSenderCompID("Trader")
		order.SetSenderSubID("John Doe")
		order.SetTargetCompID("Broker")
		order.SetSymbol("GC_SEP18")
		order.SetOrderQty(decimal.NewFromFloat(23.14), 2)
		msg := order.ToMessage()

		err := quickfix.Send(msg)
		if err != nil {
			log.Println(err)
			break
		}
	}

	initiator.Stop()
}
