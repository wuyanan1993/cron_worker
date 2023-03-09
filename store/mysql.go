package store

import "cron_worker/tasks"

type MysqlStore struct {
	ConnectionString string
}

func (s *MysqlStore) AddTask(tasks []tasks.Task) error {
	if tasks != nil {

	}
	return nil
}
