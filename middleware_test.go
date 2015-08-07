package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-martini/martini"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTokenValidator struct {
	mock.Mock
}

func (m *MockTokenValidator) CheckToken(token string) bool {

	args := m.Called(token)
	return args.Bool(0)

}

type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockCache) Set(key, val string) error {
	args := m.Called(key, val)
	return args.Error(0)
}

func TestKeystoneCachedToken(t *testing.T) {
	token := "token"

	// Set up mock token validator - no auth attempts are expected in this test
	testTokenValidator := new(MockTokenValidator)

	// Set up mock cache - the middleware should try to get the token
	testCache := new(MockCache)
	testCache.On("Get", token).Return(token, nil)

	// Dummy martini
	m := martini.New()
	m.Use(Keystone(testTokenValidator, testCache))
	m.Use(func(res http.ResponseWriter, req *http.Request, token Token) {
		res.Write([]byte(token))
	})

	// Send request to martini
	recorder := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "foo", nil)
	r.Header.Add("X-Auth-Token", token)
	m.ServeHTTP(recorder, r)

	// Assert mock expectations
	testTokenValidator.AssertExpectations(t)
	testCache.AssertExpectations(t)

	// Assert http response expectations
	assert.Equal(t, recorder.Code, 200, "Expected: 200")
	assert.Equal(t, recorder.Body.String(), token, "Expected: "+token)
}

func TestKeystoneValidUncachedToken(t *testing.T) {
	token := "token"

	// Set up mock token validator - one auth attempt is expected in this test
	testTokenValidator := new(MockTokenValidator)
	testTokenValidator.On("CheckToken", token).Return(true)

	// Set up mock cache - the middleware should try to get and set the token
	testCache := new(MockCache)
	testCache.On("Get", token).Return("", errors.New("test error"))
	testCache.On("Set", token, "authorized").Return(nil)

	// Dummy martini
	m := martini.New()
	m.Use(Keystone(testTokenValidator, testCache))
	m.Use(func(res http.ResponseWriter, req *http.Request, token Token) {
		res.Write([]byte(token))
	})

	// Send request to martini
	recorder := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "foo", nil)
	r.Header.Add("X-Auth-Token", token)
	m.ServeHTTP(recorder, r)

	// Assert mock expectations
	testTokenValidator.AssertExpectations(t)
	testCache.AssertExpectations(t)

	// Assert http response expectations
	assert.Equal(t, recorder.Code, 200, "Expected: 200")
	assert.Equal(t, recorder.Body.String(), token, "Expected: "+token)
}

func TestKeystoneInvalidToken(t *testing.T) {
	token := "token"

	// Set up mock token validator - one auth attempt is expected in this test
	testTokenValidator := new(MockTokenValidator)
	testTokenValidator.On("CheckToken", token).Return(false)

	// Set up mock cache - the middleware should try to get the token
	testCache := new(MockCache)
	testCache.On("Get", token).Return("", errors.New("test error"))

	// Dummy martini
	m := martini.New()
	m.Use(Keystone(testTokenValidator, testCache))
	m.Use(func(res http.ResponseWriter, req *http.Request, token Token) {
		res.Write([]byte(token))
	})

	// Send request to martini
	recorder := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "foo", nil)
	r.Header.Add("X-Auth-Token", token)
	m.ServeHTTP(recorder, r)

	// Assert mock expectations
	testTokenValidator.AssertExpectations(t)
	testCache.AssertExpectations(t)

	// Assert http response expectations
	assert.Equal(t, recorder.Code, 401, "Expected: 401")
	assert.Equal(t, recorder.Body.String(), "Not Authorized\n", "Expected: 'Not Authorized\n'")
}

func TestKeystoneMissingToken(t *testing.T) {
	testTokenValidator := new(MockTokenValidator)
	testCache := new(MockCache)

	// Dummy martini
	m := martini.New()
	m.Use(Keystone(testTokenValidator, testCache))
	m.Use(func(res http.ResponseWriter, req *http.Request, token Token) {
		res.Write([]byte(token))
	})

	// Send request to martini
	recorder := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "foo", nil)
	m.ServeHTTP(recorder, r)

	// Assert mock expectations
	testTokenValidator.AssertExpectations(t)
	testCache.AssertExpectations(t)

	// Assert http response expectations
	assert.Equal(t, recorder.Code, 401, "Expected: 401")
	assert.Equal(t, recorder.Body.String(), "Not Authorized\n", "Expected: 'Not Authorized\n'")
}
