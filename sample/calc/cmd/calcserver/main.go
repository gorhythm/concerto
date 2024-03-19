// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package main

import (
	"log/slog"
	"net"
	"os"

	"github.com/apache/thrift/lib/go/thrift"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/oklog/run"
	"google.golang.org/grpc"

	"github.com/gorhythm/concerto/sample/calc"
	pb "github.com/gorhythm/concerto/sample/calc/concerto/proto/gen-go/concerto/sample/calc/v1"
	"github.com/gorhythm/concerto/sample/calc/concerto/service/calculator/endpoint"
	grpcserver "github.com/gorhythm/concerto/sample/calc/concerto/service/calculator/transport/grpc/server"
	thriftserver "github.com/gorhythm/concerto/sample/calc/concerto/service/calculator/transport/thrift/server"
	thriftgen "github.com/gorhythm/concerto/sample/calc/concerto/thrift/gen-go/concerto/sample/calculator/v1"
)

func main() {
	var (
		grpcAddr   = ":8081"
		thriftAddr = ":8082"

		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

		svc       = calc.NewCalculatorService(logger)
		endpoints = endpoint.New(svc)

		grpcServer   = grpcserver.New(endpoints)
		thriftServer = thriftserver.New(endpoints)
	)

	var g run.Group
	{
		grpcListener, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			logger.Error("listen tcp for grpc server", "addr", grpcAddr, "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Info("starting gRPC server", "addr", grpcAddr)
			baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))

			pb.RegisterCalculatorServiceServer(baseServer, grpcServer)
			return baseServer.Serve(grpcListener)
		}, func(error) {
			grpcListener.Close()
		})
	}
	{
		thriftSocket, err := thrift.NewTServerSocket(thriftAddr)
		if err != nil {
			logger.Error("create Thrift socker server", "addr", thriftAddr, "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Info("starting Thrift server", "addr", thriftAddr)
			protocolFactory := thrift.NewTHeaderProtocolFactoryConf(nil)
			transportFactory := thrift.NewTFramedTransportFactoryConf(
				thrift.NewTBufferedTransportFactory(8192), nil,
			)
			return thrift.NewTSimpleServer4(
				thriftgen.NewCalculatorServiceProcessor(thriftServer),
				thriftSocket,
				transportFactory,
				protocolFactory,
			).Serve()
		}, func(error) {
			thriftSocket.Close()
		})
	}

	if err := g.Run(); err != nil {
		logger.Error("exit", "err", err)
	}
}
