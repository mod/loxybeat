package beater

import (
	"fmt"
	"html"
	"net/http"
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

var instance *Loxybeat

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
	bt.server.SetHandleFunc("/log", bt.handleLog)
	instance = bt
	return bt, nil
}

func (bt *Loxybeat) Run(b *beat.Beat) error {
	logp.Info("loxybeat is running! Hit CTRL-C to stop it.")
	bt.beat = b
	bt.client = b.Publisher.Connect()
	bt.server.Start()
	return nil
}

func (bt *Loxybeat) handleLog(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q\n", html.EscapeString(r.URL.Path))
	event := common.MapStr{
		"@timestamp": common.Time(time.Now()),
		"type":       bt.beat.Name,
		"counter":    1,
	}
	bt.client.PublishEvent(event)
	logp.Info("Event sent")
}

func (bt *Loxybeat) Stop() {
	logp.Info("beater.Stop() Recieved")
	bt.server.Stop()
	bt.client.Close()
	close(bt.done)
}
