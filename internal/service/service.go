package service

import (
	"context"
	"errors"
	"log"
	pb "rkv/gen"
	"rkv/internal/store"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	pb.UnimplementedShardServiceServer
	store *store.Store
	addr  string
}

func New(store *store.Store, addr string) *Service {
	return &Service{store: store, addr: addr}
}

func (s *Service) Put(
	ctx context.Context,
	req *pb.PutRequest,
) (*pb.PutResponse, error) {
	key, data := req.Key, req.Data
	if key == "" || len(data) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty key or data")
	}

	s.store.Put(key, []byte(data))
	return &pb.PutResponse{}, nil
}

func (s *Service) Get(
	ctx context.Context,
	req *pb.GetRequest,
) (*pb.GetResponse, error) {
	key := req.Key
	if key == "" {
		return nil, status.Error(codes.InvalidArgument, "empty key")
	}

	data, err := s.store.Get(key)
	if err != nil {
		if errors.Is(err, store.ErrKeyNotFound) {
			return nil, status.Error(codes.NotFound, "data not found for given key")
		}
		// this should never hit
		log.Printf("unexpcted error while getting: %v\n", err)
		return nil, status.Error(codes.Internal, "internal service error")
	}

	return &pb.GetResponse{Data: data}, nil
}

func (s *Service) Exists(
	ctx context.Context,
	req *pb.ExistsRequest,
) (*pb.ExistsResponse, error) {
	key := req.Key
	if key == "" {
		return nil, status.Error(codes.InvalidArgument, "empty key")
	}

	exists := s.store.Exists(key)
	return &pb.ExistsResponse{Exists: exists}, nil
}

func (s *Service) Delete(
	ctx context.Context,
	req *pb.DeleteRequest,
) (*pb.DeleteResponse, error) {
	key := req.Key
	if key == "" {
		return nil, status.Error(codes.InvalidArgument, "empty key")
	}

	if err := s.store.Delete(key); err != nil {
		if errors.Is(err, store.ErrKeyNotFound) {
			return nil, status.Error(codes.NotFound, "data not found for given key")
		}
		// this should never hit
		log.Printf("unexpcted error while deleting: %v\n", err)
		return nil, status.Error(codes.Internal, "internal service error")
	}

	return &pb.DeleteResponse{}, nil
}

func (s *Service) Info(
	ctx context.Context,
	req *pb.InfoRequest,
) (*pb.InfoResponse, error) {
	return &pb.InfoResponse{Name: s.addr}, nil
}
