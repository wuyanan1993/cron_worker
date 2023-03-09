package store

import "cron_worker/tasks"

type RemoteStore interface {
	AddTask(tasks []tasks.Task) error
}
