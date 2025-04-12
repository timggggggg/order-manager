package mw

import (
	"context"
	"log"

	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/timofey15g/homework/internal/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func Logging(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, info.FullMethod)
	defer span.Finish()

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("method: %s; metadata: %v", info.FullMethod, md)
		span.LogKV("method", info.FullMethod, "metadata", md)
	}

	rewReq, _ := protojson.Marshal((req).(proto.Message))
	log.Printf("method: %s; request: %s", info.FullMethod, string(rewReq))
	span.LogKV("request", string(rewReq))

	res, err := handler(ctx, req)
	if err != nil {
		log.Printf("method: %s; error: %s", info.FullMethod, err.Error())
		span.SetTag("error", true)
		span.LogKV("error_message", err.Error())

		metrics.IncBadRespByHandler(info.FullMethod)
		return nil, err
	}

	respReq, _ := protojson.Marshal((res).(proto.Message))
	log.Printf("method: %s; response: %s", info.FullMethod, string(respReq))
	span.LogKV("response", string(respReq))

	metrics.IncOkRespByHandler(info.FullMethod)
	return res, nil
}
