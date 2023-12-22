package gapi

import (
	"fmt"

	db "github.com/Jingyii800/simplebank/db/sqlc"
	"github.com/Jingyii800/simplebank/pb"
	"github.com/Jingyii800/simplebank/token"
	"github.com/Jingyii800/simplebank/util"
	"github.com/Jingyii800/simplebank/worker"
)

// Server serves gRPC requests for the banking service.
type Server struct {
	pb.UnimplementedSimpleBankServer // necessary to show even the api is unimplemented
	config                           util.Config
	store                            db.Store
	tokenMaker                       token.Maker
	taskDistributor                  worker.TaskDistributor
}

// NewServer creates a new gRPC server.
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
