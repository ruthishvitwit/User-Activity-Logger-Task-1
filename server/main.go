package main

import (
	"context"
	"fmt"
	"log"
	"main/protofile"
	"net"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type user_item struct {
	Id    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name"`
	Email string             `bson:"email"`
	Phone int64              `bson:"phone"`
}
type activity_item struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	Activitytype string             `bson:"activitytype"`
	Duration     int32              `bson:"duration"`
	Label        string             `bson:"label"`
	Timestamp    string             `bson:"timestamp"`
	Email        string             `bson:"email"`
}

func pushUserToDb(ctx context.Context, item user_item) string {
	email := item.Email
	filter := bson.M{
		"email": email,
	}

	var result_data []user_item
	cursor, err := U_collection.Find(context.TODO(), filter)
	handleError(err)

	cursor.All(context.Background(), &result_data)

	if len(result_data) != 0 {
		result := "User already exist"
		return result
	}

	U_collection.InsertOne(ctx, item)
	result := "User created"
	return result
}
func pushActToDb(ctx context.Context, item activity_item) string {
	A_collection.InsertOne(ctx, item)
	result := "User activity added"
	return result
}
func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type server struct {
	protofile.UnimplementedUserServiceServer
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	handleError(err)
	return os.Getenv(key)
}
func (*server) UserData(ctx context.Context, req *protofile.UserRequest) (*protofile.UserResponse, error) {
	fmt.Println(req)
	name := req.GetUser().GetName()
	email := req.GetUser().GetEmail()
	phone := req.GetUser().GetPhone()

	newUserItem := user_item{
		Name:  name,
		Email: email,
		Phone: phone,
	}
	dbres := pushUserToDb(ctx, newUserItem)
	result := fmt.Sprintf("%v", dbres)

	userAddResponse := protofile.UserResponse{
		Result: result,
	}
	return &userAddResponse, nil
}
func (*server) ActData(ctx context.Context, req *protofile.ActRequest) (*protofile.ActResponse, error) {
	fmt.Println(req)
	activitytype := req.GetActivity().GetActivitytype()
	duration := req.GetActivity().GetDuration()
	label := req.GetActivity().GetLabel()
	timestamp := req.GetActivity().GetTimestamp()
	email := req.GetActivity().GetEmail()
	newActItem := activity_item{
		Activitytype: activitytype,
		Duration:     duration,
		Label:        label,
		Timestamp:    timestamp,
		Email:        email,
	}
	dbres := pushActToDb(ctx, newActItem)
	result := fmt.Sprintf("%v", dbres)
	actAddResponse := protofile.ActResponse{
		Result: result,
	}
	return &actAddResponse, nil

}

var U_collection *mongo.Collection
var A_collection *mongo.Collection

func (*server) GetActivity(ctx context.Context, req *protofile.GetActivityRequest) (*protofile.GetActivityResponse, error) {
	fmt.Println(req)
	email := req.GetEmail()
	filter := bson.M{
		"email": email,
	}
	var result_data []activity_item
	cursor, err := A_collection.Find(context.TODO(), filter)
	handleError(err)
	cursor.All(context.Background(), &result_data)
	if len(result_data) == 0 {
		getActivityResponse := protofile.GetActivityResponse{
			Status:   false,
			Activity: nil,
		}
		return &getActivityResponse, nil
	} else {
		getActivityResponse := protofile.GetActivityResponse{
			Status: true,
			Activity: &protofile.Activity{
				Activitytype: result_data[0].Activitytype,
				Timestamp:    result_data[0].Timestamp,
				Duration:     result_data[0].Duration,
				Label:        result_data[0].Label,
				Email:        result_data[0].Email,
			},
		}
		return &getActivityResponse, nil
	}
}
func (*server) GetUser(ctx context.Context, req *protofile.GetUserRequest) (*protofile.GetUserResponse, error) {
	fmt.Println(req)
	email := req.GetEmail()
	filter := bson.M{
		"email": email,
	}
	var result_data []user_item
	cursor, err := U_collection.Find(context.TODO(), filter)
	handleError(err)
	cursor.All(context.Background(), &result_data)
	if len(result_data) == 0 {
		getUserResponse := protofile.GetUserResponse{
			Status: false,
			User:   nil,
		}
		return &getUserResponse, nil
	} else {
		getUserResponse := protofile.GetUserResponse{
			Status: true,
			User: &protofile.User{
				Email: result_data[0].Email,
				Name:  result_data[0].Name,
				Phone: result_data[0].Phone,
			},
		}
		return &getUserResponse, nil
	}
}
func (*server) RemoveUser(ctx context.Context, req *protofile.RemoveUserRequest) (*protofile.RemoveUserResponse, error) {
	fmt.Println(req)
	email := req.GetEmail()
	filter := bson.M{
		"email": email,
	}
	A_collection.DeleteMany(context.TODO(), filter)
	u_r, err := U_collection.DeleteOne(context.TODO(), filter)
	handleError(err)
	var result string
	if u_r.DeletedCount == 0 {
		result = "User does not exist"
	} else {
		result = "User deleted successfully"
	}
	removeUserResponse := protofile.RemoveUserResponse{
		Result: result,
	}
	return &removeUserResponse, nil
}
func main() {
	godotenv.Load(".env")
	fmt.Println("GRPC Server Started...")
	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)

	protofile.RegisterUserServiceServer(s, &server{})

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	handleError(err)

	go func() {
		if err := s.Serve(lis); err != nil {
			handleError(err)
		}
	}()

	mongo_uri := goDotEnvVariable("MONGOURL")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongo_uri))
	handleError(err)

	fmt.Println("MongoDB Connected")

	err = client.Connect(context.TODO())
	handleError(err)

	U_collection = client.Database("useractivity").Collection("userdata")
	A_collection = client.Database("useractivity").Collection("useractivitydata")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	fmt.Println("Closing mongodb connection")
	if err := client.Disconnect(context.TODO()); err != nil {
		handleError(err)
	}

	s.Stop()

}
