package enrollment

import (
	"context"
	"log"

	"github.com/MartinZitterkopf/gocurse_domain/domain"

	courseSdk "github.com/MartinZitterkopf/gocurse_sdk-microservice-/curse"
	userSdk "github.com/MartinZitterkopf/gocurse_sdk-microservice-/user"
)

type (
	Filters struct {
		UserID   string
		CourseID string
	}

	Service interface {
		Create(ctx context.Context, userID, courseID string) (*domain.Enrollment, error)
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Enrollment, error)
		Update(ctx context.Context, id string, status *string) error
		Count(ctx context.Context, filters Filters) (int, error)
	}

	service struct {
		log         *log.Logger
		userTrans   userSdk.Transport
		courseTrans courseSdk.Transport
		repo        Repository
	}
)

func NewService(l *log.Logger, userTrans userSdk.Transport, courseTrans courseSdk.Transport, repo Repository) Service {
	return &service{
		log:         l,
		userTrans:   userTrans,
		courseTrans: courseTrans,
		repo:        repo,
	}
}

func (s service) Create(ctx context.Context, userID, courseID string) (*domain.Enrollment, error) {

	enroll := &domain.Enrollment{
		UserID:  userID,
		CurseID: courseID,
		Status:  domain.Pending,
	}

	if _, err := s.userTrans.Get(userID); err != nil {
		return nil, err
	}

	if _, err := s.courseTrans.Get(courseID); err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, enroll); err != nil {
		return nil, err
	}

	return enroll, nil
}

func (s service) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Enrollment, error) {
	enrollments, err := s.repo.GetAll(ctx, filters, offset, limit)
	if err != nil {
		return nil, err
	}
	return enrollments, nil
}

func (s service) Update(ctx context.Context, id string, status *string) error {

	if status != nil {
		switch domain.EnrollStatus(*status) {
		case domain.Pending, domain.Active, domain.Studying, domain.Inactive:
		default:
			return ErrInvalidStatus{*status}
		}
	}

	if err := s.repo.Update(ctx, id, status); err != nil {
		return err
	}

	return nil
}

func (s service) Count(ctx context.Context, filters Filters) (int, error) {
	return s.repo.Count(ctx, filters)
}
