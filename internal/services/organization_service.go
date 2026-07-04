package services

import (
	"context"
	"forgeflow/internal/database"
	"forgeflow/internal/repositories"

	"github.com/google/uuid"
)

type OrganizationService interface {
	CreateOrganization(ctx context.Context, name string, ownerID uuid.UUID) (*database.Organization, error)
	AddMember(ctx context.Context, orgID, userID uuid.UUID, role database.OrganizationRole) error
	RemoveMember(ctx context.Context, orgID, userID uuid.UUID) error
	VerifyRBACPermissions(ctx context.Context, orgID, userID uuid.UUID, requiredRole database.OrganizationRole) (bool, error)
	ListOrganizationsForUser(ctx context.Context, userID uuid.UUID) ([]database.Organization, error)
}

type organizationService struct {
	orgRepo repositories.OrganizationRepository
}

func NewOrganizationService(orgRepo repositories.OrganizationRepository) OrganizationService {
	return &organizationService{orgRepo: orgRepo}
}

func (s *organizationService) CreateOrganization(ctx context.Context, name string, ownerID uuid.UUID) (*database.Organization, error) {
	org := &database.Organization{
		Name: name,
		Members: []database.OrganizationMember{
			{
				UserID: ownerID,
				Role:   database.OrganizationRoleOwner,
			},
		},
	}
	if err := s.orgRepo.Create(ctx, org); err != nil {
		return nil, ErrInternalError
	}
	return org, nil
}

func (s *organizationService) AddMember(ctx context.Context, orgID, userID uuid.UUID, role database.OrganizationRole) error {
	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return ErrNotFound
	}

	for _, m := range org.Members {
		if m.UserID == userID {
			return ErrDuplicateResource
		}
	}

	org.Members = append(org.Members, database.OrganizationMember{
		OrganizationID: orgID,
		UserID:         userID,
		Role:           role,
	})

	if err := s.orgRepo.Update(ctx, org); err != nil {
		return ErrInternalError
	}
	return nil
}

func (s *organizationService) RemoveMember(ctx context.Context, orgID, userID uuid.UUID) error {
	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return ErrNotFound
	}

	var newMembers []database.OrganizationMember
	found := false
	for _, m := range org.Members {
		if m.UserID == userID {
			found = true
			continue
		}
		newMembers = append(newMembers, m)
	}

	if !found {
		return ErrNotFound
	}

	org.Members = newMembers
	if err := s.orgRepo.Update(ctx, org); err != nil {
		return ErrInternalError
	}
	return nil
}

func (s *organizationService) VerifyRBACPermissions(ctx context.Context, orgID, userID uuid.UUID, requiredRole database.OrganizationRole) (bool, error) {
	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return false, ErrNotFound
	}

	// Simplify RBAC hierarchy: Owner > Maintainer > Developer > Viewer
	rolePriority := map[database.OrganizationRole]int{
		database.OrganizationRoleViewer:     1,
		database.OrganizationRoleDeveloper:  2,
		database.OrganizationRoleMaintainer: 3,
		database.OrganizationRoleOwner:      4,
	}

	for _, m := range org.Members {
		if m.UserID == userID {
			if rolePriority[m.Role] >= rolePriority[requiredRole] {
				return true, nil
			}
			return false, nil
		}
	}
	return false, nil
}

func (s *organizationService) ListOrganizationsForUser(ctx context.Context, userID uuid.UUID) ([]database.Organization, error) {
	// Simple implementation; in reality, we'd query via a specific repo method.
	orgs, err := s.orgRepo.List(ctx, 100, 0)
	if err != nil {
		return nil, ErrInternalError
	}

	var userOrgs []database.Organization
	for _, org := range orgs {
		for _, m := range org.Members {
			if m.UserID == userID {
				userOrgs = append(userOrgs, org)
				break
			}
		}
	}
	return userOrgs, nil
}
