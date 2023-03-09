package store

import "cron_worker/tasks"

type RedisStore struct {
	ConnectionString string
}

func (r *RedisStore) AddTask(tasks []tasks.Task) error {
	if tasks != nil {

	}
	return nil
}
