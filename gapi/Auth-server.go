package gapi

import (
	"context"
	"monprot/config"
	"monprot/pb"
	"monprot/services"

	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
	config         config.Config
	authService    services.AuthService
	userService    services.UserService
	userCollection *mongo.Collection
	ctx            *context.Context
}

func NewGrpcServer(config config.Config, authService services.AuthService,
	userService services.UserService, userCollection *mongo.Collection, ctx *context.Context) (*Server, error) {

	server := &Server{
		config:         config,
		authService:    authService,
		userService:    userService,
		userCollection: userCollection,
		ctx:            ctx,
	}

	return server, nil
}
