package gapi

import (
	"context"
	"time"

	db "github.com/Jingyii800/simplebank/db/sqlc"
	"github.com/Jingyii800/simplebank/pb"
	"github.com/Jingyii800/simplebank/util"
	"github.com/Jingyii800/simplebank/val"
	"github.com/Jingyii800/simplebank/worker"
	"github.com/hibiken/asynq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := ValidateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		status.Errorf(codes.Internal, "failed to hash password, %s", err)
	}
	// create user and send task to redis in 1 single DB transaction
	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.GetUsername(),
			HashedPassword: hashedPassword,
			FullName:       req.GetFullName(),
			Email:          req.GetEmail(),
		},
		AfterCreate: func(user db.User) error {
			// Send Email to verify
			taskPayload := &worker.PayloadSendVerifyEmail{
				Username: user.Username,
			}
			// add some options for asynq
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second), // delay 10s to execute
				asynq.Queue(worker.QueueCritical),
			}
			return server.taskDistributor.DistributeTaskSendverifyEmail(ctx, taskPayload, opts...)
		},
	}

	// Create User
	txResult, err := server.store.CreateUserTx(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, "username already exists, %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to create user, %s", err)
	}

	rsp := &pb.CreateUserResponse{
		// need to convert to pb.User
		User: convertUser(txResult.User),
	}

	return rsp, nil
}

func ValidateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := val.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations

}
