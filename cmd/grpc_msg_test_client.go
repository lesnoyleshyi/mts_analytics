package main

import (
	"context"
	"github.com/google/uuid"
	pb "gitlab.com/g6834/team17/grpc/analytics_messaging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

const gRPCtarget = `localhost:50051`

func main() {
	sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	ctx := sigCtx
	//ctx, cancel := context.WithTimeout(sigCtx, time.Second*10)
	//defer cancel()

	conn, err := grpc.DialContext(sigCtx, gRPCtarget, grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Can't connect:", err)
	}

	client := pb.NewAnalyticsMsgServiceClient(conn)

	for {
		ts := time.Now().Add(time.Second * time.Duration(rand.Intn(1000)))
		msg := pb.EventMessage{
			TaskUuid:  uuid.NewString(),
			EventType: pb.EventType(rand.Intn(7)),
			UserUuid:  uuid.NewString(),
			Timestamp: timestamppb.New(ts),
		}

		resp, err := client.SendMessage(ctx, &msg)
		if err != nil {
			log.Println(err, "\tresp.Success:", false)
		} else {
			log.Println("resp.Success:", resp.Success)
		}

		select {
		case <-ctx.Done():
			log.Println(conn.Close())
			return
		default:
		}
		time.Sleep(time.Second * 1)
	}

}
