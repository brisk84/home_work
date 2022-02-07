package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/brisk84/home_work/hw12_13_14_15_calendar/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	conn, err := grpc.Dial("localhost:4343", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCalendarClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, _ = c.AddEvent(ctx, &pb.Event{
		Id:           "1",
		Title:        "Event01 [grpc]",
		TimeStart:    timestamppb.New(time.Date(2022, 0o1, 23, 17, 45, 0o0, 0o0, time.Local)),
		TimeEnd:      timestamppb.New(time.Date(2022, 0o1, 23, 18, 0o0, 0o0, 0o0, time.Local)),
		Description:  "Description of event 01",
		UserId:       "123",
		NotifyBefore: timestamppb.New(time.Date(2022, 0o1, 23, 17, 30, 0o0, 0o0, time.Local)),
	})
	r, err := c.AddEvent(ctx, &pb.Event{
		Id:           "2",
		Title:        "Event02 [grpc]",
		TimeStart:    timestamppb.New(time.Date(2022, 0o1, 24, 17, 45, 0o0, 0o0, time.Local)),
		TimeEnd:      timestamppb.New(time.Date(2022, 0o1, 24, 18, 0o0, 0o0, 0o0, time.Local)),
		Description:  "Description of event 02",
		UserId:       "123",
		NotifyBefore: timestamppb.New(time.Date(2022, 0o1, 24, 17, 30, 0o0, 0o0, time.Local)),
	})
	fmt.Println(err, r)

	c.EditEvent(ctx, &pb.Event{
		Id:           "1",
		Title:        "Event01 [edited]",
		TimeStart:    timestamppb.New(time.Date(2022, 0o1, 23, 17, 45, 0o0, 0o0, time.Local)),
		TimeEnd:      timestamppb.New(time.Date(2022, 0o1, 23, 18, 0o0, 0o0, 0o0, time.Local)),
		Description:  "Description of event 01 [edited]",
		UserId:       "123",
		NotifyBefore: timestamppb.New(time.Date(2022, 0o1, 23, 17, 30, 0o0, 0o0, time.Local)),
	})

	ev, err := c.GetEvent(ctx, &pb.EventID{Id: "1"})
	fmt.Println(err, ev)

	c.DeleteEvent(ctx, &pb.EventID{Id: "1"})

	evs, err := c.ListEvents(ctx, &emptypb.Empty{})

	fmt.Println(evs, err)

}
