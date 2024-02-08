-- name: GetAbilities :many
SELECT *
FROM "ability";

-- name: GetAbilityById :one
SELECT *
FROM "ability"
WHERE "id" = $1;

-- name: CreateAbility :one
INSERT INTO "ability"("name")
VALUES ($1)
RETURNING *;

-- name: CreateAbilities :copyfrom
INSERT INTO "ability"("name")
VALUES ($1);

-- name: CreateAbilityWithoutConflict :one
INSERT INTO "ability"("name")
VALUES ($1)
ON CONFLICT ("name") DO UPDATE SET "name" = $1
RETURNING *;

-- name: UpdateAbility :one
UPDATE "ability"
SET "name" = $2
WHERE "id" = $1
RETURNING *;

-- name: DeleteAbility :exec
DELETE
FROM "ability"
WHERE "id" = $1;
