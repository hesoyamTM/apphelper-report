package report

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func CheckIdPermission(ctx context.Context, ids ...uuid.UUID) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Internal, "internal error")
	}

	mdUid, ok := md["uid"]
	if !ok {
		return status.Error(codes.Unauthenticated, "internal error")
	}

	uid, err := uuid.Parse(mdUid[0])
	if err != nil {
		return status.Error(codes.InvalidArgument, "validation error")
	}

	if len(ids) == 0 {
		return status.Error(codes.Internal, "internal error")
	}

	for _, id := range ids {
		if id == uid {
			return nil
		}
	}

	return status.Error(codes.PermissionDenied, "permission denied")
}
