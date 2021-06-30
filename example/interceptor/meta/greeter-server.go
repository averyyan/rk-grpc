// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"context"
	"fmt"
	"github.com/rookie-ninja/rk-entry/entry"
	"github.com/rookie-ninja/rk-grpc/example/interceptor/proto/gen"
	"github.com/rookie-ninja/rk-grpc/interceptor/context"
	"github.com/rookie-ninja/rk-grpc/interceptor/log/zap"
	"github.com/rookie-ninja/rk-grpc/interceptor/meta"
	"google.golang.org/grpc"
	"log"
	"net"
)

// In this example, we will create a simple gRpc server and enable RK style logging interceptor and extension interceptor.
// Then, we will try to send requests to server and monitor what kinds of logging we would get.
func main() {
	// ******************************************************
	// ********** Override App name and version *************
	// ******************************************************
	//
	// rkentry.GlobalAppCtx.GetAppInfoEntry().AppName = "demo-app"
	// rkentry.GlobalAppCtx.GetAppInfoEntry().Version = "demo-version"

	// ********************************************
	// ********** Enable interceptors *************
	// ********************************************
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			rkgrpclog.UnaryServerInterceptor(),
			// Add extension interceptor
			rkgrpcmeta.UnaryServerInterceptor(
			// Entry name and entry type will be used for distinguishing interceptors. Recommended.
			// rkgrpcmeta.WithEntryNameAndType("greeter", "grpc"),
			//
			// We will replace X-<Prefix>-XXX with prefix user provided.
			// rkgrpcmeta.WithPrefix("Dog"),
			),
		),
	}

	// 1: Create grpc server
	server := startGreeterServer(opts...)
	defer server.GracefulStop()

	// 2: Wait for ctrl-C to shutdown server
	rkentry.GlobalAppCtx.WaitForShutdownSig()
}

// Implementation of GreeterServer.
type GreeterServer struct{}

// Handle SayHello method.
func (server *GreeterServer) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloResponse, error) {
	// ******************************************
	// ********** rpc-scoped logger *************
	// ******************************************
	//
	// RequestId will be printed if enabled by bellow codes.
	// 1: Enable rkgrpcmeta.UnaryServerInterceptor() in server side.
	// 2: rkgrpcctx.AddHeaderToClient(ctx, rkgrpcctx.RequestIdKey, rkcommon.GenerateRequestId())
	//
	rkgrpcctx.GetLogger(ctx).Info("Received request from client.")

	// Append request id with X-Request-Id to outgoing headers.
	// Important!
	// This is append operation, not set operation. As a result, client would receive both original request id generated by
	// extension interceptor and new one you specified.
	//
	// However, log interceptor would pick the latest one to attach into zap and event.
	//
	// rkgrpcctx.AddHeaderToClient(ctx, rkgrpcctx.RequestIdKey, "this-is-my-request-id-overridden")

	return &proto.HelloResponse{
		Message: fmt.Sprintf("Hello %s!", request.GetName()),
	}, nil
}

// Create and start server.
func startGreeterServer(opts ...grpc.ServerOption) *grpc.Server {
	// 1: Create listener with port 8080
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// 2: Create grpc server with grpc.ServerOption
	server := grpc.NewServer(opts...)

	// 3: Register server to proto
	proto.RegisterGreeterServer(server, &GreeterServer{})

	// 4: Start server
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	return server
}