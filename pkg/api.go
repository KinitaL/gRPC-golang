package api

import (
	"context"
	"database/sql"
	"hezzlgrpc/pb"
	"hezzlgrpc/pkg/model"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type UserServer struct{}

func (us *UserServer) Get(ctx context.Context, req *pb.GetRequest) (res *pb.GetResponse, err error) {

	/*Load env file*/

	envLoad()

	var users []string
	/*
		DB
	*/
	db := connectToDB()
	defer db.Close()
	stmt, err := db.Prepare("SELECT * FROM users;")
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	checkError(err)
	defer rows.Close()

	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		checkError(err)

		userstring := "name:" + user.Name + ", email:" + user.Email
		users = append(users, userstring)
	}
	err = rows.Err()
	checkError(err)

	/*create redis client*/

	redisHelper := envManage("REDISHOST") + ":" + envManage("REDISPORT")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHelper,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	/*add users to redis*/

	for i := 0; i < len(users); i++ {

		rdb.SetNX(ctx, "User"+strconv.Itoa(i), users[i], 60*time.Second)
	}

	/*take users from redis*/
	var usersR []string
	for i := 0; i < len(users); i++ {
		userR, err := rdb.Get(ctx, "User"+strconv.Itoa(i)).Result()
		checkError(err)
		usersR = append(usersR, userR)
	}

	return &pb.GetResponse{User: usersR}, nil

}

func (us *UserServer) Post(ctx context.Context, req *pb.PostRequest) (res *pb.PostResponse, err error) {

	/*
		handle request

		example of result of req.GetUser(): "name:Kujo,email:jotaro@mail.com"
	*/

	user := req.GetUser()
	userData := strings.Split(user, ",")
	userName := strings.Split(userData[0], ":")[1]
	userEmail := strings.Split(userData[1], ":")[1]

	/*
		DB
	*/
	db := connectToDB()
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO users (name, email) VALUES ($1, $2);")
	checkError(err)
	defer stmt.Close()
	_, err = stmt.Exec(userName, userEmail)
	checkError(err)

	/*return message*/
	return &pb.PostResponse{Message: "User " + userEmail + " was added."}, nil

}

func (us *UserServer) Delete(ctx context.Context, req *pb.DeleteRequest) (res *pb.DeleteResponse, err error) {
	/*
		handle request (Post and Delete)

		example of result of req.GetUser(): "name:Kujo,email:jotaro@mail.com"
	*/
	user := req.GetUser()
	userData := strings.Split(user, ",")
	userEmail := strings.Split(userData[1], ":")[1]

	/*
		DB
	*/
	db := connectToDB()
	defer db.Close()
	stmt, err := db.Prepare("DELETE FROM users WHERE email=$1;")
	checkError(err)
	defer stmt.Close()
	_, err = stmt.Exec(userEmail)
	checkError(err)

	/*return message*/
	return &pb.DeleteResponse{Message: "User " + userEmail + " was deleted."}, nil
}

func connectToDB() *sql.DB {
	host := envManage("SQLHOST")
	port := envManage("SQLPORT")
	user := envManage("SQLUSER")
	pass := envManage("SQLPASS")
	dbname := envManage("SQLDB")
	helper := "host=" + host + " port=" + port + " user=" + user + " password=" + pass + " dbname=" + dbname + " sslmode=disable"
	db, err := sql.Open("postgres", helper)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func envLoad() {
	if err := godotenv.Load("./mainenv.env"); err != nil {
		log.Fatal("No .env file found")
	}
}
func envManage(key string) string {

	value, exists := os.LookupEnv(key)
	if exists {
		return value
	} else {
		log.Fatal("No conn")
		return ""
	}
}
