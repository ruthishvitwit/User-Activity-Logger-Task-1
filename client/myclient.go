package main

import (
	"context"
	"fmt"
	"log"
	"main/protofile"
	"time"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func UserData(c protofile.UserServiceClient, name string, email string, phone int64) {
	UserRequest := protofile.UserRequest{
		User: &protofile.User{
			Name:  name,
			Email: email,
			Phone: phone,
		},
	}
	res, err := c.UserData(context.Background(), &UserRequest)
	handleError(err)
	fmt.Println(res)
}
func getTimeStamp() string {
	t := time.Now()
	ts := t.Format("01-02-2006 15:04:05 Monday")
	return ts
}
func ActData(c protofile.UserServiceClient, email string, at string, duration int32, label string) {
	t := time.Now()
	ts := t.Format("01-02-2006 15:04:05 Monday")
	activityAddRequest := protofile.ActRequest{
		Activity: &protofile.Activity{
			Activitytype: at,
			Timestamp:    ts,
			Duration:     duration,
			Label:        label,
			Email:        email,
		},
	}

	res, err := c.ActData(context.Background(), &activityAddRequest)
	handleError(err)
	fmt.Println(res)
}

// func main1() {
// 	fmt.Println("client")
// 	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
// 	handleError(err)
// 	defer conn.Close()

// 	c := protofile.NewUserServiceClient(conn)

// }
