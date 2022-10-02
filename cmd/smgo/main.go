package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/gmghub/smgo/internal/app"
	"github.com/gmghub/smgo/internal/collector"
	internalgrpc "github.com/gmghub/smgo/internal/server/grpc"
)

const (
	shutdownTimeout   = 5 * time.Second
	defaultBufferSize = 60
)

func main() {
	retcode := 0
	defer func() {
		os.Exit(retcode)
	}()

	log.Println("OS:", runtime.GOOS, runtime.GOARCH)

	var (
		optCollectors string
		optGrpcaddr   string
	)
	// flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
	flag.StringVar(&optGrpcaddr, "l", "127.0.0.1:50051", "Specify the address for gRPC server to listen")
	flag.StringVar(&optCollectors, "c", "", "Specify the certain collectors to be run")
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	collectors := map[string]int{
		"sysstat":  1,
		"cpustat":  1,
		"diskstat": 0,
	}

	if optCollectors != "" {
		optCollectorsList := strings.Split(optCollectors, ",")
		for k := range collectors {
			collectors[k] = 0
		}
		for _, v := range optCollectorsList {
			if _, ok := collectors[v]; ok {
				collectors[v] = 1
			}
		}
	}

	if collectors[collector.CollectorNameDiskStat] != 0 {
		log.Println(collector.CollectorNameDiskStat, "is not supported on", runtime.GOOS)
	}

	app := app.NewApp()

	ctx, signalstop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer signalstop()

	for k, v := range collectors {
		switch {
		case k == "sysstat" && v > 0:
			log.Println("Adding collector:", k)
			app.Add(collector.NewSysStatCollector(defaultBufferSize))
		case k == "cpustat" && v > 0:
			log.Println("Adding collector:", k)
			app.Add(collector.NewCPUStatCollector(defaultBufferSize))
		case k == "diskstat" && v > 0:
			log.Println("Adding collector:", k)
			app.Add(collector.NewDiskStatCollector(defaultBufferSize))
		default:
			log.Println("Skipping collector:", k)
		}
	}
	if len(app.MapCollectors()) == 0 {
		log.Println("main:", "error running server with no collectors")
		retcode = 1
		return
	}
	app.Start()

	grpcserver := internalgrpc.NewServer(app)
	log.Println("gRPC server running on ", optGrpcaddr)
	grpcserver.Start(optGrpcaddr)

	<-ctx.Done()

	ctx, shutdownstop := context.WithTimeout(context.Background(), shutdownTimeout)
	go func() {
		log.Println("stopping gRPC server")
		if err := grpcserver.Stop(); err != nil {
			log.Println(err)
		}
		log.Println("stopping collectors")
		app.Close()
		shutdownstop()
	}()

	<-ctx.Done()
	log.Println("app exit")
}
