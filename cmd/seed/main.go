package main

import (
	"context"
	"os"

	"github.com/handiism/things/db/sqlc"
	"github.com/handiism/things/util"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if cfg.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	pool, err := pgxpool.New(context.Background(), cfg.PostgresURL)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	queries := sqlc.New(pool)

	abilityEnums := []sqlc.AbilityEnum{
		sqlc.AbilityEnumCredentialCreate,
		sqlc.AbilityEnumCredentialRead,
		sqlc.AbilityEnumCredentialUpdate,
		sqlc.AbilityEnumCredentialDelete,
	}

	abilities := []sqlc.Ability{}
	for _, abilityEnum := range abilityEnums {
		ability, err := queries.CreateAbilityWithoutConflict(context.Background(), abilityEnum)
		if err != nil {
			log.Fatal().Err(err)
		}

		abilities = append(abilities, ability)
	}

	roleEnums := []sqlc.RoleEnum{
		sqlc.RoleEnumAdmin,
	}
	for _, roleEnum := range roleEnums {
		role, err := queries.UpsertRole(context.Background(), roleEnum)
		if err != nil {
			log.Fatal().Err(err)
		}

		for _, ability := range abilities {
			_, err := queries.UpsertRoleAbility(context.Background(), sqlc.UpsertRoleAbilityParams{RoleID: role.ID, AbilityID: ability.ID})
			if err != nil {
				log.Fatal().Err(err)
			}
		}
	}

	log.Info().Msg("db seeded successfully")
}
