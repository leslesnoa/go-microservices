package main

import (
	"context"
	"log"
	"net"

	_ "github.com/go-sql-driver/mysql"
	pb "github.com/leslesnoa/go-microservices/comment/commentpb"
	"github.com/leslesnoa/go-microservices/comment/server/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// gRPCサーバ
type server struct {
}

const (
	port = ":9091"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterCommentServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %c", err)
	}
	log.Println("starting on gRPC server on " + port)
}

func (s *server) CreateComment(ctx context.Context, r *pb.Comment) (*pb.Empty, error) {
	log.Printf("Recieved CreatePostRequest : %s", r)
	log.Println(r.GetCommentByPostId())
	log.Println(r.GetContent())
	// DB接続
	conn := db.Connect()
	db.CreateRow(conn, r)
	return &pb.Empty{}, nil
}

func (s *server) GetAllComment(ctx context.Context, r *pb.Empty) (*pb.Comments, error) {
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
	var result []*pb.Comment
	for _, p := range rows {
		log.Println(p.CommentByPostId, p.Content)
		result = append(result, p)
	}
	defer conn.Close()
	return &pb.Comments{Comments: rows}, nil

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

func (s *server) GetCommentByPostId(ctx context.Context, r *pb.GetCommentByPostIdRequest) (*pb.Comments, error) {
	log.Printf("Recieved GetChatRequest : %s", r)

	// DB接続
	conn := db.Connect()
	defer conn.Close()

	// 接続確認
	err := conn.Ping()
	if err != nil {
		log.Println("connection failed.")
		panic(err)
	} else {
		log.Println("connection success.")
	}

	// 行データ取得
	row := db.GetRowByPostId(conn, r)
	log.Println(row)
	var result []*pb.Comment
	result = append(result, row)

	return &pb.Comments{Comments: result}, nil
}
