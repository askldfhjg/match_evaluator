package main

import (
	"match_evaluator/handler"
	pb "match_evaluator/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("match_evaluator"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterMatchEvaluatorHandler(srv.Server(), new(handler.Match_evaluator))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
