package test

import (
	"bytes"
	"encoding/json"
	"firstbeegoapi/internal/shared"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	_ "firstbeegoapi/routers"

	beego "github.com/beego/beego/v2/server/web"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

// TestGet is a sample to run an endpoint test
func TestGet(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v1/ordering/1", nil)
	r.Header.Set("Authorization", "Bearer "+testToken(t))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
	})
}

func TestGetWithInvalidID(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v1/ordering/abc", nil)
	r.Header.Set("Authorization", "Bearer "+testToken(t))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var body map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &body)

	Convey("Subject: Test Ordering Error Response\n", t, func() {
		Convey("Status Code Should Be 400", func() {
			So(w.Code, ShouldEqual, 400)
		})
		Convey("The Result Should Be JSON", func() {
			So(err, ShouldBeNil)
		})
		Convey("The Error Should Use The Shared Error Handler", func() {
			errorBody := body["error"].(map[string]any)
			So(errorBody["code"], ShouldEqual, "invalid_object_id")
		})
	})
}

func TestLogin(t *testing.T) {
	body := bytes.NewBufferString(`{"email":"admin@example.com","password":"password"}`)
	r, _ := http.NewRequest("POST", "/v1/auth/login", body)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Auth Login\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Be JSON", func() {
			So(err, ShouldBeNil)
		})
		Convey("The Result Should Contain Access Token", func() {
			data := response["data"].(map[string]any)
			So(data["access_token"], ShouldNotBeEmpty)
			So(data["token_type"], ShouldEqual, "Bearer")
		})
	})
}

func TestOrderingWithoutToken(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v1/ordering/1", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var body map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &body)

	Convey("Subject: Test JWT Middleware\n", t, func() {
		Convey("Status Code Should Be 401", func() {
			So(w.Code, ShouldEqual, 401)
		})
		Convey("The Result Should Be JSON", func() {
			So(err, ShouldBeNil)
		})
		Convey("The Error Should Be Missing Authorization", func() {
			errorBody := body["error"].(map[string]any)
			So(errorBody["code"], ShouldEqual, "missing_authorization")
		})
	})
}

func testToken(t *testing.T) string {
	t.Helper()

	token, err := shared.GenerateJWT(shared.JWTUser{
		ID:    1,
		Name:  "Admin",
		Email: "admin@example.com",
	}, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}

	return token
}
