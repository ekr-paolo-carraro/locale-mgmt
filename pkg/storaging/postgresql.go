package storaging

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/lib/pq"
)

//LocalePersistencer interface for persistence service
type LocalePersistencer interface {
	PostLocaleItem(item LocaleItem) (*LocaleItem, error)
	//PostLocaleItems(items []LocaleItem) (int, error)
	GetLocaleItem(key, bundle, lang string) (*LocaleItem, error)
	//GetLocaleItems(bundle, lang string) ([]LocaleItem, error)
}

//LocalePersistenceService manages persistence with db
type LocalePersistenceService struct {
	DBDelegate *sql.DB
}

//NewPostgresPersistenceService return a new persistence service for postgresql db
func NewPostgresPersistenceService() (*LocalePersistenceService, error) {
	connStr := os.Getenv("DATABASE_URL")
	log.Println(connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	initStmt, err := ioutil.ReadFile("pkg/storaging/sql/init.sql")
	_, err = db.Exec(string(initStmt))
	if err != nil {
		return nil, err
	}

	lps := LocalePersistenceService{db}

	return &lps, nil
}

//PostLocaleItem implements LocalePersistencer interface with postgresql implementation
func (lps LocalePersistenceService) PostLocaleItem(item LocaleItem) (*LocaleItem, error) {

	insertStmtStr, err := ioutil.ReadFile("pkg/storaging/sql/upsert.sql")
	if err != nil {
		return nil, err
	}

	insertStmt, err := lps.DBDelegate.Prepare(string(insertStmtStr))
	if err != nil {
		return nil, err
	}

	insertResult := insertStmt.QueryRow(item.Key, item.Bundle, item.Lang, item.Content)

	err = insertResult.Scan(&item.ID)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

//GetLocaleItem return one localeitem for key, bundle, lang
func (lps LocalePersistenceService) GetLocaleItem(key, bundle, lang string) (*LocaleItem, error) {
	selectStmt := "SELECT id, bundle, lang, key, content FROM localeitems WHERE localeitems.key = $1 AND localeitems.bundle = $2 AND localeitems.lang = $3;"
	sqlResult, err := lps.DBDelegate.Query(selectStmt)
	if err != nil {
		return nil, err
	}
	defer sqlResult.Close()

	items, err := parseResult(sqlResult)
	if err != nil {
		return nil, err
	}

	li := items[0]

	return &li, nil
}

func parseResult(res *sql.Rows) ([]LocaleItem, error) {
	result := make([]LocaleItem, 0)

	if res == nil {
		return nil, errors.New("Error on query result: no item to parse")
	}

	for res.Next() {
		var li LocaleItem
		err := res.Scan(
			&li.ID,
			&li.Bundle,
			&li.Lang,
			&li.Key,
			&li.Content,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, li)
	}

	if len(result) == 0 {
		return nil, errors.New("Error on query result: zero items")
	}

	return result, nil
}
