package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/ankitanwar/GrpcWithMongoDBAndGin/blogpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Title    string             `bson:"title"`
	Content  string             `bson:"content"`
}

var (
	client     *mongo.Client
	collection *mongo.Collection
)

type server struct {
}

func (s *server) Create(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	blog := req.GetBlog()
	data := blogItem{
		AuthorID: blog.GetAuthorID(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}
	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			fmt.Sprintf("Error while Inserting the blog into the database %v", err),
		)
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Error(
			codes.Internal,
			fmt.Sprintf("Error while Inserting the blog into the database %v", err),
		)
	}
	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			ID:       oid.Hex(),
			AuthorID: blog.GetAuthorID(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
		},
	}, nil

}
func (s *server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	blogID := req.GetBlogID()
	oid, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot convert to id %v", err),
		)
	}
	filter := bson.M{"_id": oid}
	res := &blogItem{}
	findErr := collection.FindOne(context.Background(), filter).Decode(res)
	if findErr != nil {
		return nil, status.Error(
			codes.NotFound,
			fmt.Sprintf("Cannot find the blog %v", err),
		)
	}
	return &blogpb.ReadBlogResponse{
		Blog: &blogpb.Blog{
			ID:       res.ID.Hex(),
			AuthorID: res.AuthorID,
			Title:    res.Title,
			Content:  res.Content,
		},
	}, nil
}

func main() {
	//setting up mongo Server
	mongoSetup()

	lis, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatalln("Unable to listen")

		panic(err)
	}
	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalln("Unable to server")
			panic(err)
		}
	}()

	//Waiting for the stop signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("closing the listner")
	lis.Close()
	fmt.Println("Closing MongoDB Sever")
	client.Disconnect(context.TODO())

}

func mongoSetup() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalln("Unable to connect to mongoDB Sever")
		panic(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalln("Error while ping")
		panic(err)
	}
	collection = client.Database("grpc").Collection("blog")
}
