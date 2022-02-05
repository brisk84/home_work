package internalgrpc

import (
	"context"
	"net"

	"github.com/brisk84/home_work/hw12_13_14_15_calendar/api"
	pb "github.com/brisk84/home_work/hw12_13_14_15_calendar/api"
	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedCalendarServer
	addr       string
	grpcServer *grpc.Server
	logg       Logger
	app        Application
	ctx        context.Context
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Application interface {
	AddEvent(context.Context, storage.Event) error
	GetEvent(context.Context, string) (storage.Event, error)
	EditEvent(context.Context, storage.Event) error
	DeleteEvent(context.Context, string) error
	ListEvents(context.Context) ([]storage.Event, error)
}

func NewServer(logger Logger, app Application, addr string) *Server {
	server := &Server{
		addr: addr,
		logg: logger,
		app:  app,
	}

	return server
}

func (s *Server) Start(ctx context.Context) error {
	s.logg.Info("Start grpc server")

	loggingInterceptor := grpc.ChainUnaryInterceptor(s.loggingMiddleware)
	s.grpcServer = grpc.NewServer(loggingInterceptor)

	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		s.logg.Error(err.Error())
		return err
	}
	pb.RegisterCalendarServer(s.grpcServer, s)
	err = s.grpcServer.Serve(lis)
	<-ctx.Done()
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	s.logg.Info("Stop grpc server")
	s.grpcServer.Stop()
	return nil
}

func PbToStorage(event *pb.Event) storage.Event {
	return storage.Event{
		ID:           event.Id,
		Title:        event.Title,
		TimeStart:    event.TimeStart.AsTime(),
		TimeEnd:      event.TimeEnd.AsTime(),
		Description:  event.Description,
		UserID:       event.UserId,
		NotifyBefore: event.NotifyBefore.AsTime(),
	}
}

func StorageToPb(event storage.Event) *pb.Event {
	return &pb.Event{
		Id:           event.ID,
		Title:        event.Title,
		TimeStart:    timestamppb.New(event.TimeStart),
		TimeEnd:      timestamppb.New(event.TimeEnd),
		Description:  event.Description,
		UserId:       event.UserID,
		NotifyBefore: timestamppb.New(event.NotifyBefore),
	}
}

func (s *Server) AddEvent(ctx context.Context, event *pb.Event) (*pb.Error, error) {
	s.logg.Info("gprc: AddEvent")
	err := s.app.AddEvent(ctx, PbToStorage(event))
	return &pb.Error{}, err
}

func (s *Server) GetEvent(ctx context.Context, eventId *pb.EventID) (*pb.Event, error) {
	s.logg.Info("gprc: GetEvent")
	if eventId == nil {
		return nil, nil
	}
	ev, err := s.app.GetEvent(ctx, eventId.Id)
	return StorageToPb(ev), err
}

func (s *Server) EditEvent(ctx context.Context, event *pb.Event) (*pb.Error, error) {
	s.logg.Info("gprc: EditEvent")
	err := s.app.EditEvent(ctx, PbToStorage(event))
	return nil, err
}

func (s *Server) DeleteEvent(ctx context.Context, eventId *pb.EventID) (*pb.Error, error) {
	s.logg.Info("gprc: DeleteEvent")
	if eventId == nil {
		return nil, nil
	}
	err := s.app.DeleteEvent(ctx, eventId.Id)
	return nil, err
}

func (s *Server) ListEvents(ctx context.Context, empty *emptypb.Empty) (*api.Events, error) {
	s.logg.Info("gprc: ListEvents")
	evs, err := s.app.ListEvents(ctx)
	res := []*pb.Event{}
	for _, ev := range evs {
		res = append(res, StorageToPb(ev))
	}
	events := &api.Events{Events: res}
	return events, err
}
