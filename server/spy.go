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

	if targetUser, err := p.API.GetUserByUsername(target); err != nil {
		p.API.LogError(err.Error())
	} else {
		if status, err := p.API.GetUserStatus(targetUser.Id); err != nil {
			p.API.LogError(err.Error())
		} else {
			if bytes, err := p.API.KVGet(WatchedTargets); err != nil {
				p.API.LogError(err.Error())
			} else {
				var targets []TargetWatch
				json.Unmarshal(bytes, &targets)
				targetWatch := TargetWatch{
					Target:  target,
					Status:  status.Status,
					Watcher: watcher,
				}
				targets = append(targets, targetWatch)
				if jsonTargets, err := json.Marshal(targets); err != nil {
					p.API.LogError(err.Error())
				} else {
					if err := p.API.KVSet(WatchedTargets, jsonTargets); err != nil {
						p.API.LogError(err.Error())
					}
				}
			}
		}
	}

}
func (p *Plugin) unSpy(target string, watcher string) {

	if bytes, err := p.API.KVGet(WatchedTargets); err != nil {
		p.API.LogError(err.Error())
	} else {
		var targets []TargetWatch
		var updateTargets []TargetWatch
		json.Unmarshal(bytes, &targets)

		for _, targetWatch := range targets {
			if targetWatch.Watcher != watcher && targetWatch.Target != target {
				updateTargets = append(updateTargets, targetWatch)
			}
		}
		if jsonTargets, err := json.Marshal(updateTargets); err != nil {
			p.API.LogError(err.Error())
		} else {
			if err := p.API.KVSet(WatchedTargets, jsonTargets); err != nil {
				p.API.LogError(err.Error())
			}
		}
	}
}

func (p *Plugin) trigger() {

	bytes, _ := p.API.KVGet(WatchedTargets)
	var targets []TargetWatch
	json.Unmarshal(bytes, &targets)

	var updatedTargets []TargetWatch
	for _, targetWatch := range targets {
		if targetUser, err := p.API.GetUserByUsername(targetWatch.Target); err != nil {
			p.API.LogError(err.Error())
		} else {
			status, _ := p.API.GetUserStatus(targetUser.Id)
			if status.Status != targetWatch.Status && status.Status != "away" {

				watcherUser, _ := p.API.GetUserByUsername(targetWatch.Watcher)

				channel, cErr := p.API.GetDirectChannel(p.spyUserId, watcherUser.Id)
				if cErr != nil {
					p.API.LogError("failed to create channel ")
					continue
				}

				prefix := "user"
				if targetUser.IsBot {
					prefix = "bot"
				}
				post := model.Post{
					ChannelId: channel.Id,
					UserId:    p.spyUserId,
					Message:   prefix + " @" + targetWatch.Target + " is now " + status.Status + ".",
				}
				p.API.CreatePost(&post)

				targetWatch.Status = status.Status
			}
			updatedTargets = append(updatedTargets, targetWatch)
		}

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
