package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	spew "github.com/davecgh/go-spew/spew"
	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/Liooo/golang-orm-samples/sqlboiler/models"
)

func main() {
	var (
		dbSchema   = os.Getenv("PSQL_SCHEMA")
		dbName     = os.Getenv("PSQL_DBNAME")
		dbUser     = os.Getenv("PSQL_USER")
		dbHost     = os.Getenv("PSQL_HOST")
		dbPort     = os.Getenv("PSQL_PORT")
		dbPassword = os.Getenv("PSQL_PASSWORD")
		dbSSLMode  = os.Getenv("PSQL_SSLMODE")
	)
	setDebug(true)

	pqOpts := fmt.Sprintf("host=%s port=%s user=%s password=%s schema=%s dbname=%s sslmode=%s", dbHost, dbPort, dbUser, dbPassword, dbSchema, dbName, dbSSLMode)
	println("connectiing... " + pqOpts + "\n")
	db, err := sql.Open(
		"postgres",
		pqOpts,
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = db.Close()
		dieOnErr(err)
	}()

	ctx := context.Background()
	setupSeed(ctx, db, false)

	// _, err = models.Pilots(
	// 	qm.Load(models.PilotRels.Jets),
	//     qm.Load(models.PilotRels.Languages),
	//     models.PilotWhere.Name.EQ("taro"),
	// ).All(ctx, db)
	// dieOnErr(err)
	// SELECT * FROM "pilots" WHERE ("pilots"."name" = $1);
	// [taro]
	// SELECT * FROM "jets" WHERE ("jets"."pilot_id" IN ($1));
	// [49]
	// SELECT "languages".id, "languages".language, "a"."pilot_id" FROM "languages" INNER JOIN "pilots_languages" as "a" on "languages"."id" = "a"."language_id" WHERE ("a"."pilot_id" IN ($1));
	// [49]

	jets, err := models.Jets(
		qm.Load(
			qm.Rels(
				models.JetRels.Pilot,
				models.PilotRels.Languages,
			),
			models.LanguageWhere.Language.EQ("English"),
			models.LanguageWhere.ID.NIN([]int{1, 2, 3}),
		),
	).All(ctx, db)
	dieOnErr(err)
	spew.Dump(jets)
}

func setDebug(debug bool) {
	boil.DebugMode = debug

	if debug {
		boil.DebugWriter = os.Stdout
	} else {
		boil.DebugWriter = nil
	}
}

func doTruncate(db *sql.DB, tbl string) error {
	_, err := queries.Raw(fmt.Sprintf("TRUNCATE TABLE %s CASCADE;", tbl)).Exec(db)
	return err
}

func setupSeed(ctx context.Context, db *sql.DB, showLog bool) {
	d := boil.IsDebug(ctx)
	setDebug(showLog)
	defer setDebug(d) // revert to whatever it was

	var err error
	_, err = queries.Raw("SET CONSTRAINTS ALL DEFERRED;").Exec(db)
	dieOnErr(err)

	defer func() {
		_, err := queries.Raw("SET CONSTRAINTS ALL IMMEDIATE;").Exec(db)
		dieOnErr(err)
	}()

	for _, t := range models.AllTableNames {
		err := doTruncate(db, t)
		dieOnErr(err)
	}

	dieOnErr(err)

	languages := []*models.Language{
		{Language: "Japanese"},
		{Language: "English"},
		{Language: "Hindi"},
	}
	for _, l := range languages {
		err := l.Insert(ctx, db, boil.Infer())
		dieOnErr(err)
	}

	names := []string{"a", "b", "c"}
	for _, name := range names {
		p := models.Pilot{Name: name}
		err := p.Insert(ctx, db, boil.Infer())
		dieOnErr(err)

		jets := []*models.Jet{
			{
				PilotID: p.ID,
				Age:     0,
				Name:    name + "-jet1",
				Color:   "blue",
			},
			{
				PilotID: p.ID,
				Age:     1,
				Name:    name + "-jet2",
				Color:   "blue",
			},
		}
		err = p.AddJets(ctx, db, true, jets...)
		dieOnErr(err)

		err = p.AddLanguages(ctx, db, false, languages...) // the 3rd `false` then creates only pilots_languages
		dieOnErr(err)
		// err = p.AddLanguages(ctx, db, false, languages...)
		// INSERT INTO "languages" ("id","language") VALUES ($1,$2)
		// [44 Japanese]
		// failed to insert into foreign table: models: unable to insert into languages: pq: duplicate key value violates unique constraint "languages_pkey"

		myLang := &models.Language{Language: name + "-lang"}
		err = p.AddLanguages(ctx, db, true, myLang) // the 3rd `true` then creates both language and pilots_languages
		dieOnErr(err)
		// err = p.AddLanguages(ctx, db, false, myLang)
		// INSERT INTO "pilots_languages" ("pilot_id", "language_id") VALUES ($1, $2)
		// [18 0]
		// failed to insert into foreign table: models: unable to insert into languages: pq: duplicate key value violates unique constraint "languages_pkey"
		// why duplicated key for id=0...? anyway
	}

}

func dieOnErr(err error) {
	if err != nil {
		println("---------- error happened!!! ----------")
		panic(err)
	}
}
