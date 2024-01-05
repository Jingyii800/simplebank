package gapi

import (
	"context"
	"errors"
	"time"

	db "github.com/Jingyii800/simplebank/db/sqlc"
	"github.com/Jingyii800/simplebank/pb"
	"github.com/Jingyii800/simplebank/util"
	"github.com/Jingyii800/simplebank/val"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// add authorization to protect api
	authpayload, err := server.authorizeUser(ctx, []string{util.BankRole, util.DepositorRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}
	// check violations for input fields
	violations := ValidateUpdateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	// if the authpayload doesn't match the username in req field, access deny
	if authpayload.Role != util.BankRole && authpayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's info")
	}

	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: pgtype.Text{
			String: "",
			Valid:  false,
		},
		Email: pgtype.Text{
			String: "",
			Valid:  false,
		},
	}

	// Update FullName if provided
	if req.FullName != nil {
		arg.FullName.String = *req.FullName
		arg.FullName.Valid = true
	}

	// Update Email if provided
	if req.Email != nil {
		arg.Email.String = *req.Email
		arg.Email.Valid = true
	}

	// Check if Password is provided, hash it and assign it
	if req.Password != nil {
		hashedPassword, err := util.HashPassword(*req.Password)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
		}

		arg.HashedPassword = pgtype.Text{
			String: hashedPassword,
			Valid:  true,
		}

		arg.PasswordChangedAt = pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		status.Errorf(codes.Internal, "failed to update user, %s", err)
	}

	rsp := &pb.UpdateUserResponse{
		// need to convert to pb.User
		User: convertUser(user),
	}

	return rsp, nil
}

func ValidateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	// check other fields is not nil
	if req.Password != nil {
		if err := val.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}

	if req.FullName != nil {
		if err := val.ValidateFullName(req.GetFullName()); err != nil {
			violations = append(violations, fieldViolation("full_name", err))
		}
	}

	if req.Email != nil {
		if err := val.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}

	return violations

}
