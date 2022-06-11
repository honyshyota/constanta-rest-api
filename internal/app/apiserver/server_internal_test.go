package apiserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/honyshyota/constanta-rest-api/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestServer_HandleUsersCreate(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users", nil)
	srv := newServer(teststore.New())
	srv.ServeHTTP(rec, req)
	assert.Equal(t, rec.Code, http.StatusOK)
}
