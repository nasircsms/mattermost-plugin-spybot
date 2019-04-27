package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/mattermost/mattermost-server/model"
)

const TriggerHostName = "__SPYHOST__"
const WatchedTargets = "__WATCHED__"

type TargetWatch struct {
	Target  string
	Status  string
	Watcher string
}

func (p *Plugin) spy(target string, watcher string) {

	targetUser, _ := p.API.GetUserByUsername(target)
	status, _ := p.API.GetUserStatus(targetUser.Id)

	bytes, _ := p.API.KVGet(WatchedTargets)

	var targets []TargetWatch
	json.Unmarshal(bytes, &targets)
	targetWatch := TargetWatch{
		Target:  target,
		Status:  status.Status,
		Watcher: watcher,
	}
	targets = append(targets, targetWatch)
	ro, _ := json.Marshal(targets)

	p.API.KVSet(WatchedTargets, ro)

}

func (p *Plugin) trigger() {

	bytes, _ := p.API.KVGet(WatchedTargets)
	var targets []TargetWatch
	json.Unmarshal(bytes, &targets)

	var updatedTargets []TargetWatch
	for _, targetWatch := range targets {
		targetUser, _ := p.API.GetUserByUsername(targetWatch.Target)
		status, _ := p.API.GetUserStatus(targetUser.Id)
		if status.Status != targetWatch.Status {

			watcherUser, _ := p.API.GetUserByUsername(targetWatch.Watcher)

			channel, cErr := p.API.GetDirectChannel(p.spyUserId, watcherUser.Id)
			if cErr != nil {
				p.API.LogError("failed to create channel " + cErr.Error())
				continue
			}

			post := model.Post{
				ChannelId: channel.Id,
				UserId:    p.spyUserId,
				Message:   "user @" + targetWatch.Target + " is now " + status.Status + ".",
			}
			p.API.CreatePost(&post)

			targetWatch.Status = status.Status
		}
		updatedTargets = append(updatedTargets, targetWatch)
	}
	ro, _ := json.Marshal(updatedTargets)

	p.API.KVSet(WatchedTargets, ro)
}

func (p *Plugin) Run() {

	hostname, _ := os.Hostname()
	bytes, bErr := p.API.KVGet(TriggerHostName)
	if bErr != nil {
		p.API.LogError("failed KVGet %s", bErr)
		return
	}
	if string(bytes) != "" && string(bytes) != hostname {
		return
	}
	p.API.KVSet(TriggerHostName, []byte(hostname))

	if !p.running {
		p.running = true
		p.runner()
	}

}

func (p *Plugin) Stop() {
	p.API.KVSet(TriggerHostName, []byte(""))
	p.running = false
}

func (p *Plugin) runner() {

	go func() {
		<-time.NewTimer(time.Second).C
		p.trigger()
		if !p.running {
			return
		}
		p.runner()
	}()
}
