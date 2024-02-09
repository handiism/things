-- name: CreateCredential :one
INSERT INTO "credential"("email", "password", "role_id")
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetCredentials :many
SELECT *
FROM "credential";

-- name: GetCredentialById :one
SELECT *
FROM "credential"
WHERE "credential"."id" = $1;

-- name: GetCredentialByUsername :one
SELECT *
FROM "credential"
WHERE "username" = $1;

-- name: GetCredentialPasswordByEmail :one
SELECT "password"
FROM "credential"
WHERE "email" = $1;

-- name: SetPictureByCredentialId :one
UPDATE "credential"
SET "picture" = $2
WHERE "id" = $1
RETURNING *;

-- name: GetCredentialByEmail :one
SELECT *
FROM "credential"
WHERE "email" = $1;

-- name: UpdateCredential :exec
UPDATE "credential"
SET "name"     = $2,
    "email"    = $2,
    "username" = $3,
    "role_id"  = $4
WHERE "id" = $1;

-- name: DeleteCredential :exec
DELETE
FROM "credential"
WHERE "id" = $1;
