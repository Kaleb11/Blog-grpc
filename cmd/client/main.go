package main

import (
	"log"

	"monprot/client"
	"monprot/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "0.0.0.0:8080"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	defer conn.Close()

	//Sign Up
	// if true {
	// 	signUpUserClient := client.NewSignUpUserClient(conn)
	// 	newUser := &pb.SignUpUserInput{
	// 		Name:            "Km",
	// 		Email:           "km@gmail.com",
	// 		Password:        "password123",
	// 		PasswordConfirm: "password123",
	// 	}
	// 	signUpUserClient.SignUpUser(newUser)
	// }

	//Sign In
	// if true {
	// 	signInUserClient := client.NewSignInUserClient(conn)

	// 	credentials := &pb.SignInUserInput{
	// 		Email:    "kalebtilahun29@gmail.com",
	// 		Password: "12345678",
	// 	}
	// 	signInUserClient.SignInUser(credentials)
	// }

	// Get Me
	if true {

		getMeClient := client.NewGetMeClient(conn)
		id := &pb.GetMeRequest{
			Id: "63907c51262031fbe144456b",
		}
		getMeClient.GetMeUser(id)

	}

}
