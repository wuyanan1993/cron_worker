package store

import "cron_worker/tasks"

type MongodbStore struct {
	ConnectionString string
}

func (m *MongodbStore) AddTask(tasks []tasks.Task) error {
	if tasks != nil {

	}
	return nil
}
