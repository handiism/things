-- name: CreateRoleAbility :one
INSERT INTO "role_ability"("role_id", "ability_id")
VALUES ($1, $2)
RETURNING *;

-- name: UpsertRoleAbility :one
INSERT INTO "role_ability"("role_id", "ability_id")
VALUES ($1, $2)
ON CONFLICT ("role_id","ability_id") DO UPDATE SET "role_id"    = $1,
                                                   "ability_id" = $2
RETURNING *;

-- name: GetRoleAbilities :many
SELECT *
FROM "role_ability";

-- name: UpdateRoleAbility :execresult
UPDATE "role_ability"
SET "role_id"    = $2,
    "ability_id" = $3
WHERE "id" = $1;


-- name: DeleteRoleAbilitiesByRoleId :execresult
DELETE
FROM "role_ability"
WHERE "role_id" = $1;
