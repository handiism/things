package main

import (
	"context"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/handiism/smi/api"
	"github.com/handiism/smi/db/sqlc"
	"github.com/handiism/smi/util"
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

	migrateDatabase(cfg.MigrationPath, cfg.PostgresURL)
	seedDatabase(pool)
	runHttpServer(cfg, pool)
}

func seedDatabase(pool *pgxpool.Pool) error {
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
			return err
		}

		abilities = append(abilities, ability)
	}

	roleEnums := []sqlc.RoleEnum{
		sqlc.RoleEnumAdmin,
	}
	for _, roleEnum := range roleEnums {
		role, err := queries.UpsertRole(context.Background(), roleEnum)
		if err != nil {
			return err
		}

		for _, ability := range abilities {
			_, err := queries.UpsertRoleAbility(context.Background(), sqlc.UpsertRoleAbilityParams{RoleID: role.ID, AbilityID: ability.ID})
			if err != nil {
				return err
			}

		}

	}

	log.Info().Msg("db seeded successfully")

	return nil
}

func migrateDatabase(migrationPath string, postgresUrl string) {
	mgr, err := migrate.New(migrationPath, postgresUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to run new migration instance")
	}

	if err := mgr.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}

	log.Info().Msg("db migrated successfully")
}

func runHttpServer(config util.Config, pool *pgxpool.Pool) {

	server, err := api.NewServer(config, pool)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}
}
