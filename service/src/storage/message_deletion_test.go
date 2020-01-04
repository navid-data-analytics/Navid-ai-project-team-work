package storage_test

import (
	"context"
	"github.com/callstats-io/ai-decision/service/src/storage"
	"github.com/callstats-io/go-common/postgres"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDeleteMessagesByID(t *testing.T) {
	emptyIDs := []int32{}
	testTemplate := &storage.MessageTemplate{Type: "test_fa_create_tmpls_2", Version: 1, Template: "{.String \"val12\"}"}
	_, err := testPostgresDB.Model(testTemplate).Returning("*").Insert()
	require.Nil(t, err)

	createdMessages := []*storage.Message{
		{AppID: 123, TemplateID: testTemplate.ID, GeneratedAt: time.Now().Add(-5 * time.Minute), Data: []byte(`{"val1":"abc1"}`)},
		{AppID: 123, TemplateID: testTemplate.ID, GeneratedAt: time.Now().Add(-2 * time.Minute), Data: []byte(`{"val2":"abc2"}`)},
		{AppID: 123, TemplateID: testTemplate.ID, GeneratedAt: time.Now(), Data: []byte(`{"val3":"abc3"}`)},
	}
	_, err = testPostgresDB.Model(&createdMessages).Returning("*").Insert()
	require.Nil(t, err)
	checkMessages(t, testCtx, testPostgresClient, []int32{1, 2, 3})
	storage.DeleteMessageByID(testCtx, testPostgresClient, []int32{1, 3})
	checkMessages(t, testCtx, testPostgresClient, []int32{2})
	storage.DeleteMessageByID(testCtx, testPostgresClient, []int32{2})
	checkMessages(t, testCtx, testPostgresClient, emptyIDs)
}

func checkMessages(t *testing.T, ctx context.Context, postgresClient postgres.Client, ids []int32) {
	assert := require.New(t)
	allMessageIDs := storage.CheckExistingMessageIDs(ctx, postgresClient, []int32{})
	assert.Equal(allMessageIDs, ids)
}
