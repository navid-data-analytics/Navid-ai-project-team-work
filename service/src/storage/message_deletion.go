package storage

import (
	"context"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/postgres"

	"github.com/go-pg/pg"
)

// CheckExistingMessageIDs fetches messages from postgres
func CheckExistingMessageIDs(ctx context.Context, postgresClient postgres.Client, ids []int32) []int32 {
	logger := log.FromContext(ctx)
	logger.Info("Selecting all messages...")
	db, err := postgresClient.DB(ctx)

	err = db.Model(&Message{}).Column("id").Select(&ids)
	if err != nil {
		panic(err)
	}
	logger.Info("Messages retrieved from database:" + string(ids))
	return ids
}

// DeleteMessageByID deletes messages from postgres
func DeleteMessageByID(ctx context.Context, postgresClient postgres.Client, ids []int32) {
	logger := log.FromContext(ctx)
	logger.Debug("Deleting message of ids" + string(ids))

	deleteIds := pg.In(ids)
	db, err := postgresClient.DB(ctx)
	_, err = db.Model((*Message)(nil)).Where("id IN (?)", deleteIds).Delete()
	logger.Info("messages deleted")
	if err != nil {
		panic(err)
	}
	logger.Debug("Messages (by id):" + string(ids) + "deleted!")

}
