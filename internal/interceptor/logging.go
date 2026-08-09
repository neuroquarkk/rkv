package interceptor

import (
	"context"
	"log"
	"path"
	pb "rkv/gen"
	"time"

	"google.golang.org/grpc"
)

func LoggingInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start)

	result := "ok"
	if err != nil {
		result = "error"
	}

	method := path.Base(info.FullMethod)

	log.Printf("method=%s | key=%s | duration=%s | result=%s",
		method, extractKey(req), duration, result)

	return resp, err
}

func extractKey(req any) string {
	switch r := req.(type) {
	case *pb.GetRequest:
	case *pb.PutRequest:
	case *pb.DeleteRequest:
		return r.Key
	default:
		return "-"
	}
	return "-"
}
