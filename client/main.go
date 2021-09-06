package main

import (
	"context"
	"fmt"
	"hezzlgrpc/pb"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	envLoad()
	serverPort := envManage("SERVERPORT")
	conn, err := grpc.Dial(":"+serverPort, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewUserServiceClient(conn)

	fmt.Println("You are connected to server.")
	for {
		fmt.Println("Write method you want to use (get, post or delete).")
		var method string
		fmt.Scanf("%s\n", &method)
		fmt.Println("You have choosen method: " + method)

		switch method {
		case "post":
			fmt.Println("Write user data without spaces. Pattern: name:Ivan,email:ivan@mail.com")
			var user string
			fmt.Scanf("%s\n", &user)
			res, err := client.Post(context.Background(), &pb.PostRequest{User: user})
			if err != nil {
				log.Fatal(err)
			}

			log.Println(res.GetMessage())

		case "get":
			res, err := client.Get(context.Background(), &pb.GetRequest{})
			if err != nil {
				log.Fatal(err)
			}
			for _, user := range res.GetUser() {
				fmt.Println(user)
			}
		case "delete":
			fmt.Println("Write user data without spaces. Pattern: name:Ivan,email:ivan@mail.com")
			var user string
			fmt.Scanf("%s\n", &user)
			res, err := client.Delete(context.Background(), &pb.DeleteRequest{User: user})
			if err != nil {
				log.Fatal(err)
			}

			log.Println(res.GetMessage())
		default:
			log.Println("wrong method was choosen")
		}
	}

}

func envLoad() {
	if err := godotenv.Load(".././mainenv.env"); err != nil {
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
