package main

import "fmt"

func (p *Plugin) spy(username string) {

	user, _ := p.API.GetUserByUsername(username)

	p.API.LogInfo(fmt.Sprintf("%v", user.LastActivityAt))
}
