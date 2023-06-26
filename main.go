package main

import (
	"match_evaluator/evaluator"
	"match_evaluator/handler"
	"match_evaluator/internal/db"
	"match_evaluator/internal/db/redis"
	"match_evaluator/internal/manager"
	pb "match_evaluator/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	evaluator.DefaultManager = manager.NewManager()
	// Create service
	srv := service.New(
		service.Name("match_evaluator"),
		service.Version("latest"),
		service.BeforeStart(func() error {
			svr, err := redis.New(
				db.WithAddress("127.0.0.1:6379"),
				db.WithPoolMaxActive(5),
				db.WithPoolMaxIdle(100),
				db.WithPoolIdleTimeout(300))
			if err != nil {
				return err
			}
			db.Default = svr
			return nil
		}),
		service.AfterStart(evaluator.DefaultManager.Start),
		service.BeforeStop(evaluator.DefaultManager.Stop),
	)

	// Register handler
	pb.RegisterMatchEvaluatorHandler(srv.Server(), new(handler.Match_evaluator))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
