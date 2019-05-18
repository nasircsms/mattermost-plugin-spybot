package main

import (
	"encoding/json"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestSpy(t *testing.T) {

	user := &model.User{
		Email:    "-@-.-",
		Nickname: "TestUser",
		Password: model.NewId(),
		Username: "testuser",
		Roles:    model.SYSTEM_USER_ROLE_ID,
		Locale:   "en",
	}

	targetWatches := []TargetWatch{
		{
			Target:  "test",
			Status:  "online",
			Watcher: "testuser",
		},
	}

	targetStatus := &model.Status{
		UserId: model.NewId(),
		Status: "online",
		Manual: false,
	}

	stringTargets, _ := json.Marshal(targetWatches)

	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
		api.On("GetUserStatus", mock.AnythingOfType("string")).Return(targetStatus, nil)
		api.On("KVGet", mock.Anything).Return(stringTargets, nil)
		api.On("KVSet", mock.Anything, mock.Anything).Return(nil)

		return api
	}

	t.Run("if update list happens", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		p.spy(targetWatches[0].Target, targetWatches[0].Watcher)

	})

}

func TestUnSpy(t *testing.T) {

	targetWatches := []TargetWatch{
		{
			Target:  "test",
			Status:  "online",
			Watcher: "testuser",
		},
	}

	stringTargets, _ := json.Marshal(targetWatches)

	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("KVGet", mock.Anything).Return(stringTargets, nil)
		api.On("KVSet", mock.Anything, mock.Anything).Return(nil)

		return api
	}

	t.Run("if update list happens", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		p.unSpy(targetWatches[0].Target, targetWatches[0].Watcher)

	})

}

func TestSpyList(t *testing.T) {

	targetWatches := []TargetWatch{
		{
			Target:  "test",
			Status:  "online",
			Watcher: "testuser",
		},
	}

	stringTargets, _ := json.Marshal(targetWatches)

	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("KVGet", mock.Anything).Return(stringTargets, nil)

		return api
	}

	t.Run("if update list happens", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		p.list(targetWatches[0].Watcher)

	})

}
