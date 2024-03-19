// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"

	"github.com/apache/thrift/lib/go/thrift"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gorhythm/concerto/transport"

	"github.com/gorhythm/concerto/sample/calc"
	grpcclient "github.com/gorhythm/concerto/sample/calc/concerto/service/calculator/transport/grpc/client"
	thriftclient "github.com/gorhythm/concerto/sample/calc/concerto/service/calculator/transport/thrift/client"
	thriftgen "github.com/gorhythm/concerto/sample/calc/concerto/thrift/gen-go/concerto/sample/calculator/v1"
)

func main() {
	var (
		grpcAddr   = ":8081"
		thriftAddr = ":8082"

		logger  = slog.New(slog.NewTextHandler(os.Stdout, nil))
		clients = map[transport.Transport]calc.CalculatorService{}
	)

	{
		conn, err := grpc.DialContext(
			context.Background(),
			grpcAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			logger.Error("connect to gRPC server", "addr", grpcAddr, "err", err)
			os.Exit(1)
		}
		defer conn.Close()

		clients[transport.TransportGRPC] = grpcclient.New(conn)
	}
	{
		var (
			tTransport       thrift.TTransport = thrift.NewTSocketConf(thriftAddr, nil)
			protocolFactory                    = thrift.NewTHeaderProtocolFactoryConf(nil)
			transportFactory                   = thrift.NewTFramedTransportFactoryConf(
				thrift.NewTBufferedTransportFactory(8192), nil,
			)
		)

		tTransport, err := transportFactory.GetTransport(tTransport)
		if err != nil {
			logger.Error("create Thrift transport", "err", err)
			os.Exit(1)
		}

		if err := tTransport.Open(); err != nil {
			logger.Error("connect to Thrift server", "err", err)
			os.Exit(1)
		}
		defer tTransport.Close()

		clients[transport.TransportThrift] = thriftclient.New(
			thriftgen.NewCalculatorServiceClientFactory(tTransport, protocolFactory),
		)
	}

	for trans, client := range clients {
		var (
			op = [...]calc.Op{
				calc.OpAdd, calc.OpSubtract, calc.OpMultiply, calc.OpDivide,
			}[rand.Intn(4)]
			num1 = rand.Int63n(10000)
			num2 = rand.Int63n(10000)
		)

		result, err := client.Calculate(context.Background(), op, num1, num2)
		if err != nil {
			logger.Error(
				"call CalculatorService/Calculate",
				"transport", trans, "err", err,
			)
			continue
		}

		logger.Info(
			fmt.Sprintf("calculate(%s, %d, %d) = %d", op, num1, num2, result),
			"transport", trans,
		)
	}
}
