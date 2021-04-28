package app

// Get implements gondola/kvs.Storage.Get
func (app *App) Get(key interface{}) interface{} {
	return app.kv.Get(key)
}

// Set implements gondola/kvs.Storage.Set
func (app *App) Set(key, value interface{}) {
	app.kv.Set(key, value)
}

// Get implements gondola/kvs.Storage.Get
func (c *Context) Get(key interface{}) interface{} {
	return c.kv.Get(key)
}

// Set implements gondola/kvs.Storage.Set
func (c *Context) Set(key, value interface{}) {
	c.kv.Set(key, value)
}
