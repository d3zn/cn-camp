package service

import "m02/quark"

func (c *Controller) Register(engine *quark.Engine) {
	engine.GET("/header", c.Header)
	engine.POST("/env", c.Env)
	engine.GET("/health", c.Health)
	engine.GET("/random", c.Random)
}
