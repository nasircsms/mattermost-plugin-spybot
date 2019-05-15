package main

import (
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/pkg/errors"
	"strings"
)

const CommandTrigger = "spy"
const UnCommandTrigger = "unspy"

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

	if err := p.API.RegisterCommand(&model.Command{
		TeamId:           teamId,
		Trigger:          UnCommandTrigger,
		Username:         botName,
		AutoComplete:     true,
		AutoCompleteHint: "[@someone]",
		AutoCompleteDesc: "UnListen for another users presence",
		DisplayName:      "Spy Plugin Command",
		Description:      "A command used to stop notifying a user of another users presence",
	}); err != nil {
		return errors.Wrap(err, "failed to register command")
	}

	return nil
}

func (p *Plugin) unregisterCommand(teamId string) error {
	p.API.UnregisterCommand(teamId, UnCommandTrigger)
	return p.API.UnregisterCommand(teamId, CommandTrigger)
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {

	_, uErr := p.API.GetUser(args.UserId)
	if uErr != nil {
		return &model.CommandResponse{}, uErr
	}

	if strings.HasPrefix(args.Command, "/"+CommandTrigger+" @") {
		target := strings.Split(args.Command, "/"+CommandTrigger+" @")[1]

		user, uErr := p.API.GetUserByUsername(target)
		if uErr != nil {
			post := model.Post{
				ChannelId: args.ChannelId,
				UserId:    p.spyUserId,
				Message:   "invalid user @" + target,
			}
			p.API.SendEphemeralPost(user.Id, &post)
			return &model.CommandResponse{}, uErr
		} else {
			p.spy(target, user.Username)
		}

		post := model.Post{
			ChannelId: args.ChannelId,
			UserId:    p.spyUserId,
			Message:   "spying on @" + target,
		}
		p.API.SendEphemeralPost(user.Id, &post)
	} else if strings.HasPrefix(args.Command, "/"+UnCommandTrigger+" @") {
		target := strings.Split(args.Command, "/"+UnCommandTrigger+" @")[1]

		user, uErr := p.API.GetUserByUsername(target)
		if uErr != nil {
			post := model.Post{
				ChannelId: args.ChannelId,
				UserId:    p.spyUserId,
				Message:   "invalid user @" + target,
			}
			p.API.SendEphemeralPost(user.Id, &post)
			return &model.CommandResponse{}, uErr
		} else {
			p.unSpy(target, user.Username)
		}

		post := model.Post{
			ChannelId: args.ChannelId,
			UserId:    p.spyUserId,
			Message:   "unspying on @" + target,
		}
		p.API.SendEphemeralPost(user.Id, &post)
	}

	return &model.CommandResponse{}, nil

}
