package handling

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/ekr-paolo-carraro/locale-mgmt/pkg/storaging"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"
)

type testItem struct {
	Name     string
	TextFunc func(*testing.T)
}

var r *gin.Engine

var localeItemToCompare []storaging.LocaleItem

func TestMain(m *testing.M) {

	var err error
	err = gotenv.Load()
	os.Setenv("test", "on")

	if err != nil {
		log.Panicln(err)
	}

	r, err = NewHandler()
	if err != nil {
		log.Panicln(err)
	}

	var compareJson string = `[{
		"id":"1",
		"bundle": "message",
		"key": "@ALERT_ERROR@",
		"lang": "it-IT",
		"content": "This is an error"
	}]`

	localeItemToCompare, err = buildDataToCompare([]byte(compareJson))
	if err != nil {
		log.Panicln(err)
	}

	os.Exit(m.Run())
}

func TestRoutes(t *testing.T) {

	apiTest := []testItem{
		{"welcome", testWelcome},
		{"post localeitem correct payload", testCorrectPostLocaleItem},
		{"post localeitem wrong payload", testWrongPostLocaleItem},
		{"post localeitem missing payload", testMissingPostLocaleItem},
		{"post localeitems correct payload", testCorrectPostLocaleItems},
		{"post localeitems wrong payload", testWrongPostLocaleItems},
		{"bundles", testBundles},
		{"langs", testLangs},
		{"get locale item by bundle", testGetLangByBundle},
		{"get locale item by bundle lang", testGetLangByBundleLang},
		{"get locale item by bundle key", testGetLangByBundleKey},
		{"delete locale item by bundle", testDeleteLangByBundle},
	}

	for _, ct := range apiTest {
		t.Run(ct.Name, ct.TextFunc)
	}
}

var welcomeMSG string = `{"Message":"Hi, welcome to locale-mgmt, don't know who you are so go to login"}`

func testWelcome(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, welcomeMSG, w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/welcome", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, welcomeMSG, w.Body.String())
}

func testBundles(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/bundles", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "label")
}

func testLangs(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/langs", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "it-IT")
}

func testCorrectPostLocaleItem(t *testing.T) {

	jdata, err := ioutil.ReadFile("test-data/single-item.json")
	if err != nil {
		t.Fatalf("error on load json file: %v\n", err)
	}

	bodyReader := strings.NewReader(string(jdata))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/locale-item", bodyReader)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func testWrongPostLocaleItem(t *testing.T) {
	jdata, err := ioutil.ReadFile("test-data/single-item-wrong.json")
	if err != nil {
		t.Fatalf("error on load json file: %v\n", err)
	}

	bodyReader := strings.NewReader(string(jdata))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/locale-item", bodyReader)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func testMissingPostLocaleItem(t *testing.T) {

	bodyReader := strings.NewReader("")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/locale-item", bodyReader)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func testCorrectPostLocaleItems(t *testing.T) {

	jdata, err := ioutil.ReadFile("test-data/multiple-item.json")
	if err != nil {
		t.Fatalf("error on load json file: %v\n", err)
	}

	bodyReader := strings.NewReader(string(jdata))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/locale-items", bodyReader)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `{
		"num_successfull": 4,
		"num_failed": 0
	}`, w.Body.String())
}

func testWrongPostLocaleItems(t *testing.T) {

	jdata, err := ioutil.ReadFile("test-data/multiple-item-wrong.json")
	if err != nil {
		t.Fatalf("error on load json file: %v\n", err)
	}

	bodyReader := strings.NewReader(string(jdata))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/locale-items", bodyReader)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `{
		"num_successfull": 0,
		"num_failed": 2
	}`, w.Body.String())
}

func buildDataToCompare(rawdata []byte) ([]storaging.LocaleItem, error) {
	comparingLocaleIten := make([]storaging.LocaleItem, 1)
	err := json.Unmarshal(rawdata, &comparingLocaleIten)
	if err != nil {
		return nil, err
	}
	return comparingLocaleIten, nil
}

func testGetLangByBundle(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/locale-items/message", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	gotData, err := buildDataToCompare(w.Body.Bytes())
	if err != nil {
		t.Fatalf("error on build json: %v\n", err)
	}

	assert.Equal(t, localeItemToCompare, gotData)
}

func testGetLangByBundleLang(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/locale-items/message/lang/it-IT", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	gotData, err := buildDataToCompare(w.Body.Bytes())
	if err != nil {
		t.Fatalf("error on build json: %v\n", err)
	}

	assert.Equal(t, localeItemToCompare, gotData)
}

func testGetLangByBundleKey(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/locale-items/message/key/@ALERT_ERROR@", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	gotData, err := buildDataToCompare(w.Body.Bytes())
	if err != nil {
		t.Fatalf("error on build json: %v\n", err)
	}

	assert.Equal(t, localeItemToCompare, gotData)
}

func testDeleteLangByBundle(t *testing.T) {

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/locale-items/message/lang/it-IT", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{
		"num_successfull": 1,
		"num_failed": 0
	}`, w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/locale-items/message/key/@ALERT_ERROR@", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{
		"num_successfull": 0,
		"num_failed": 0
	}`, w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/locale-items/message", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{
		"num_successfull": 0,
		"num_failed": 0
	}`, w.Body.String())
}
