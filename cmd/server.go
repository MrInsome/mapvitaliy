package main

import (
	"apitraning/internal"
	"apitraning/pkg"
	"apitraning/pkg/api"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
)

type serverStruct struct {
	repository *pkg.Repository
	api.UnsubscribeServiceServer
}

func (s *serverStruct) Unsubscribe(ctx context.Context, account *internal.Account) (*emptypb.Empty, error) {
	db := s.repository.ReturnDB()
	r := &pkg.Repository{}
	err := r.UnsubscribeAccount(db, account.AccountID)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func OpenGRPC() {
	lis, err := net.Listen("tcp", ":50051")
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
