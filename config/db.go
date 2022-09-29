package config

import (
	"log"
	"os"

	controllers "github.com/cavdy-play/go_db/controllers"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

// Connecting to db
func Connect() *pg.DB {
	opts := &pg.Options{
		User:     "krishna",
		Password: "12345",
		Addr:     "localhost:5432",
		Database: "postgres",
	}

	var db *pg.DB = pg.Connect(opts)
	if db == nil {
		log.Printf("Failed to connect")
		os.Exit(100)
	}
	log.Printf("Connected to Database")
	controllers.CreateAccountTable(db)

	//Added later
	err := createSchema(db)
	if err != nil {
		panic(err)
	}

	controllers.InitiateDB(db)
	return db
}

//To create all database at once
func createSchema(db *pg.DB) error {

	orm.RegisterTable((*controllers.Acc)(nil))

	models := []interface{}{
		(*controllers.Bank)(nil),      // First in the sequence of tables.
		(*controllers.Customers)(nil), // Second
		(*controllers.Acc)(nil),
		(*controllers.Transactions)(nil),
	}

	opt := &orm.CreateTableOptions{

		IfNotExists:   true,
		FKConstraints: true,
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(opt)
		if err != nil {
			return err
		}
	}

	return nil
}