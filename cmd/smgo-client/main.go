package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	pb "github.com/gmghub/smgo/pkg/smgo/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	retcode := 0
	defer func() {
		os.Exit(retcode)
	}()

	log.Println("OS:", runtime.GOOS, runtime.GOARCH)

	var (
		optGrpcaddr   string
		optCollectors string
		optInterval   uint
		optPeriod     uint
	)

	flag.StringVar(&optCollectors, "c", "", "Specify certain collectors to be received")
	flag.UintVar(&optInterval, "i", 1, "Specify the interval to receive stats")
	flag.UintVar(&optPeriod, "p", 5, "Specify the period for averaging receiving stats")
	flag.StringVar(&optGrpcaddr, "s", "127.0.0.1:50051", "Specify <addr:port> of the server")
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	ctx, signalstop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer signalstop()

	// conn, err := grpc.Dial(optGrpcaddr, grpc.WithInsecure())
	conn, err := grpc.Dial(optGrpcaddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("can not connect to the server:", err)
		retcode = 1
		return
	}

	client := pb.NewSmgoServiceClient(conn)
	req := &pb.GetSysStatRequest{
		Statinterval: uint32(optInterval),
		Statperiod:   uint32(optPeriod),
		Collectors:   optCollectors,
	}
	stream, err := client.GetSysStat(context.Background(), req)
	if err != nil {
		log.Println("open stream error:", err)
		retcode = 1
		return
	}

	go func() {
		for {
			resp, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				log.Println("EOF received")
				signalstop()
				return
			}
			if err != nil {
				log.Println("cannot receive:", err)
				signalstop()
				return
			}
			log.Println("resp received:", string(resp.Sysstat))
		}
	}()

	<-ctx.Done()
	stream.CloseSend()
	log.Println("app exit")
}
