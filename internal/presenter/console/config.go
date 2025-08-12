package console

import (
	"rumi-go/internal/config"

	"gopkg.in/ukautz/clif.v1"
)

// StartConfigShow shows the current configuration
func (c *Console) StartConfigShow() *clif.Command {
	return clif.NewCommand("config:show", "show current configuration", func(o *clif.Command, in clif.Input, out clif.Output) error {
		conf := config.Get()

		out.Printf("Current Configuration:\n")
		out.Printf("====================\n")
		out.Printf("App:\n")
		out.Printf("  Name: %s\n", conf.App.Name)
		out.Printf("  Environment: %s\n", conf.App.Env)
		out.Printf("  API Prefix: %s\n", conf.App.ApiPrefix)
		out.Printf("\n")
		out.Printf("Server:\n")
		out.Printf("  Port: %d\n", conf.Server.Port)
		out.Printf("\n")
		out.Printf("Database:\n")
		out.Printf("  Driver: %s\n", conf.Database.Driver)
		out.Printf("  Name: %s\n", conf.Database.Name)
		out.Printf("  Host: %s\n", conf.Database.Host)
		out.Printf("  Port: %d\n", conf.Database.Port)
		out.Printf("  User: %s\n", conf.Database.User)
		out.Printf("  Max Open: %d\n", conf.Database.MaxOpen)
		out.Printf("  Max Idle: %d\n", conf.Database.MaxIdle)
		out.Printf("\n")
		out.Printf("Redis:\n")
		out.Printf("  Address: %s\n", conf.Redis.Address)
		out.Printf("  Port: %s\n", conf.Redis.Port)
		out.Printf("  DB: %d\n", conf.Redis.DB)
		out.Printf("\n")
		out.Printf("JWT:\n")
		out.Printf("  Secret: %s\n", conf.JWT.Secret)
		out.Printf("  Expiration: %d hours\n", conf.JWT.Expiration)

		return nil
	})
}
