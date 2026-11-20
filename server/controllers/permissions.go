package controllers

import (
	"context"
	"errors"

	"github.com/zeiss/builder/server/ports"
	authz "github.com/zeiss/fiber-authz/v3"
	humax "github.com/zeiss/fiber-authz/v3/oas/huma"
)

var _ ports.Permissions = (*PermissionsController)(nil)

// ErrNotAllowed is returned when the user is not allowed to perform the given action on the given resource.
var ErrNotAllowed = errors.New("not allowed")

// PermissionsController is a controller for managing sites.
type PermissionsController struct {
	checker authz.Checker
}

// NewPermissionsController creates a new instance of PermissionsController.
func NewPermissionsController(checker authz.Checker) *PermissionsController {
	return &PermissionsController{
		checker: checker,
	}
}

// Allowed checks if the user is allowed to perform the given action on the given resource.
func (c *PermissionsController) Allowed(ctx context.Context, siteName, permission string) error {
	email, err := humax.GetClaim(ctx, "email")
	if err != nil {
		return ErrNotAllowed
	}

	allowed, err := c.checker.Allowed(ctx, authz.Principal(email), authz.Object(siteName), authz.Action(permission))
	if err != nil {
		return ErrNotAllowed
	}

	if !allowed {
		return ErrNotAllowed
	}

	return nil
}
