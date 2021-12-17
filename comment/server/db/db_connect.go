package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/leslesnoa/go-microservices/comment/commentpb"
	pb "github.com/leslesnoa/go-microservices/comment/commentpb"
)

// type Post struct {
// 	Id   int
// 	Name string
// 	Text string
// }

func Connect() *sql.DB {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	log.Println("start dbConnectFunc.")
	user := "test"
	password := "test"
	host := "localhost"
	port := "3307"
	dbName := "testdb"

	conn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbName + "?charset=utf8mb4"

	db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Println("DB接続時エラー")
		panic(err.Error())
	}
	return db
}

func GetRows(db *sql.DB) []*pb.Comment {
	log.Println("start GetRowsFunc.")
	cmd := "SELECT * FROM comments;"
	rows, err := db.Query(cmd)
	if err != nil {
		log.Println("Query実行時エラー")
		log.Println(err.Error())
	}
	defer rows.Close()
	log.Println("get rows query success.")

	var result []*pb.Comment
	for rows.Next() {
		var p *pb.Comment
		p = new(pb.Comment)
		err := rows.Scan(&p.Id, &p.CommentByPostId, &p.Content)
		if err != nil {
			log.Println("クエリ実行時エラー")
			log.Fatal(err.Error())
		}
		result = append(result, p)
	}
	return result
}

func CreateRow(db *sql.DB, r *pb.Comment) {
	log.Println("Starting CreateRowFunc.")
	// db.Prepare()
	stmtInsert, err := db.Prepare("INSERT INTO comments(comment_by_post_id, content) VALUES(?, ?);")
	if err != nil {
		panic(err.Error())
	}
	defer stmtInsert.Close()

	post := pb.Comment{CommentByPostId: r.CommentByPostId, Content: r.Content}
	result, err := stmtInsert.Exec(post.CommentByPostId, post.Content)
	if err != nil {
		panic(err.Error())
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(lastInsertID)

	// return &pb.User{Id: int32(lastInsertID), Name: r.Name, Email: r.Email}, nil
}

func GetRowByPostId(db *sql.DB, r *pb.GetCommentByPostIdRequest) *pb.Comment {
	log.Println("start GetRowByPostIdFunc.")
	var p *commentpb.Comment
	p = new(pb.Comment)
	err := db.QueryRow("SELECT * FROM comments WHERE comment_by_post_id=?;", r.CommentByPostId).Scan(&p.Id, &p.CommentByPostId, &p.Content)
	if err != nil {
		panic(err)
	}
	log.Println("get rows query success.")

	// err := row.Scan(r.CommentByPostId)
	// if err != nil {
	// 	log.Println("クエリ実行時エラー")
	// 	log.Fatal(err.Error())
	// }
	return p
}
