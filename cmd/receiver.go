package cmd

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/micro/go-micro"
	"github.com/quickfixgo/quickfix"
	"github.com/rudeigerc/broker-gateway/mapper"
	"github.com/rudeigerc/broker-gateway/receiver"
	"github.com/spf13/cobra"
)

var receiverCmd = &cobra.Command{
	Use:   "receiver",
	Short: "Run receiver",
	Long:  "Run receiver",
	Run: func(cmd *cobra.Command, args []string) {
		cfgFileName := path.Join("config", "receiver.cfg")

		cfg, err := os.Open(cfgFileName)
		if err != nil {
			log.Printf("[cmd.receiver.receiverCmd] [ERROR] Error opening %v, %v\n", cfgFileName, err)
			return
		}

		appSettings, err := quickfix.ParseSettings(cfg)
		if err != nil {
			log.Println("[cmd.receiver.receiverCmd] [ERROR] Error reading cfg,", err)
			return
		}

		logFactory := quickfix.NewScreenLogFactory()

		r := receiver.NewReceiver()
		mapper.NewDB()

		defer func() {
			r.Stop()
			mapper.DB.Close()
		}()

		acceptor, err := quickfix.NewAcceptor(r, quickfix.NewMemoryStoreFactory(), appSettings, logFactory)
		if err != nil {
			log.Printf("[cmd.receiver.receiverCmd] [ERROR] Unable to create Acceptor: %s\n", err)
			return
		}

		err = acceptor.Start()
		defer acceptor.Stop()

		if err != nil {
			log.Printf("[cmd.receiver.receiverCmd] [ERROR] Unable to start Acceptor: %s\n", err)
			return
		}

		service := micro.NewService(
			micro.Name("github.com.rudeigerc.broker-gateway.receiver"),
			micro.RegisterTTL(time.Minute),
			micro.RegisterInterval(time.Second*30),
		)

		if err := service.Run(); err != nil {
			log.Fatalf("[cmd.receiver.receiverCmd] [FETAL] %s", err)
		}

	},
}
