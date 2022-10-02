package internalgrpc

import (
	"strings"
	"time"

	pb "github.com/gmghub/smgo/pkg/smgo/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	pb.UnimplementedSmgoServiceServer
	app  Application
	quit chan struct{}
}

func NewService(app Application) *Service {
	return &Service{
		app:  app,
		quit: make(chan struct{}),
	}
}

func (s *Service) Close() {
	close(s.quit)
}

func (s *Service) GetSysStat(req *pb.GetSysStatRequest, srv pb.SmgoService_GetSysStatServer) (err error) {
	// check collectors list, remove unknown, log errors and run
	var collectors []string
	if req.Collectors != "" {
		collectors = strings.Split(req.Collectors, ",")
		collectorsmap := s.app.MapCollectors()
		for _, c := range collectors {
			if _, ok := collectorsmap[c]; !ok {
				err = status.Error(codes.NotFound, "unknown collector name "+c)
				return err
			}
		}
	}
	ticker := time.NewTicker(time.Duration(req.Statinterval) * time.Second)
out:
	for {
		select {
		case <-s.quit:
			break out
		case <-ticker.C:
			stat := s.app.StatJSON(int(req.Statperiod), collectors...)
			if err = srv.Send(&pb.GetSysStatResponse{Sysstat: stat}); err != nil {
				break out
			}
		}
	}
	return err
}
