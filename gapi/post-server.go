package gapi

import (
	"monprot/pb"
	"monprot/services"

	"go.mongodb.org/mongo-driver/mongo"
)

type PostServer struct {
	pb.UnimplementedPostServiceServer
	postCollection *mongo.Collection
	postService    services.PostService
}

func NewGrpcPostServer(postCollection *mongo.Collection, postService services.PostService) (*PostServer, error) {
	postServer := &PostServer{
		postCollection: postCollection,
		postService:    postService,
	}

	return postServer, nil
}
