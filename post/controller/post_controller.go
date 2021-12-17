package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	pb "github.com/leslesnoa/go-microservices/post/postpb"
	"google.golang.org/grpc"
)

const (
	port = "localhost:9090"
)

type ChatData struct {
	Chat_Name string `json:"name"`
	Chat_Text string `json:"text"`
}
type ChatList struct {
	Chat_List []*ChatData `json:"chat_list"`
}
type User struct {
}

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

func GetPostAll(c echo.Context) error {
	log.Println("starting on getPostFunc.")
	conn, err := grpc.Dial(port, grpc.WithInsecure(), grpc.WithUnaryInterceptor(unaryInterceptor))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewPostServiceClient(conn)
	ctx := context.Background()
	res, err := client.GetAllPosts(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("ResponseValue: %v", res)
	return c.JSON(http.StatusOK, res)
}

func CreatePost(c echo.Context) error {
	log.Println("starting on CreatePostFunc.")
	var p pb.Post

	log.Printf("request body: %v", p)
	log.Printf("request body: %v", c.Get("text"))

	if err := c.Bind(&p); err != nil {
		fmt.Sprintf("Bind error! %s", err)
	}

	conn, err := grpc.Dial(port, grpc.WithInsecure(), grpc.WithUnaryInterceptor(unaryInterceptor))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewPostServiceClient(conn)

	// contextの準備
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.CreatePost(ctx, &p)
	// res, err := client.CreatePost(ctx, &pb.Post{Title: "test", Text: "test"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	return c.JSON(http.StatusOK, res)
}

func main() {
	e := echo.New()
	// ミドルウェアを設定
	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.BodyDump(bodyDumpHandler))
	e.GET("/posts", GetPostAll)
	e.POST("/posts", CreatePost)
	e.Logger.Fatal(e.Start(":4000"))
}
