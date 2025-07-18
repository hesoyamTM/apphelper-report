package report

import (
	"context"

	"github.com/google/uuid"
	"github.com/hesoyamTM/apphelper-report/internal/models"

	reportv1 "github.com/hesoyamTM/apphelper-protos/gen/go/report"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Report interface {
	CreateReport(ctx context.Context, groupId, studentId, trainerId uuid.UUID, desc string) error
	GetReports(ctx context.Context, groupId, studentId, trainerId uuid.UUID) ([]models.Report, error)
}

type serverAPI struct {
	reportv1.UnimplementedReportServer
	reportService Report
}

func RegisterServer(gRpc *grpc.Server, service Report) {
	reportv1.RegisterReportServer(gRpc, &serverAPI{reportService: service})
}

func (s *serverAPI) CreateReport(ctx context.Context, req *reportv1.CreateReportRequest) (*reportv1.Empty, error) {
	studentId, err := uuid.Parse(req.GetStudentId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "validation error")
	}
	trainerId, err := uuid.Parse(req.GetTrainerId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "validation error")
	}
	groupId, err := uuid.Parse(req.GetGroupId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "validation error")
	}
	desc := req.GetDescription()

	if err := CheckIdPermission(ctx, studentId, trainerId); err != nil {
		return nil, err
	}

	if err := s.reportService.CreateReport(ctx, groupId, studentId, trainerId, desc); err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &reportv1.Empty{}, nil
}

func (s *serverAPI) GetReports(ctx context.Context, req *reportv1.GetReportsRequest) (*reportv1.GetReportsResponse, error) {
	studentId, err := uuid.Parse(req.GetStudentId())
	if err != nil {
		studentId = uuid.Nil
	}
	trainerId, err := uuid.Parse(req.GetTrainerId())
	if err != nil {
		trainerId = uuid.Nil
	}
	groupId, err := uuid.Parse(req.GetGroupId())
	if err != nil {
		groupId = uuid.Nil
	}

	if err := CheckIdPermission(ctx, studentId, trainerId); err != nil {
		return nil, err
	}

	if studentId == uuid.Nil && trainerId == uuid.Nil {
		return nil, status.Error(codes.InvalidArgument, "validation error")
	}

	reports, err := s.reportService.GetReports(ctx, groupId, studentId, trainerId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	res := make([]*reportv1.Report, len(reports))
	for i := range reports {
		res[i] = &reportv1.Report{
			GroupId:     reports[i].GroupId.String(),
			StudentId:   reports[i].StudentId.String(),
			TrainerId:   reports[i].TrainerId.String(),
			Description: reports[i].Description,
			Date:        timestamppb.New(reports[i].Date),
		}
	}

	return &reportv1.GetReportsResponse{
		Reports: res,
	}, nil
}
