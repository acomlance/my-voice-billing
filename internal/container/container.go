package container

import (
	"context"
	"sync"

	"my-voice-billing/internal/config"
	"my-voice-billing/internal/repository/db"
	"my-voice-billing/internal/repository/repo"
	"my-voice-billing/internal/services/logic"
)

var (
	c        *Container
	initOnce sync.Once
)

type Container struct {
	Config *config.Config

	DB          *db.Manager
	AccountRepo *repo.AccountRepo
	TaskRepo    *repo.TaskRepo

	AccountLogic *logic.AccountLogic
	TaskLogic    *logic.TaskLogic
}

func Init(ctx context.Context, cfg *config.Config) error {
	var initErr error
	initOnce.Do(func() {
		c = &Container{Config: cfg}
		c.DB = &db.Manager{}
		if err := c.DB.Connect(ctx, cfg); err != nil {
			initErr = err
			return
		}
		c.AccountRepo = repo.NewAccountRepo(c.DB)
		c.TaskRepo = repo.NewTaskRepo(c.DB)
		c.AccountLogic = logic.NewAccountLogic(c.AccountRepo)
		c.TaskLogic = logic.NewTaskLogic(c.TaskRepo, c.AccountRepo)
	})
	return initErr
}

func Get() *Container {
	if c == nil {
		panic("container not initialized: call Init first")
	}
	return c
}

func Shutdown() {
	if c != nil && c.DB != nil {
		c.DB.Close()
	}
}
