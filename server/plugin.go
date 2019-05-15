package main

import (
	"github.com/mattermost/mattermost-server/plugin"
	"io/ioutil"
)

type Plugin struct {
	plugin.MattermostPlugin

	spyUserId string

	running bool

	readFile func(path string) ([]byte, error)
}

func NewPlugin() *Plugin {
	return &Plugin{
		readFile: ioutil.ReadFile,
	}
}
