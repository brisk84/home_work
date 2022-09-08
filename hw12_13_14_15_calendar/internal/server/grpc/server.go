package internalgrpc

import (
	"context"
	"net"

	pb "github.com/brisk84/home_work/hw12_13_14_15_calendar/api"
	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/app"
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
	appl       *app.App
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
	GetEventsOnDay(context.Context, string) ([]storage.Event, error)
	GetEventsOnWeek(context.Context, string) ([]storage.Event, error)
	GetEventsOnMonth(context.Context, string) ([]storage.Event, error)
}

func NewServer(logger Logger, app Application, addr string) *Server {
	server := &Server{
		addr: addr,
		logg: logger,
		appl: appl,
	}

	return server
}

func (s *Server) Start(ctx context.Context) error {
	s.logg.Info("Start grpc server")
	s.ctx = ctx

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
	err := s.appl.AddEvent(ctx, PbToStorage(event))
	return &pb.Error{}, err
}

func (s *Server) GetEvent(ctx context.Context, eventID *pb.EventID) (*pb.Event, error) {
	s.logg.Info("gprc: GetEvent")
	if eventID == nil {
		return nil, nil
	}
	ev, err := s.appl.GetEvent(ctx, eventID.Id)
	return StorageToPb(ev), err
}

func (s *Server) EditEvent(ctx context.Context, event *pb.Event) (*pb.Error, error) {
	s.logg.Info("gprc: EditEvent")
	err := s.appl.EditEvent(ctx, PbToStorage(event))
	return nil, err
}

func (s *Server) DeleteEvent(ctx context.Context, eventID *pb.EventID) (*pb.Error, error) {
	s.logg.Info("gprc: DeleteEvent")
	if eventID == nil {
		return nil, nil
	}
	err := s.appl.DeleteEvent(ctx, eventID.Id)
	return nil, err
}

func (s *Server) ListEvents(ctx context.Context, empty *emptypb.Empty) (*pb.Events, error) {
	s.logg.Info("gprc: ListEvents")
	evs, err := s.appl.ListEvents(ctx)
	res := []*pb.Event{}
	for _, ev := range evs {
		res = append(res, StorageToPb(ev))
	}
	events := &pb.Events{Events: res}
	return events, err
}

func (s *Server) GetEventsOnDay(ctx context.Context, day *pb.Day) (*pb.Events, error) {
	s.logg.Info("gprc: GetEventsOnDay")
	evs, err := s.app.GetEventsOnDay(ctx, day.Day)
	res := []*pb.Event{}
	for _, ev := range evs {
		res = append(res, StorageToPb(ev))
	}
	events := &pb.Events{Events: res}
	return events, err
}

func (s *Server) GetEventsOnWeek(ctx context.Context, day *pb.Day) (*pb.Events, error) {
	s.logg.Info("gprc: GetEventsOnWeek")
	evs, err := s.app.GetEventsOnWeek(ctx, day.Day)
	res := []*pb.Event{}
	for _, ev := range evs {
		res = append(res, StorageToPb(ev))
	}
	events := &pb.Events{Events: res}
	return events, err
}

func (s *Server) GetEventsOnMonth(ctx context.Context, day *pb.Day) (*pb.Events, error) {
	s.logg.Info("gprc: GetEventsOnMonth")
	evs, err := s.app.GetEventsOnMonth(ctx, day.Day)
	res := []*pb.Event{}
	for _, ev := range evs {
		res = append(res, StorageToPb(ev))
	}
	events := &pb.Events{Events: res}
	return events, err
}
