package internalgrpc

import (
	"log"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func loggingStreamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	var err error
	start := time.Now().UTC()

	ctx := ss.Context()
	p, _ := peer.FromContext(ctx)
	ip := strings.SplitN(p.Addr.String(), ":", 2)[0]
	req := info.FullMethod
	useragent := "UNKNOWN AGENT"
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if useragents, ok := md["user-agent"]; ok && len(useragents) > 0 {
			useragent = useragents[0]
		}
	}
	status := "OK"

	log.Printf("gRPC connect: %s %s \"%s\"\n",
		ip,
		req,
		useragent,
	)

	defer func() {
		stop := time.Now().UTC()
		log.Printf("gRPC disconn: %s %s %v %s (%v)\n",
			ip,
			req,
			stop.Sub(start),
			status,
			err,
		)
	}()

	err = handler(srv, ss)
	if err != nil {
		status = "FAIL"
	}

	return err
}
