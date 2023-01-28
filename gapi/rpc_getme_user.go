package gapi

import (
	"context"

	"monprot/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (userServer *UserServer) GetMe(ctx context.Context, req *pb.GetMeRequest) (*pb.GenericUserResponse, error) {
	id := req.GetId()
	user, err := userServer.userService.FindUserById(id)

	if err != nil {
		return nil, status.Errorf(codes.Unimplemented, err.Error())
	}

	res := &pb.GenericUserResponse{
		User: &pb.User{
			Id:        user.ID.Hex(),
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
	}
	return res, nil
}
