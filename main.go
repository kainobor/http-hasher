package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kainobor/http-hasher/src"
)

const (
	minArgAMount             = 2
	flagPrefix               = "-"
	parallelFlagName         = "parallel"
	defaultParallelReqAmount = 10
	maxParallelProcesses     = 1000
	defaultReqTimeoutSec     = 10
)

type (
	Config struct {
		parallel int
		sites    []string
	}
)

func main() {
	ctx, exit := context.WithCancel(context.Background())

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	cfg, err := config()
	if err != nil {
		panic(err)
	}

	w := src.NewIOWriter()
	h := src.NewHasher(src.WithMD5())
	g := src.NewHttpGetter(defaultReqTimeoutSec)
	p := src.NewProcessor(cfg.parallel, w, h, g)

	p.Start(ctx)
	for _, curSite := range cfg.sites {
		select {
		case <-done:
			exit()
		default:
			p.Process(curSite)
			continue
		}
		break
	}

	p.Shutdown(ctx)
}

func config() (Config, error) {
	var cfg Config
	flag.IntVar(&cfg.parallel, parallelFlagName, defaultParallelReqAmount, "How many parallel requests can be processed")
	flag.Parse()

	if len(os.Args) < minArgAMount {
		return Config{}, fmt.Errorf("at least %d arguments should be used", minArgAMount)
	}

	if cfg.parallel > maxParallelProcesses {
		return Config{}, fmt.Errorf("too much parallel processed asked: %d; max is %d", cfg.parallel, maxParallelProcesses)
	}

	cfg.sites = os.Args[1:]
	if os.Args[1] == flagPrefix+parallelFlagName {
		cfg.sites = cfg.sites[2:]
	}

	return cfg, nil
}
