package gapi

import (
	"context"
	"monprot/models"
	"monprot/pb"
	"monprot/utils"

	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) ResetPassword(ctx context.Context, req *pb.ResetPasswordInput) (*pb.GenericResponse, error) {
	// code := server.reqs.GetResetCode()
	user := models.ResetPasswordInput{
		Password:        req.GetPassword(),
		PasswordConfirm: req.GetPasswordConfirm(),
	}

	if user.Password != user.PasswordConfirm {
		return nil, status.Errorf(codes.InvalidArgument, "Passwords do not match: %s")
	}

	hashedPassword, _ := utils.HashPassword(user.Password)

	passwordResetToken := utils.Encode(req.GetResetCode())

	// Update User in Database
	query := bson.D{{Key: "passwordResetToken", Value: passwordResetToken}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: hashedPassword}}}, {Key: "$unset", Value: bson.D{{Key: "passwordResetToken", Value: ""}, {Key: "passwordResetAt", Value: ""}}}}
	result, err := server.userCollection.UpdateOne(*server.ctx, query, update)

	if result.MatchedCount == 0 {
		return nil, status.Errorf(codes.Canceled, "%s", "Token is invalid or has expired")
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	// ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	// ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	// ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)
	res := &pb.GenericResponse{
		Status:  "success",
		Message: "Password data updated successfully",
	}
	return res, nil
}
