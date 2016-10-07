package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/mod/loxybeat/config"
	"github.com/mod/loxybeat/httpd"
)

type Loxybeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
	beat   *beat.Beat
	server *httpd.Server
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Loxybeat{
		done:   make(chan struct{}),
		config: config,
		server: httpd.New(&config),
	}
	go bt.listenPipe()
	return bt, nil
}

func (bt *Loxybeat) Run(b *beat.Beat) error {
	logp.Info("loxybeat is running! Hit CTRL-C to stop it.")
	bt.beat = b
	bt.client = b.Publisher.Connect()
	bt.server.Start()
	return nil
}

func (bt *Loxybeat) listenPipe() {
	for {
		select {
		case payload := <-bt.server.Pipe:
			logp.Debug("beater", "Payload recv: %s", payload)

			event := common.MapStr{
				"@timestamp": common.Time(time.Now()),
				"type":       bt.beat.Name,
				"payload":    payload,
			}
			bt.client.PublishEvent(event)
			logp.Info("Event sent")

		case <-bt.done:
			logp.Info("listenPipe() returning")
			return
		}
	}
}

func (bt *Loxybeat) Stop() {
	logp.Info("beater.Stop() Recieved")
	bt.server.Stop()
	bt.client.Close()
	close(bt.done)
}
