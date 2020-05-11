package handling

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/ekr-paolo-carraro/locale-mgmt/pkg/session"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type testItem struct {
	Name     string
	TextFunc func(*testing.T)
}

type CustomClaims struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}

var router *gin.Engine

func retriveToken() {

	// err := gotenv.Load()
	// if err != nil {
	// 	log.Panicln(err)
	// }

	// url := os.Getenv("AUTH0_DOMAIN") + "oauth/token"
	// urlAudience := os.Getenv("AUTH0_DOMAIN") + "api/v2/"
	// payloadString := fmt.Sprintf("{\"client_id\":\"%v\",\"client_secret\":\"%v\",\"audience\":\"%v\",\"grant_type\":\"client_credentials\"}", os.Getenv("AUTH0_CLIENT_ID_M2M"), os.Getenv("AUTH0_CLIENT_SECRET_M2M"), urlAudience)
	// payload := strings.NewReader(payloadString)

	// req, _ := http.NewRequest("POST", url, payload)
	// req.Header.Add("content-type", "application/json")
	// res, _ := http.DefaultClient.Do(req)

	// ss, err := session.Store.Get(req, "auth-session")
	// if err != nil {
	// 	log.Fatalf("error on set session: %v", err.Error())
	// }

	// defer res.Body.Close()

	// responseData, _ := ioutil.ReadAll(res.Body)
	// accessToken := oauth2.Token{}
	// json.Unmarshal(responseData, &accessToken)

	session.InitSessionStorage()

	ss, err := session.Store.Get(nil, "auth-session")
	if err != nil {
		log.Fatalf("error on set session: %v", err.Error())
	}
	ss.Values["access_token"] = "fake"

	// jtoken, err := jwt.Parse(accessToken.AccessToken, func(tk *jwt.Token) (interface{}, error) {
	// 	if _, ok := tk.Method.(*jwt.SigningMethodRSA); !ok {
	// 		return nil, fmt.Errorf("Error on parse token")
	// 	}

	// 	//k := os.Getenv("CERT")
	// 	k := "My secret"
	// 	log.Println(k)
	// 	return []byte(k), nil
	// })

	// if claims, ok := jtoken.Claims.(jwt.MapClaims); ok && jtoken.Valid {
	// 	fmt.Println(claims)
	// } else {
	// 	fmt.Println(err)
	// }

}

// func TestMain(m *testing.M) {
// 	retriveToken()
// 	ret := m.Run()
// 	os.Exit(ret)
// }

func TestSimpleRoutes(t *testing.T) {

	var err error
	router, err = NewHandler()
	if err != nil {
		t.Fatalf("startup router give error:%s\n", err)
	}

	simpleApiTest := []testItem{
		{"welcome", testWelcome},
		{"infoNotLogged", testInfoNotLogged},
		{"restrictedNotLogged", testRestrictedNotLogged},
		// {"infoLogged", testInfoLogged},
		// {"restrictedLogged", testRestrictedLogged},
	}

	for _, ct := range simpleApiTest {
		t.Run(ct.Name, ct.TextFunc)
	}

}

var welcomeMSG string = `{"Message":"Hi, welcome to locale-mgmt, don't know who you are so go to login"}`

func testWelcome(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, welcomeMSG, w.Body.String())

	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/welcome", nil)
	router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, welcomeMSG, w.Body.String())
}

func testInfoNotLogged(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/info", nil)
	router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"version":"0.0.1"}`, w.Body.String())
}

func testInfoLogged(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/info", nil)
	router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}

func testRestrictedNotLogged(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/api/v1/restricted", nil)
	router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
}

func testRestrictedLogged(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/api/v1/restricted", nil)
	router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "You are in the restricted area")
}
