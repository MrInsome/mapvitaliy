package pkg

import (
	"apitraning/pkg/api"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"log"
	"net"
)

type serverStruct struct {
	repository *Repository
	db         *gorm.DB
	api.UnsubscribeServiceServer
}

func (s *serverStruct) Unsubscribe(ctx context.Context, account *api.AccountRequest) (*emptypb.Empty, error) {
	db := s.db
	r := &Repository{}
	err := r.UnsubscribeAccount(db, int(account.AccountId))
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func OpenGRPC() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("%v", err)
	}
	s := grpc.NewServer()
	srv := &serverStruct{}
	api.RegisterUnsubscribeServiceServer(s, srv)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("%v", err)
	}
}
