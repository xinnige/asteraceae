package auth0api

import (
	"errors"
	// "fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock "github.com/xinnige/asteraceae/calendula/mock"
	utils "github.com/xinnige/asteraceae/calendula/utils"
)

func fakeUser() string {
	return `{"nickname":"yamada_taro","groups":[],"dn":"uid=yamada_taro,ou=People,dc=asteraceae,dc=local","organizationUnits":"uid=yamada_taro,ou=People,dc=asteraceae,dc=local","updated_at":"2018-10-01T00:02:03.091Z","name":"","picture":"","user_id":"ad|auth0-ldap01|yamada_taro","identities":[{"user_id":"auth0-ldap01|yamada_taro","provider":"ad","connection":"auth0-ldap01","isSocial":false}],"created_at":"2018-05-29T09:17:12.941Z","user_metadata":{"surname":"yamada","givenname":"taro"},"last_login":"2018-11-01T00:00:00.090Z","last_ip":"000.000.00.0","logins_count":8, "app_metadata":{"lambda_authorizer":true,"apps":["app1","app2"]}}`
}

func TestParseUser(t *testing.T) {
	client := &Auth0Client{
		SerialAPI: &utils.JSONAPI{},
	}
	user := &User{}
	err := client.ParseUser([]byte(fakeUser()), user)
	assert.NotNil(t, user)
	assert.Nil(t, err)
	assert.Equal(t, "yamada_taro", user.Nickname)
	appmeta, ok := user.AppMeta.(*AuthAppMeta)
	assert.Equal(t, true, appmeta.LambdaAuthorizer)
	assert.Equal(t, true, ok)
	usermeta, ok := user.UserMeta.(*SimpleUserMeta)
	assert.Equal(t, "taro", usermeta.Givenname)
	assert.Equal(t, true, ok)

	// mock unmarshal err
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockSiface := mock.NewMockSerialInterface(mockCtrl)

	// unmarshal user error
	mockSiface.EXPECT().Unmarshal(gomock.Any(), gomock.Any()).Return(
		errors.New("FakeJSONError")).Times(1)

	// unmarshal user ok, appmeta error
	mockSiface.EXPECT().Unmarshal(gomock.Any(), gomock.Any()).Return(
		nil).Times(1)
	mockSiface.EXPECT().Unmarshal(gomock.Any(), gomock.Any()).Return(
		errors.New("FakeJSONError")).Times(1)

	// unmarshal user,appmeta ok; usermeta error
	mockSiface.EXPECT().Unmarshal(gomock.Any(), gomock.Any()).Return(
		nil).Times(2)
	mockSiface.EXPECT().Unmarshal(gomock.Any(), gomock.Any()).Return(
		errors.New("FakeJSONError")).Times(1)

	client.SerialAPI = mockSiface
	err = client.ParseUser([]byte(fakeUser()), user)
	assert.NotNil(t, err)

	err = client.ParseUser([]byte(fakeUser()), user)
	assert.NotNil(t, err)

	err = client.ParseUser([]byte(fakeUser()), user)
	assert.NotNil(t, err)
}

// func TestParseUsers(t *testing.T) {
// 	fakeUsers := ``
//
// 	client := &Auth0Client{
// 		SerialAPI: &utils.JSONAPI{},
// 	}
// 	users := make([]User, 0)
// 	err := client.ParseUsers([]byte(fakeUsers), &users)
// 	fmt.Println(err)
//
// }
