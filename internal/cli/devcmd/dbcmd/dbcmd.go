package dbcmd

import "github.com/gookit/gcli/v3"

// NewDBCmd the db command
func NewDBCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "db",
		Desc:    "Provide some useful database commands for mysql,pgsql,redis",
		Aliases: []string{"database"},
		Config: func(c *gcli.Command) {
			c.AddArg("command", "The command to execute", true)
		},
		// Func: func(c *gcli.Command, args []string) error {
		// 	return c.RunSubCmd()
		// },
		Subs: []*gcli.Command{
			// TODO
			// NewMysqlCmd(),
			// NewPgsqlCmd(),
			// NewRedisCmd(),
		},
	}
}
