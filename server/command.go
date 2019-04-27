package main

import (
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/pkg/errors"
	"strings"
)

const CommandTrigger = "spy"

func (p *Plugin) registerCommand(teamId string) error {
	if err := p.API.RegisterCommand(&model.Command{
		TeamId:           teamId,
		Trigger:          CommandTrigger,
		Username:         botName,
		AutoComplete:     true,
		AutoCompleteHint: "[@someone]",
		AutoCompleteDesc: "Listen for another users presence",
		DisplayName:      "Spy Plugin Command",
		Description:      "A command used to notify a user of another users presence",
	}); err != nil {
		return errors.Wrap(err, "failed to register command")
	}

	return nil
}

func (p *Plugin) unregisterCommand(teamId string) error {
	return p.API.UnregisterCommand(teamId, CommandTrigger)
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {

	user, uErr := p.API.GetUser(args.UserId)
	if uErr != nil {
		return &model.CommandResponse{}, uErr
	}

	target := strings.Split(args.Command, "/"+CommandTrigger+" @")[1]

	p.spy(target, user.Username)

	post := model.Post{
		ChannelId: args.ChannelId,
		UserId:    p.spyUserId,
		Message:   "spying on @" + target,
	}
	p.API.SendEphemeralPost(user.Id, &post)
	return &model.CommandResponse{}, nil

}
