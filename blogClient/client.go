package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ankitanwar/GrpcWithMongoDBAndGin/blogpb"
	"google.golang.org/grpc"
)

func main() {
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:8080", opts)
	if err != nil {
		fmt.Println("Error while connection to the server", err.Error())
		panic(err)
	}
	defer cc.Close()
	c := blogpb.NewBlogServiceClient(cc)
	blog := &blogpb.Blog{
		AuthorID: "hello",
		Title:    "my first blog",
		Content:  "hello world",
	}
	response, err := c.Create(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error %v", err)
		panic(err)
	}

	fmt.Println("The response from the server is ", response)

	fmt.Println("Reading the Blog")

	req := blogpb.ReadBlogRequest{
		BlogID: "5fc1fc3efd9dd54eeab86a3a",
	}

	res, readErr := c.ReadBlog(context.Background(), &req)
	if readErr != nil {
		log.Fatalln("Some error has been occured while reading the blog", readErr)
		panic(err)
	}
	fmt.Println("The value from read request is ", res)

}
