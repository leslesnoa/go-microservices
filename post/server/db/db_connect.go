package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	pb "github.com/leslesnoa/go-microservices/post/postpb"
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
	port := "3306"
	dbName := "testdb"

	conn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbName + "?charset=utf8mb4"

	db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Println("DB接続時エラー")
		panic(err.Error())
	}
	return db
}

func GetRows(db *sql.DB) []*pb.Post {
	log.Println("start GetRowsFunc.")
	cmd := "SELECT * FROM posts"
	rows, err := db.Query(cmd)
	if err != nil {
		log.Println("Query実行時エラー")
		log.Println(err.Error())
	}
	defer rows.Close()
	log.Println("get rows query success.")

	var result []*pb.Post
	for rows.Next() {
		var p *pb.Post
		p = new(pb.Post)
		err := rows.Scan(&p.Id, &p.Title, &p.Text)
		if err != nil {
			log.Println("クエリ実行時エラー")
			log.Fatal(err.Error())
		}
		result = append(result, p)
	}
	return result
}

func CreateRow(db *sql.DB, r *pb.Post) {
	log.Println("Starting CreateRowFunc.")
	// db.Prepare()
	stmtInsert, err := db.Prepare("INSERT INTO posts(title, text) VALUES(?,?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmtInsert.Close()

	post := pb.Post{Title: r.Title, Text: r.Text}
	result, err := stmtInsert.Exec(post.Title, post.Text)
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
