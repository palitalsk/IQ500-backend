package services

import (
	"context"
	"fmt"
	"time"

	"main/helpers"
	"main/models"
	"main/repo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BusinessService struct {
	embedRepo        *repo.EmbedRepository
	pineconeRepo     *repo.PineconeRepository
	businessRepo     *repo.BusinessRepository
	businessUserRepo *repo.BusinessUserRepository
}

func NewBusinessService(embedRepo *repo.EmbedRepository, pineconeRepo *repo.PineconeRepository, businessRepo *repo.BusinessRepository, businessUserRepo *repo.BusinessUserRepository) *BusinessService {
	return &BusinessService{
		embedRepo:        embedRepo,
		pineconeRepo:     pineconeRepo,
		businessRepo:     businessRepo,
		businessUserRepo: businessUserRepo,
	}
}

// CreateBusiness
func (s *BusinessService) CreateBusiness(ctx context.Context, userID primitive.ObjectID, req *models.CreateBusinessRequest) (*models.BusinessResponse, error) {
	// Generate unique business code
	businessCode, err := helpers.GenerateBusinessCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate business code: %w", err)
	}

	// Check if business code already exists
	exists, err := s.businessRepo.BusinessExists(ctx, businessCode)
	if err != nil {
		return nil, fmt.Errorf("failed to check business code: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("business code already exists")
	}

	business := &models.Business{
		Code:         businessCode,
		Name:         req.Name,
		BusinessType: req.BusinessType,
		Detail:       req.Detail,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save to DB
	err = s.businessRepo.CreateBusiness(ctx, business)
	if err != nil {
		return nil, fmt.Errorf("failed to create business: %w", err)
	}

	// Create business-user
	businessUser := &models.BusinessUser{
		Username:     "",
		UserID:       userID,
		BusinessCode: businessCode,
		BusinessID:   business.ID,
		Role:         "owner",
	}

	err = s.businessUserRepo.CreateBusinessUser(ctx, businessUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create business-user relationship: %w", err)
	}

	return &models.BusinessResponse{
		ID:           business.ID,
		Code:         business.Code,
		Name:         business.Name,
		BusinessType: business.BusinessType,
		Detail:       business.Detail,
		CreatedAt:    business.CreatedAt,
	}, nil
}

// GetUserBusinesses ดึงข้อมูลธุรกิจของ user_id นั้น
func (s *BusinessService) GetUserBusinesses(ctx context.Context, userID primitive.ObjectID) (*models.BusinessListResponse, error) {
	// Get business-user relationships
	businessUsers, err := s.businessUserRepo.GetBusinessUsersByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user businesses: %w", err)
	}

	var businesses []models.BusinessResponse
	for _, businessUser := range businessUsers {
		business, err := s.businessRepo.GetBusinessByID(ctx, businessUser.BusinessID)
		if err != nil {
			continue // Skip if business not found
		}

		businesses = append(businesses, models.BusinessResponse{
			ID:           business.ID,
			Code:         business.Code,
			Name:         business.Name,
			BusinessType: business.BusinessType,
			Detail:       business.Detail,
			CreatedAt:    business.CreatedAt,
		})
	}

	return &models.BusinessListResponse{
		Businesses: businesses,
		Total:      len(businesses),
	}, nil
}

// GetBusiness ดึงข้อมูลโดยใช้ businessID
func (s *BusinessService) GetBusiness(ctx context.Context, businessID primitive.ObjectID) (*models.BusinessResponse, error) {
	business, err := s.businessRepo.GetBusinessByID(ctx, businessID)
	if err != nil {
		return nil, fmt.Errorf("failed to get business: %w", err)
	}

	return &models.BusinessResponse{
		ID:           business.ID,
		Code:         business.Code,
		Name:         business.Name,
		BusinessType: business.BusinessType,
		Detail:       business.Detail,
		CreatedAt:    business.CreatedAt,
	}, nil
}

// GetBusinessByCode ดึงข้อมูลโดยใช้ business_code
func (s *BusinessService) GetBusinessByCode(ctx context.Context, businessCode string) (*models.BusinessResponse, error) {
	business, err := s.businessRepo.GetBusinessByCode(ctx, businessCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get business: %w", err)
	}

	return &models.BusinessResponse{
		ID:           business.ID,
		Code:         business.Code,
		Name:         business.Name,
		BusinessType: business.BusinessType,
		Detail:       business.Detail,
		CreatedAt:    business.CreatedAt,
	}, nil
}
