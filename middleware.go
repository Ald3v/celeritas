package celeritas

import "net/http"

func (c *Celeritas) SessionLoad(next http.Handler) http.Handler{
	c.Logger.PrintInfo("SessionLoad called",nil)
	return c.Session.LoadAndSave(next)
}