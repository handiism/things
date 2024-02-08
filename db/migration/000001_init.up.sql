CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE "role_enum" AS ENUM ('admin');

CREATE TYPE "ability_enum" AS ENUM ('credential_create','credential_read','credential_update','credential_delete');

CREATE TABLE IF NOT EXISTS "credential"
(
    "id"         uuid PRIMARY KEY             DEFAULT uuid_generate_v4(),
    "name"       varchar(255),
    "email"      varchar(320) UNIQUE NOT NULL,
    "username"   varchar(30) UNIQUE,
    "password"   char(60)            NOT NULL,
    "role_id"    int                 NOT NULL,
    "created_at" timestamptz         NOT NULL DEFAULT now(),
    "deleted_at" timestamptz
);

CREATE TABLE IF NOT EXISTS "role"
(
    "id"   serial PRIMARY KEY,
    "name" role_enum NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS "role_ability"
(
    "id"         serial PRIMARY KEY,
    "role_id"    int NOT NULL,
    "ability_id" int NOT NULL
);

CREATE TABLE IF NOT EXISTS "ability"
(
    "id"   serial PRIMARY KEY,
    "name" ability_enum NOT NULL UNIQUE
);


ALTER TABLE
    "credential"
    ADD
        CONSTRAINT fk_credential_role FOREIGN KEY ("role_id") REFERENCES "role" ("id");

ALTER TABLE
    "role_ability"
    ADD
        CONSTRAINT uq_role_ability UNIQUE("role_id","ability_id");


ALTER TABLE
    "role_ability"
    ADD
        CONSTRAINT fk_role_in_role_ability FOREIGN KEY ("role_id") REFERENCES "role" ("id");

ALTER TABLE
    "role_ability"
    ADD
        CONSTRAINT fk_ability_in_role_ability FOREIGN KEY ("ability_id") REFERENCES "ability" ("id");
