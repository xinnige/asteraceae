package slackapi

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/xinnige/asteraceae/calendula/mock"
	"github.com/xinnige/asteraceae/calendula/utils"
)

func fakeAuditLogsNext() []byte {
	jsonBytes, _ := utils.ReadFile("../test/slack/auditlogs_next.json")
	return jsonBytes
}

func fakeAuditLogs() []byte {
	jsonBytes, _ := utils.ReadFile("../test/slack/auditlogs.json")
	return jsonBytes
}

func fakeAuditActions() []byte {
	jsonBytes, _ := utils.ReadFile("../test/slack/auditactions.json")
	return jsonBytes
}

func fakeAuditSchemas() []byte {
	jsonBytes, _ := utils.ReadFile("../test/slack/auditschemas.json")
	return jsonBytes
}

func fakeResponse(content []byte) *http.Response {
	return &http.Response{
		Body:       utils.ReadCloser{Reader: bytes.NewBuffer(content)},
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
	}
}

func fakeClient() *Client {
	return &Client{
		token:     "",
		unmarshal: json.Unmarshal,
		marshal:   json.Marshal,
		debug:     false,
		log:       log.New(os.Stderr, "debug", log.LstdFlags|log.Lshortfile),
	}
}

func TestGetAuditLogsPaginated(t *testing.T) {
	client := fakeClient()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClientiface := mock.NewMockAsterClient(mockCtrl)
	mockClientiface.EXPECT().Do(gomock.Any()).Return(
		fakeResponse([]byte(fakeAuditLogsNext())), nil).Times(1)
	mockClientiface.EXPECT().Do(gomock.Any()).Return(
		fakeResponse([]byte(fakeAuditLogs())), nil).Times(1)
	client.client = mockClientiface

	p := client.GetAuditLogsPaginated(
		AuditLogsOptionLimit(10), AuditLogsOptionOldest(1559626515))
	p, err := p.Next(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 2, len(p.Entries))

	p, err = p.Next(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 1, len(p.Entries))
}

func TestListAuditLogs(t *testing.T) {
	client := fakeClient()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClientiface := mock.NewMockAsterClient(mockCtrl)
	mockClientiface.EXPECT().Do(gomock.Any()).Return(
		fakeResponse([]byte(fakeAuditLogsNext())), nil).Times(1)
	mockClientiface.EXPECT().Do(gomock.Any()).Return(
		fakeResponse([]byte(fakeAuditLogs())), nil).Times(1)
	client.client = mockClientiface

	result, err := client.ListAuditLogs(10, 1559636515, 1559626515,
		"fake-action", "fake-actor", "fake-entity")
	assert.Nil(t, err)
	assert.Equal(t, 3, len(result))
}

func TestListAuditActions(t *testing.T) {
	client := fakeClient()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClientiface := mock.NewMockAsterClient(mockCtrl)
	mockClientiface.EXPECT().Do(gomock.Any()).Return(
		fakeResponse([]byte(fakeAuditActions())), nil).Times(1)
	client.client = mockClientiface

	result, err := client.GetActions()
	assert.Nil(t, err)
	assert.Equal(t, 5, len(result["app"]))

}

func TestListAuditSchemas(t *testing.T) {
	client := fakeClient()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClientiface := mock.NewMockAsterClient(mockCtrl)
	mockClientiface.EXPECT().Do(gomock.Any()).Return(
		fakeResponse([]byte(fakeAuditSchemas())), nil).Times(1)
	client.client = mockClientiface

	result, err := client.GetSchemas()
	assert.Nil(t, err)
	assert.Equal(t, "string", result.Workspace.ID)
	assert.Equal(t, "array", result.App.Scopes)

}
