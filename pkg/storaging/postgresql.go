package storaging

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

//LocalePersistencer interface for persistence service
type LocalePersistencer interface {
	PostLocaleItem(item LocaleItem) (*LocaleItem, error)
	PostLocaleItems(items []LocaleItem) (int, error)
	GetLocaleItems(key, bundle, lang string) ([]LocaleItem, error)
	GetLangs() ([]string, error)
	GetBundles() ([]string, error)
}

//LocalePersistenceService manages persistence with db
type LocalePersistenceService struct {
	DBDelegate *sql.DB
}

//NewPostgresPersistenceService return a new persistence service for postgresql db
func NewPostgresPersistenceService() (*LocalePersistenceService, error) {
	connStr := os.Getenv("DATABASE_URL")
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
	insertResult.Scan(&item.ID)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

//PostLocaleItems implements LocalePersistencer interface with postgresql implementation
func (lps LocalePersistenceService) PostLocaleItems(items []LocaleItem) (int, error) {

	insertStmtStr, err := ioutil.ReadFile("pkg/storaging/sql/upsert.sql")
	if err != nil {
		return 0, err
	}

	tx, err := lps.DBDelegate.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	insertStmt, err := tx.Prepare(string(insertStmtStr))
	if err != nil {
		return 0, err
	}

	defer insertStmt.Close()
	for _, item := range items {
		if _, err = insertStmt.Exec(item.Key, item.Bundle, item.Lang, item.Content); err != nil {
			return 0, err
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return len(items), nil
}

//GetLocaleItem return one localeitem for key, bundle, lang
func (lps LocalePersistenceService) GetLocaleItems(key, bundle, lang string) ([]LocaleItem, error) {
	selectStmt := "SELECT id, bundle, lang, key, content FROM localeitems WHERE"

	placeHolderCounter := 0
	params := []interface{}{}
	if key != "" {
		placeHolderCounter++
		selectStmt += " localeitems.key = $" + strconv.Itoa(placeHolderCounter) + " AND"
		params = append(params, key)
	}
	if bundle != "" {
		placeHolderCounter++
		selectStmt += " localeitems.bundle = $" + strconv.Itoa(placeHolderCounter) + " AND"
		params = append(params, bundle)
	}
	if lang != "" {
		placeHolderCounter++
		selectStmt += " localeitems.lang = $" + strconv.Itoa(placeHolderCounter) + " AND"
		params = append(params, lang)
	}

	selectStmt = strings.TrimSuffix(selectStmt, " AND")

	log.Println(selectStmt)

	sqlResult, err := lps.DBDelegate.Query(selectStmt, params...)
	if err != nil {
		return nil, err
	}
	defer sqlResult.Close()

	items, err := parseResult(sqlResult)
	if err != nil {
		return nil, err
	}

	return items, nil
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

func (lps LocalePersistenceService) GetLangs() ([]string, error) {
	result := []string{}
	rows, err := lps.DBDelegate.Query("SELECT DISTINCT(lang) FROM localeitems")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var lang string
		err = rows.Scan(&lang)
		if err != nil {
			return nil, err
		}

		result = append(result, lang)
	}

	return result, nil
}

func (lps LocalePersistenceService) GetBundles() ([]string, error) {
	result := []string{}
	rows, err := lps.DBDelegate.Query("SELECT DISTINCT(bundle) FROM localeitems")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var lang string
		err = rows.Scan(&lang)
		if err != nil {
			return nil, err
		}

		result = append(result, lang)
	}

	return result, nil
}
