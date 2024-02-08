-- name: CreateCredential :one
INSERT INTO "credential"("email", "password")
VALUES ($1, $2)
RETURNING *;

-- name: GetCredentials :many
SELECT *
FROM "credential";

-- name: GetCredentialById :one
SELECT sqlc.embed(credential), sqlc.embed(role)
FROM "credential"
         JOIN "role" ON "credential"."id" = "role"."id"
WHERE "credential"."id" = $1;

-- name: GetCredentialByUsername :one
SELECT *
FROM "credential"
WHERE "username" = $1;

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
