package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"net/http"
	_ "net/http/pprof"

	"github.com/BurntSushi/toml"
	"github.com/immesys/wave/consts"
	"github.com/immesys/wave/waved"
	"github.com/immesys/wavemq/core"
	"github.com/immesys/wavemq/server"
	logging "github.com/op/go-logging"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli"
)

var lg = logging.MustGetLogger("main")

const WAVEMQPermissionSet = "\x4a\xd2\x3f\x5f\x6e\x73\x17\x38\x98\xef\x51\x8c\x6a\xe2\x7a\x7f\xcf\xf4\xfe\x9b\x86\xa3\xf1\xa2\x08\xc4\xde\x9e\xac\x95\x39\x6b"
const WAVEMQPublish = "publish"
const WAVEMQSubscribe = "subscribe"

type Configuration struct {
	RoutingConfig core.RoutingConfig
	WaveConfig    waved.Configuration
	QueueConfig   core.QManagerConfig
	LocalConfig   server.LocalServerConfig
	PeerConfig    server.PeerServerConfig
}

func main() {
	app := cli.NewApp()
	app.Name = "xbosd"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:   "run",
			Usage:  "Run WAVEMQ, WAVED daemons",
			Action: runDaemons,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "wavemq",
					Value: "wavemq.toml",
					Usage: "WAVEMQ configuration file",
				},
				cli.StringFlag{
					Name:  "waved",
					Value: "waved.toml",
					Usage: "WAVED configuration file",
				},
			},
		},
		{
			Name:   "att",
			Usage:  "Manage attestations",
			Action: runWaveAtt,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "db",
					EnvVar: "WAVEATTD_DB_LOCATION",
					Usage:  "Location of attestations database for this entity",
				},
				cli.StringFlag{
					Name:   "entity",
					EnvVar: "WAVE_DEFAULT_ENTITY",
					Usage:  "Entity granting attestations",
				},
				cli.StringFlag{
					Name:   "waved",
					EnvVar: "WAVE_AGENT",
					Value:  "localhost:410",
					Usage:  "WAVE Agent",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		lg.Fatal(err)
	}
}

func runWaveAtt(c *cli.Context) error {
	location := c.String("db")
	if location == "" {
		return fmt.Errorf("Set WAVEATTD_DB_LOCATION")
	}
	entity := c.String("entity")
	if entity == "" {
		return fmt.Errorf("Set WAVE_DEFAULT_ENTITY")
	}
	agent := c.String("waved")
	if agent == "" {
		return fmt.Errorf("Set WAVE_AGENT")
	}
	log.Info("╒ WAVEATTD_DB_LOCATION: ", location)
	log.Info("╞ WAVE_DEFAULT_ENTITY: ", entity)
	log.Info("╘ WAVE_AGENT: ", agent)
	cfg := &Config{
		Path:        location,
		Agent:       agent,
		Perspective: entity,
	}
	db, err := NewDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	go db.watch(".")

	db.RunShell()
	return nil
}

func runDaemons(c *cli.Context) error {
	fmt.Println(c.Args())
	wavemqfile := c.String("wavemq")
	fmt.Println("wavemq config: ", wavemqfile)

	wavedfile := c.String("waved")
	fmt.Println("waved config: ", wavedfile)

	sigchan := make(chan os.Signal, 30)
	wavemq_sigchan := make(chan os.Signal, 30)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// wavemq
	go func() {

		go func() {
			http.Handle("/metrics", promhttp.Handler())
			metricsAddr := "127.0.0.1:6060"
			if os.Getenv("METRICS_ADDRESS") != "" {
				metricsAddr = os.Getenv("METRICS_ADDRESS")
			}
			fmt.Printf("starting metrics on %q\n", metricsAddr)
			err := http.ListenAndServe(metricsAddr, nil)
			panic(err)
		}()

		var conf Configuration
		if _, err := toml.DecodeFile(wavemqfile, &conf); err != nil {
			fmt.Printf("failed to load configuration: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("configuration loaded\n")

		consts.DefaultToUnrevoked = conf.WaveConfig.DefaultToUnrevoked
		qm, err := core.NewQManager(&conf.QueueConfig)
		if err != nil {
			fmt.Printf("failed to initialize queues: %v\n", err)
			os.Exit(1)
		}
		am, err := core.NewAuthModule(&conf.WaveConfig)
		if err != nil {
			fmt.Printf("failed to initialize auth: %v\n", err)
			os.Exit(1)
		}
		tm, err := core.NewTerminus(qm, am, &conf.RoutingConfig)
		if err != nil {
			fmt.Printf("failed to initialize routing: %v\n", err)
			os.Exit(1)
		}
		server.NewLocalServer(tm, am, &conf.LocalConfig)
		server.NewPeerServer(tm, am, &conf.PeerConfig)
		//am.wave.StartServer(conf.WaveConfig.ListenIP, conf.WaveConfig.HTTPListenIP)
		<-wavemq_sigchan
		fmt.Printf("WAVEMQ SHUTTING DOWN\n")
		qm.Shutdown()
	}()

	// block for daemon
	sig := <-sigchan

	// kill other daemons
	wavemq_sigchan <- sig

	return nil
}
