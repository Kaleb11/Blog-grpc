package gapi

import (
	"context"
	"fmt"
	"log"
	"monprot/config"
	"monprot/models"
	"monprot/pb"
	"monprot/utils"
	"strings"
	"time"

	"github.com/thanhpk/randstr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordInput) (*pb.GenericResponse, error) {

	message := "You will receive a reset email if user with that email exist"

	user := models.ForgotPasswordInput{
		Email: req.GetEmail(),
	}

	newUser, err := server.userService.FindUserByEmail(user.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("Here is")
			return nil, status.Errorf(codes.InvalidArgument, "Invalid email or password : %s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	if !newUser.Verified {
		return nil, status.Errorf(codes.Unauthenticated, "Account not verified : %s", err.Error())
	}

	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load config", err)
	}

	// Generate Verification Code
	resetToken := randstr.String(20)

	passwordResetToken := utils.Encode(resetToken)

	// Update User in Database
	query := bson.D{{Key: "email", Value: strings.ToLower(user.Email)}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "passwordResetToken", Value: passwordResetToken}, {Key: "passwordResetAt", Value: time.Now().Add(time.Minute * 15)}}}}
	result, err := server.userCollection.UpdateOne(*server.ctx, query, update)

	if result.MatchedCount == 0 {
		return nil, status.Errorf(codes.Internal, "There was an error sending email : %s", err.Error())
	}

	if err != nil {
		return nil, status.Errorf(codes.Canceled, "%s", err.Error())
	}
	var firstName = newUser.Name

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// ? Send Email
	emailData := utils.EmailData{
		URL:       config.Origin + "/resetpassword/" + resetToken,
		FirstName: firstName,
		Subject:   "Your password reset token (valid for 10min)",
	}

	err = utils.SendMail(newUser, &emailData, "resetPassword.html")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "There was an error sending email : %s", err.Error())
	}
	res := &pb.GenericResponse{
		Status:  "success",
		Message: message,
	}
	return res, nil
}
