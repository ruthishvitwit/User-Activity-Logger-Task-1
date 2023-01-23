package main

import (
	"fmt"
	"main/protofile"
	"testing"

	"google.golang.org/grpc"
)

// var c activity_pb.UserServiceClient

// func connectToServer() {

// }

func check_string(t *testing.T, got string, want string) {
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func check_bool(t *testing.T, got bool, want bool) {
	if got != want {
		s := fmt.Sprint("got", got, ", wanted", want)
		t.Errorf(s)
	}
}

func connectToServer() (protofile.UserServiceClient, *grpc.ClientConn) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	handleError(err)
	// defer conn.Close()
	c := protofile.NewUserServiceClient(conn)
	return c, conn
}

func TestUserAdd(t *testing.T) {
	c, conn := connectToServer()
	defer conn.Close()

	got := UserData(c, "testuser1", "testuser1@gmail.com", 1212121212)
	want := "User already exist"
	check_string(t, got, want)
	got = UserData(c, "testuser2", "testuser2@gmail.com", 13131313)
	want = "User already exist"
	check_string(t, got, want)
}

func TestActivityAdd(t *testing.T) {
	c, conn := connectToServer()
	defer conn.Close()

	got := ActivityAdd(c, "testuser1@gmail.com", "Sleep", 7, "label2")
	want := "User activity added"
	check_string(t, got, want)
}

func TestGetUser(t *testing.T) {
	c, conn := connectToServer()
	defer conn.Close()

	got := GetUser(c, "testuser1@gmail.com")
	want := true
	check_bool(t, got, want)
	got = GetUser(c, "unknownuser@gmail.com")
	want = false
	check_bool(t, got, want)

}

func TestGetActivity(t *testing.T) {
	c, conn := connectToServer()
	defer conn.Close()

	got := GetActivity(c, "sai@gmail.com")
	want := true

	check_bool(t, got, want)
}
