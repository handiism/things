ALTER TABLE "role"
    DROP CONSTRAINT IF EXISTS fk_credential_role;

ALTER TABLE "role_ability"
    DROP CONSTRAINT IF EXISTS fk_role_in_role_ability;

ALTER TABLE "role_ability"
    DROP CONSTRAINT IF EXISTS fk_ability_in_role_ability;

ALTER TABLE "role_ability"
    DROP CONSTRAINT IF EXISTS uq_role_ability;

DROP TABLE IF EXISTS "credential";

DROP TABLE IF EXISTS "role";

DROP TABLE IF EXISTS "role_ability";

DROP TABLE IF EXISTS "ability";

DROP TYPE IF EXISTS "ability_enum";

DROP TYPE IF EXISTS "role_enum";