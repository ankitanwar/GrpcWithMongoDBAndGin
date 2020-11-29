package client

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/ankitanwar/GrpcWithMongoDBAndGin/blogpb"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var (
	router = gin.Default()
	//C service Client
	C  blogpb.BlogServiceClient
	cc *grpc.ClientConn
)

//StartClient : To start the client service
func StartClient() {
	urlMapping()
	connectServer()
	go func() {
		router.Run(":8081")
	}()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	fmt.Println("Closing the Connection with server")
	cc.Close()

}
func connectServer() {
	opts := grpc.WithInsecure()
	var err error
	cc, err = grpc.Dial("localhost:4040", opts)
	if err != nil {
		fmt.Println("Error while connection to the server", err.Error())
		panic(err)
	}
	C = blogpb.NewBlogServiceClient(cc)
	fmt.Println("Connection to Server is successfull")
}
