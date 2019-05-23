package auth0api

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/xinnige/asteraceae/calendula/mock"
	utils "github.com/xinnige/asteraceae/calendula/utils"
)

func fakeResponse(content []byte) *http.Response {
	return &http.Response{
		Body:       utils.ReadCloser{Reader: bytes.NewBuffer(content)},
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
	}
}

func fakeClient() *Auth0Client {
	return &Auth0Client{
		Endpoint: &Auth0Endpoint{
			URL:        "fake-url",
			Provider:   "fake-provider",
			Connection: "fake-conn",
		},
		token:     "fake-token",
		SerialAPI: &utils.JSONAPI{},
	}
}

func TestGetUserByName(t *testing.T) {
	api := fakeClient()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClientiface := mock.NewMockAsterClient(mockCtrl)

	mockClientiface.EXPECT().Do(gomock.Any()).Return(
		fakeResponse([]byte(fakeUser())), nil).Times(1)

	api.httpClient = mockClientiface
	user, err := api.GetUserByName("yamada_taro")
	assert.Nil(t, err)
	assert.Equal(t, "yamada_taro", user.Nickname)
	assert.NotEqual(t, "", user.RawAppMeta)
	assert.NotNil(t, user.AppMeta)
}
