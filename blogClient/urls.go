package client

func urlMapping() {
	router.GET("/read/blog/:blogID", ReadBlog)
	router.POST("/blog", CreateBlog)
	router.GET("/hello", hello)
}
