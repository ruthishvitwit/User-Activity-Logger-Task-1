package main

import (
	"context"
	"fmt"
	"log"
	"main/protofile"
	"net"
	"os"
	"os/signal"

	// "github.com/joho/godotenv"
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

func pushUserToDb(ctx context.Context, item user_item) string {
	email := item.Email
	filter := bson.M{
		"email": email,
	}

	var result_data []user_item
	cursor, err := collection.Find(context.TODO(), filter)
	handleError(err)

	cursor.All(context.Background(), &result_data)

	if len(result_data) != 0 {
		result := "User already exist"
		return result
	}

	collection.InsertOne(ctx, item)
	result := "User created"
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

var collection *mongo.Collection

func main() {
	// godotenv.Load(".env")

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

	// mongo_uri := goDotEnvVariable("MONGODB_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://rutish2722:rutish2722@cluster0.5mkby27.mongodb.net/?retryWrites=true&w=majority"))
	handleError(err)

	fmt.Println("MongoDB Connected")

	err = client.Connect(context.TODO())
	handleError(err)

	collection = client.Database("useractivity").Collection("userdata")
	// collection = client.Database("useractivity").Collection("useractivitydata")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	fmt.Println("Closing mongodb connection")
	if err := client.Disconnect(context.TODO()); err != nil {
		handleError(err)
	}

	s.Stop()

}
