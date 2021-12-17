package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	pb "github.com/leslesnoa/go-microservices/comment/commentpb"
	"google.golang.org/grpc"
)

const (
	port = "localhost:9091"
)

// Echoのリクエスト/レスポンスボディ変換プラグイン
func bodyDumpHandler(c echo.Context, reqBody, resBody []byte) {
	fmt.Printf("Request Body: %v\n", string(reqBody))
	fmt.Printf("Response Body: %v\n", string(resBody))
}

// Interceptorの定義
func unaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log.Printf("before call: %s, request: %+v", method, req)
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("after call: %s, response: %+v", method, reply)
	return err
}

func GetAllComment(c echo.Context) error {
	log.Println("starting on getPostFunc.")
	conn, err := grpc.Dial(port, grpc.WithInsecure(), grpc.WithUnaryInterceptor(unaryInterceptor))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewCommentServiceClient(conn)
	ctx := context.Background()
	res, err := client.GetAllComment(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("ResponseValue: %v", res)
	return c.JSON(http.StatusOK, res)
}

func CreateComment(c echo.Context) error {
	log.Println("starting on CreatePostFunc.")
	var p pb.Comment

	log.Printf("request body: %v", p)
	log.Printf("request body: %v", c.Get("text"))

	if err := c.Bind(&p); err != nil {
		fmt.Sprintf("Bind error! %s", err)
	}
	postId, err := strconv.ParseInt(c.Param("post_id"), 10, 64)
	if err != nil {
		fmt.Println("parse error!")
	}
	p.CommentByPostId = postId

	conn, err := grpc.Dial(port, grpc.WithInsecure(), grpc.WithUnaryInterceptor(unaryInterceptor))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewCommentServiceClient(conn)

	// contextの準備
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.CreateComment(ctx, &p)
	// res, err := client.CreatePost(ctx, &pb.Post{Title: "test", Text: "test"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	return c.JSON(http.StatusOK, res)
}

func GetCommentByPostId(c echo.Context) error {
	log.Println("starting on GetCommentByPostIdFunc.")
	conn, err := grpc.Dial(port, grpc.WithInsecure(), grpc.WithUnaryInterceptor(unaryInterceptor))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewCommentServiceClient(conn)
	ctx := context.Background()

	// Get Path Parameter
	// postId := c.Param("post_id")
	postId, err := strconv.ParseInt(c.Param("post_id"), 10, 64)
	if err != nil {
		fmt.Println("parse error!")
	}

	res, err := client.GetCommentByPostId(ctx, &pb.GetCommentByPostIdRequest{CommentByPostId: postId})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("ResponseValue: %v", res)
	return c.JSON(http.StatusOK, res)
}

func main() {
	e := echo.New()
	// ミドルウェアを設定
	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.BodyDump(bodyDumpHandler))
	e.GET("posts/comments", GetAllComment)
	e.POST("posts/:post_id/comments", CreateComment)
	e.GET("posts/:post_id/comments", GetCommentByPostId)
	e.Logger.Fatal(e.Start(":4001"))
}
