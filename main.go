package main

import (
	"hezzlgrpc/pb"
	api "hezzlgrpc/pkg"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"

	"google.golang.org/grpc"
)

func main() {
	server := grpc.NewServer()
	service := &api.UserServer{}

	pb.RegisterUserServiceServer(server, service)

	if err := godotenv.Load("./mainenv.env"); err != nil {
		log.Print("No .env file found")
	}
	serverPort, exists := os.LookupEnv("SERVERPORT")
	if exists {
		listener, err := net.Listen("tcp", ":"+serverPort)

		if err != nil {
			log.Fatal(err)
		}

		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("No conn")
	}

}
