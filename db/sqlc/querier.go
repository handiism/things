// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CreateAbilities(ctx context.Context, name []AbilityEnum) (int64, error)
	CreateAbility(ctx context.Context, name AbilityEnum) (Ability, error)
	CreateAbilityWithoutConflict(ctx context.Context, name AbilityEnum) (Ability, error)
	CreateCredential(ctx context.Context, arg CreateCredentialParams) (Credential, error)
	CreateRole(ctx context.Context, name RoleEnum) (Role, error)
	CreateRoleAbility(ctx context.Context, arg CreateRoleAbilityParams) (RoleAbility, error)
	DeleteAbility(ctx context.Context, id int32) error
	DeleteCredential(ctx context.Context, id pgtype.UUID) error
	DeleteRole(ctx context.Context, id int32) error
	DeleteRoleAbility(ctx context.Context, id int32) error
	GetAbilities(ctx context.Context) ([]Ability, error)
	GetAbilityById(ctx context.Context, id int32) (Ability, error)
	GetCredentialByEmail(ctx context.Context, email string) (Credential, error)
	GetCredentialById(ctx context.Context, id pgtype.UUID) (GetCredentialByIdRow, error)
	GetCredentialByUsername(ctx context.Context, username pgtype.Text) (Credential, error)
	GetCredentials(ctx context.Context) ([]Credential, error)
	GetRoleAbilities(ctx context.Context) ([]RoleAbility, error)
	GetRoleById(ctx context.Context, id int32) (Role, error)
	GetRoleByName(ctx context.Context, name RoleEnum) (Role, error)
	GetRoles(ctx context.Context) ([]Role, error)
	UpdateAbility(ctx context.Context, arg UpdateAbilityParams) error
	UpdateCredential(ctx context.Context, arg UpdateCredentialParams) error
	UpdateRole(ctx context.Context, arg UpdateRoleParams) error
	UpdateRoleAbility(ctx context.Context, arg UpdateRoleAbilityParams) error
	UpsertRole(ctx context.Context, name RoleEnum) (Role, error)
	UpsertRoleAbility(ctx context.Context, arg UpsertRoleAbilityParams) (RoleAbility, error)
}

var _ Querier = (*Queries)(nil)