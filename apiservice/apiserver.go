package apiservice

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/iotexproject/iotex-analyser-api/api"
	"github.com/iotexproject/iotex-analyser-api/config"
	graphqlruntime "github.com/ysugimoto/grpc-graphql-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func registerAPIService(ctx context.Context, grpcServer *grpc.Server) {
	api.RegisterAccountServiceServer(grpcServer, &AccountService{})
	api.RegisterStakingServiceServer(grpcServer, &StakingService{})
	api.RegisterActionsServiceServer(grpcServer, &ActionsService{})
}

func registerProxyAPIService(ctx context.Context, mux *runtime.ServeMux) error {
	if err := api.RegisterAccountServiceHandlerServer(ctx, mux, &AccountService{}); err != nil {
		return err
	}
	if err := api.RegisterStakingServiceHandlerServer(ctx, mux, &StakingService{}); err != nil {
		return err
	}
	if err := api.RegisterActionsServiceHandlerServer(ctx, mux, &ActionsService{}); err != nil {
		return err
	}
	return nil
}

func registerGraphQLAPIService(ctx context.Context, mux *graphqlruntime.ServeMux) error {
	addr := fmt.Sprintf("127.0.0.1:%d", config.Default.Server.GrpcAPIPort)
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	if err := api.RegisterActionsServiceGraphqlHandler(mux, conn); err != nil {
		return err
	}
	return nil
}

func StartGRPCService(ctx context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Default.Server.GrpcAPIPort))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	registerAPIService(ctx, grpcServer)
	reflection.Register(grpcServer)
	return grpcServer.Serve(lis)
}

func StartGRPCProxyService() error {
	gwmux := runtime.NewServeMux()
	ctx := context.Background()
	if err := registerProxyAPIService(ctx, gwmux); err != nil {
		return err
	}

	graphqlMux := graphqlruntime.NewServeMux()
	if err := registerGraphQLAPIService(ctx, graphqlMux); err != nil {
		return err
	}

	http.Handle("/graphql", graphqlMux)
	http.Handle("/", gwmux)

	port := fmt.Sprintf(":%d", config.Default.Server.HTTPAPIPort)
	return http.ListenAndServe(port, nil)
}
