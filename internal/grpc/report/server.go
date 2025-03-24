package report

import (
	"context"

	"github.com/hesoyamTM/apphelper-report/internal/models"

	reportv1 "github.com/hesoyamTM/apphelper-protos/gen/go/report"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Report interface {
	CreateReport(ctx context.Context, groupId, studentId, trainerId int64, desc string) error
	GetReports(ctx context.Context, groupId, studentId, trainerId int64) ([]models.Report, error)
}

type serverAPI struct {
	reportv1.UnimplementedReportServer
	reportService Report
}

func RegisterServer(gRpc *grpc.Server, service Report) {
	reportv1.RegisterReportServer(gRpc, &serverAPI{reportService: service})
}

func (s *serverAPI) CreateReport(ctx context.Context, req *reportv1.CreateReportRequest) (*reportv1.Empty, error) {
	studentId := req.GetStudentId()
	trainerId := req.GetTrainerId()
	groupId := req.GetGroupId()
	desc := req.GetDescription()

	if err := CheckIdPermission(ctx, studentId, trainerId); err != nil {
		return nil, err
	}

	// validation
	if studentId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "validation error")
	}
	if trainerId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "validation error")
	}

	if err := s.reportService.CreateReport(ctx, groupId, studentId, trainerId, desc); err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &reportv1.Empty{}, nil
}

func (s *serverAPI) GetReports(ctx context.Context, req *reportv1.GetReportsRequest) (*reportv1.GetReportsResponse, error) {
	studentId := req.GetStudentId()
	trainerId := req.GetTrainerId()
	groupId := req.GetGroupId()

	if err := CheckIdPermission(ctx, studentId, trainerId); err != nil {
		return nil, err
	}

	if studentId <= 0 && trainerId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "validation error")
	}

	reports, err := s.reportService.GetReports(ctx, groupId, studentId, trainerId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	res := make([]*reportv1.Report, len(reports))
	for i := range reports {
		res[i] = &reportv1.Report{
			GroupId:     reports[i].GroupId,
			StudentId:   reports[i].StudentId,
			TrainerId:   reports[i].TrainerId,
			Description: reports[i].Description,
			Date:        timestamppb.New(reports[i].Date),
		}
	}

	return &reportv1.GetReportsResponse{
		Reports: res,
	}, nil
}
