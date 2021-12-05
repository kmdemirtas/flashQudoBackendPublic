package database

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/ironstone95/FlashQudoV2/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	l        *log.Logger
	db       *gorm.DB
	debugLog *log.Logger
}

// ConnectDB connects to the database with dsn configuration.
// If migrateDB is true, the tables will be created for given models.
// If deletePrevData is true, the data will be deleted but tables will remain.
// If deletePrevData and addTestData is true, the test data will be added.
// If fullLog is true, all actions will be logged.
func ConnectDB(l *log.Logger, migrateDB, deletePrevData, addTestData, fullLog bool) *Database {
	rd := Database{}
	rd.l = l
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_DATABASE_NAME")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432", host, user, password, dbname)
	dbLogger := logger.New(l, logger.Config{
		IgnoreRecordNotFoundError: true,
	})
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: dbLogger})
	if err != nil {
		l.Fatal(err)
	}
	rd.db = db
	models := []interface{}{&model.Token{}, &model.Member{}, &model.Card{}, &model.Bundle{}, &model.Group{}, &model.User{}}
	if migrateDB {
		err := rd.migrate(models)
		if err != nil {
			l.Fatal(err)
		}
	}
	if deletePrevData {
		rd.clear()
		if addTestData {
			rd.addTestData()
		}
	}

	if fullLog {
		rd.debugLog = log.New(os.Stdout, "[Database Error] ", log.Default().Flags())
		l.SetFlags(log.Lshortfile)
	}
	return &rd
}

// migrate creates tables inside the database with given models.
func (db *Database) migrate(models []interface{}) error {
	for _, m := range models {
		if err := db.db.Statement.AutoMigrate(m); err != nil {
			return err
		}
	}
	err := db.db.Statement.SetupJoinTable(&model.User{}, "Groups", &model.Member{})
	if err != nil {
		return err
	}

	err = db.db.Statement.SetupJoinTable(&model.Group{}, "Members", &model.Member{})
	if err != nil {
		return err
	}
	return nil
}

// clear deletes all data from the tables. queries must be written by hand.
func (db *Database) clear() {
	if err := db.db.Exec("Delete From tokens").Error; err != nil {
		db.l.Fatal(err)
	}
	if err := db.db.Exec("Delete From members").Error; err != nil {
		db.l.Fatal(err)
	}
	if err := db.db.Exec("Delete From cards").Error; err != nil {
		db.l.Fatal(err)
	}
	if err := db.db.Exec("Delete From bundles").Error; err != nil {
		db.l.Fatal(err)
	}
	if err := db.db.Exec("Delete From groups").Error; err != nil {
		db.l.Fatal(err)
	}
	if err := db.db.Exec("Delete From users").Error; err != nil {
		db.l.Fatal(err)
	}
}

// addTestData adds the test data to the database. Data must be written by hand.
func (db *Database) addTestData() {
	tokens := []model.Token{{Token: "token1", UserID: "user1", Expires: 1660403886},
		{Token: "token2", UserID: "user2", Expires: 1660403886},
		{Token: "token3", UserID: "user3", Expires: 1660403886},
	}
	if err := db.db.Create(&tokens).Error; err != nil {
		db.l.Fatal(err)
	}

	users := []model.User{{ID: "user1", Username: "user1", ImageURL: "Default"},
		{ID: "user2", Username: "user2", ImageURL: "Default"},
		{ID: "user3", Username: "user3", ImageURL: "Default"},
	}
	if err := db.db.Create(&users).Error; err != nil {
		db.l.Fatal(err)
	}

	groups := []model.Group{{ID: "group1", Name: "GroupName1"}}
	if err := db.db.Create(&groups).Error; err != nil {
		db.l.Fatal(err)
	}

	members := []model.Member{{GroupID: "group1", UserID: "user1", IsAdmin: true, MemberSince: time.Now()},
		{GroupID: "group1", UserID: "user2", IsAdmin: false, MemberSince: time.Now()},
	}
	if err := db.db.Create(&members).Error; err != nil {
		db.l.Fatal(err)
	}

	bundles := []model.Bundle{{
		ID:          "bundle1",
		Title:       "bTitle1",
		Description: "bDesc1",
		GroupID:     "group1",
	},
		{
			ID:          "bundle2",
			Title:       "bTitle2",
			Description: "bDesc1",
			GroupID:     "group1",
		},
	}

	if err := db.db.Create(&bundles).Error; err != nil {
		db.l.Fatal(err)
	}

	cards := []model.Card{}
	for i := 0; i < 50; i++ {
		c := model.Card{}
		err := faker.FakeData(&c)
		if err != nil {
			db.l.Fatal(err)
		}

		if rand.Int()%2 == 0 {
			c.BundleID = "bundle1"
		} else {
			c.BundleID = "bundle2"
		}

		cards = append(cards, c)
	}

	if err := db.db.Create(&cards).Error; err != nil {
		db.l.Fatal(err)
	}

}
