package job

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository"
)

// JobUseCase は求人情報に関するユースケースを定義するインターフェースです
type JobUseCase interface {
	// GetJobs は求人情報の一覧を取得します
	GetJobs(ctx context.Context, page, limit int) (*entity.JobResponse, error)

	// GetJob は指定されたIDの求人情報を取得します
	GetJob(ctx context.Context, id int) (*entity.JobPosting, error)

	// CreateJob は新しい求人情報を作成します
	CreateJob(ctx context.Context, job *entity.JobPosting) error

	// UpdateJob は既存の求人情報を更新します
	UpdateJob(ctx context.Context, job *entity.JobPosting) error

	// DeleteJob は指定されたIDの求人情報を削除します
	DeleteJob(ctx context.Context, id int) error

	// GetJobsByCompanyID は指定された企業IDの求人情報を取得します
	GetJobsByCompanyID(ctx context.Context, companyID, page, limit int) (*entity.JobResponse, error)

	// GetCompanyWithJobs は指定された企業IDの会社情報と求人情報を取得します
	GetCompanyWithJobs(ctx context.Context, companyID int) (*entity.Company, []entity.JobPosting, error)
}

// jobUseCase は求人情報に関するユースケースの実装です
type jobUseCase struct {
	jobRepo        repository.JobRepository
	getCompanyRepo repository.GetCompanyRepository
}

// NewJobUseCase は求人ユースケースの新しいインスタンスを作成します
func NewJobUseCase(jobRepo repository.JobRepository, getCompanyRepo repository.GetCompanyRepository) JobUseCase {
	return &jobUseCase{
		jobRepo:        jobRepo,
		getCompanyRepo: getCompanyRepo,
	}
}

// GetJobs は求人情報の一覧を取得します
func (u *jobUseCase) GetJobs(ctx context.Context, page, limit int) (*entity.JobResponse, error) {
	// ページとリミットのデフォルト値を設定
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	// リポジトリから求人情報を取得
	jobs, total, err := u.jobRepo.FindAll(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	// レスポンスを作成
	response := &entity.JobResponse{
		Jobs:  jobs,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	return response, nil
}

// GetJob は指定されたIDの求人情報を取得します
func (u *jobUseCase) GetJob(ctx context.Context, id int) (*entity.JobPosting, error) {
	return u.jobRepo.FindByID(ctx, id)
}

// CreateJob は新しい求人情報を作成します
func (u *jobUseCase) CreateJob(ctx context.Context, job *entity.JobPosting) error {
	return u.jobRepo.Create(ctx, job)
}

// UpdateJob は既存の求人情報を更新します
func (u *jobUseCase) UpdateJob(ctx context.Context, job *entity.JobPosting) error {
	return u.jobRepo.Update(ctx, job)
}

// DeleteJob は指定されたIDの求人情報を削除します
func (u *jobUseCase) DeleteJob(ctx context.Context, id int) error {
	return u.jobRepo.Delete(ctx, id)
}

// GetJobsByCompanyID は指定された企業IDの求人情報を取得します
func (u *jobUseCase) GetJobsByCompanyID(ctx context.Context, companyID, page, limit int) (*entity.JobResponse, error) {
	// ページとリミットのデフォルト値を設定
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	// リポジトリから求人情報を取得
	jobs, total, err := u.jobRepo.FindByCompanyID(ctx, companyID, page, limit)
	if err != nil {
		return nil, err
	}

	// レスポンスを作成
	response := &entity.JobResponse{
		Jobs:  jobs,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	return response, nil
}

// GetCompanyWithJobs は指定された企業IDの会社情報と求人情報を取得します
func (u *jobUseCase) GetCompanyWithJobs(ctx context.Context, companyID int) (*entity.Company, []entity.JobPosting, error) {
	company, err := u.getCompanyRepo.FindByID(ctx, companyID)
	if err != nil {
		return nil, nil, err
	}

	jobs, _, err := u.jobRepo.FindByCompanyID(ctx, companyID, 1, 100)
	if err != nil {
		return nil, nil, err
	}

	return company, jobs, nil
}
