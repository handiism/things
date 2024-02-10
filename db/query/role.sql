-- name: GetRoles :many
SELECT *
FROM "role";

-- name: GetRoleById :one
SELECT *
FROM "role"
WHERE "id" = $1;

-- name: GetRoleByName :one
SELECT *
FROM "role"
WHERE "name" = $1;

-- name: CreateRole :one
INSERT INTO "role"("name")
VALUES ($1)
RETURNING *;

-- name: UpsertRole :one
INSERT INTO "role"("name")
VALUES ($1)
ON CONFLICT ("name") DO UPDATE SET "name" = $1
RETURNING *;

-- name: UpdateRole :one
UPDATE "role"
SET "name" = $2
WHERE "id" = $1
RETURNING *;

-- name: SetRoleAbilities :many
INSERT INTO role_ability (role_id, ability_id)
SELECT sqlc.arg('role_id')::int as role_id, id as ability_id
FROM ability
WHERE name = ANY (sqlc.arg('role_abilities')::text[]::ability_enum[])
RETURNING *;

-- name: DeleteRole :exec
DELETE
FROM "role"
WHERE "id" = $1;
