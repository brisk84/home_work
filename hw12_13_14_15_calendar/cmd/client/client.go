package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/brisk84/home_work/hw12_13_14_15_calendar/api"
	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func intTests(ctx context.Context, c pb.CalendarClient) int {
	errUUIDBusy := storage.ErrUUIDBusy
	errNotFound := storage.ErrNotFound

	_, _ = c.DeleteEvent(ctx, &pb.EventID{Id: "1"})
	_, _ = c.DeleteEvent(ctx, &pb.EventID{Id: "2"})
	_, _ = c.DeleteEvent(ctx, &pb.EventID{Id: "3"})

	_, err := c.AddEvent(ctx, &pb.Event{
		Id:           "1",
		Title:        "Event01 [grpc]",
		TimeStart:    timestamppb.New(time.Date(2022, 0o1, 23, 17, 45, 0o0, 0o0, time.Local)),
		TimeEnd:      timestamppb.New(time.Date(2022, 0o1, 23, 18, 0o0, 0o0, 0o0, time.Local)),
		Description:  "Description of event 01",
		UserId:       "123",
		NotifyBefore: timestamppb.New(time.Date(2022, 0o1, 23, 17, 30, 0o0, 0o0, time.Local)),
	})
	if err != nil {
		return 1
	}

	_, err = c.AddEvent(ctx, &pb.Event{
		Id:           "2",
		Title:        "Event02 [grpc]",
		TimeStart:    timestamppb.New(time.Date(2022, 0o1, 24, 17, 45, 0o0, 0o0, time.Local)),
		TimeEnd:      timestamppb.New(time.Date(2022, 0o1, 24, 18, 0o0, 0o0, 0o0, time.Local)),
		Description:  "Description of event 02",
		UserId:       "123",
		NotifyBefore: timestamppb.New(time.Now()),
	})
	if err != nil {
		return 2
	}

	_, err = c.AddEvent(ctx, &pb.Event{
		Id:           "2",
		Title:        "Event02 [grpc]",
		TimeStart:    timestamppb.New(time.Date(2022, 0o1, 24, 17, 45, 0o0, 0o0, time.Local)),
		TimeEnd:      timestamppb.New(time.Date(2022, 0o1, 24, 18, 0o0, 0o0, 0o0, time.Local)),
		Description:  "Description of event 02",
		UserId:       "123",
		NotifyBefore: timestamppb.New(time.Now()),
	})
	if !errors.As(err, &errUUIDBusy) {
		return 3
	}

	_, err = c.AddEvent(ctx, &pb.Event{
		Id:           "3",
		Title:        "Event03 [grpc]",
		TimeStart:    timestamppb.New(time.Date(2022, 0o1, 25, 17, 45, 0o0, 0o0, time.Local)),
		TimeEnd:      timestamppb.New(time.Date(2022, 0o1, 25, 18, 0o0, 0o0, 0o0, time.Local)),
		Description:  "Description of event 03",
		UserId:       "123",
		NotifyBefore: timestamppb.New(time.Now()),
	})
	if err != nil {
		return 4
	}

	_, err = c.GetEvent(ctx, &pb.EventID{Id: "1"})
	if err != nil {
		return 5
	}

	_, err = c.GetEvent(ctx, &pb.EventID{Id: "10"})
	if !errors.As(err, &errNotFound) {
		return 6
	}

	_, err = c.EditEvent(ctx, &pb.Event{
		Id:           "1",
		Title:        "Event01 [edited]",
		TimeStart:    timestamppb.New(time.Date(2022, 0o1, 23, 17, 45, 0o0, 0o0, time.Local)),
		TimeEnd:      timestamppb.New(time.Date(2022, 0o1, 23, 18, 0o0, 0o0, 0o0, time.Local)),
		Description:  "Description of event 01 [edited]",
		UserId:       "123",
		NotifyBefore: timestamppb.New(time.Now()),
	})
	if err != nil {
		return 7
	}

	evs, err := c.ListEvents(ctx, &emptypb.Empty{})
	if err != nil || len(evs.Events) != 3 {
		return 8
	}

	evs, err = c.GetEventsOnDay(ctx, &pb.Day{Day: "2022-01-25"})
	if err != nil || len(evs.Events) != 1 {
		return 9
	}

	evs, err = c.GetEventsOnWeek(ctx, &pb.Day{Day: "2022-01-24"})
	if err != nil || len(evs.Events) != 2 {
		return 10
	}

	evs, err = c.GetEventsOnMonth(ctx, &pb.Day{Day: "2022-01-25"})
	if err != nil || len(evs.Events) != 3 {
		return 11
	}

	evs, err = c.GetEventsOnMonth(ctx, &pb.Day{Day: "2022-02-25"})
	if err != nil || len(evs.Events) != 0 {
		return 12
	}

	_, err = c.DeleteEvent(ctx, &pb.EventID{Id: "1"})
	if err != nil {
		return 13
	}

	return 0
}

func main() {
	conn, err := grpc.Dial("localhost:4343", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := pb.NewCalendarClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	exitCode := intTests(ctx, c)
	fmt.Println("ExitCode:", exitCode)

	cancel()
	conn.Close()
	os.Exit(exitCode)
}
