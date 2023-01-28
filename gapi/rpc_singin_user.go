package gapi

import (
	"context"
	"monprot/config"
	"monprot/models"
	"monprot/pb"
	"monprot/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) SignInUser(ctx context.Context, req *pb.SignInUserInput) (*pb.GenericSignInResponse, error) {
	user := models.SignInInput{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	newUser, err := server.userService.FindUserByEmail(user.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.InvalidArgument, "Invalid email or password : %s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	if err := utils.VerifyPassword(newUser.Password, user.Password); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid email or password : %s", err.Error())
	}

	config, _ := config.LoadConfig(".")

	// Generate Tokens

	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, newUser.ID, config.AccessTokenPrivateKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	refresh_token, err := utils.CreateToken(config.RefreshTokenExpiresIn, newUser.ID, config.AccessTokenPrivateKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	res := &pb.GenericSignInResponse{
		Status:       "success",
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}
	return res, nil
}
