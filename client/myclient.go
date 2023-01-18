package main

import (
	"context"
	"fmt"
	"log"
	"main/protofile"
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

// func main1() {
// 	fmt.Println("client")
// 	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
// 	handleError(err)
// 	defer conn.Close()

// 	c := protofile.NewUserServiceClient(conn)

// }
