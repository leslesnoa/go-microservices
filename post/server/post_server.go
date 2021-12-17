package main

import (
	"context"
	"log"
	"net"

	_ "github.com/go-sql-driver/mysql"
	pb "github.com/leslesnoa/go-microservices/post/postpb"
	"github.com/leslesnoa/go-microservices/post/server/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// gRPCサーバ
type server struct {
}

const (
	port = ":9090"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPostServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %c", err)
	}
	log.Println("starting on gRPC server on " + port)
}

func (s *server) CreatePost(ctx context.Context, r *pb.Post) (*pb.Empty, error) {
	log.Printf("Recieved CreatePostRequest : %s", r)
	log.Println(r.GetTitle())
	log.Println(r.GetText())
	// DB接続
	conn := db.Connect()
	db.CreateRow(conn, r)
	return &pb.Empty{}, nil
}

func (s *server) GetAllPosts(ctx context.Context, r *pb.Empty) (*pb.Posts, error) {
	log.Printf("Recieved GetChatRequest : %s", r)

	// DB接続
	conn := db.Connect()

	// 接続確認
	err := conn.Ping()
	if err != nil {
		log.Println("connection failed.")
		panic(err)
	} else {
		log.Println("connection success.")
	}

	// 行データ取得
	rows := db.GetRows(conn)
	log.Println(rows)
	var result []*pb.Post
	for _, p := range rows {
		log.Println(p.Title, p.Text)
		result = append(result, p)
	}
	defer conn.Close()
	return &pb.Posts{Posts: rows}, nil

	// return &pb.Posts{Posts: []*pb.Post{
	// 	{
	// 		Id:    1,
	// 		Title: "testTitle",
	// 		Text:  "testText",
	// 	},
	// 	{
	// 		Id:    2,
	// 		Title: "testTitle",
	// 		Text:  "testText",
	// 	},
	// }}, nil
}
