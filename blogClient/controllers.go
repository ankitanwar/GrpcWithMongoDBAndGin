package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ankitanwar/GrpcWithMongoDBAndGin/blogpb"
	"github.com/gin-gonic/gin"
)

//ReadBlog : To  read the blog of the given  iD
func ReadBlog(c *gin.Context) {

	fmt.Println("Reading the Blog")
	id := c.Param("blogID")

	req := blogpb.ReadBlogRequest{
		BlogID: id,
	}

	res, readErr := C.ReadBlog(context.Background(), &req)
	if readErr != nil {
		c.JSON(http.StatusInternalServerError, readErr)
		return
	}
	c.JSON(http.StatusOK, res)
}

//CreateBlog : To create the new blog
func CreateBlog(c *gin.Context) {
	blog := &blogpb.Blog{}
	err := c.ShouldBindJSON(blog)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	response, err := C.Create(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, response)
}
func hello(c *gin.Context) {
	c.JSON(http.StatusOK, "Hello world")
}
