package toolcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/timex"
)

// DateCmd command
var DateCmd = &gcli.Command{
	Name:    "date",
	Aliases: []string{"dt", "time"},
	Desc:    "display or parse the date and time",
	Subs: []*gcli.Command{
		DatePrintCmd,
	},
}

var dateParseOpts = struct {
	Format string `flag:"desc=the date format string;default=Y-m-d H:i:s;shorts=f"`
	Inline bool   `flag:"desc=inline output;shorts=i"`
}{}

// DatePrintCmd command
var DatePrintCmd = &gcli.Command{
	Name: "print",
	Desc: "print the current date and time",
	Config: func(c *gcli.Command) {
		c.MustFromStruct(&dateParseOpts)
	},
	Func: func(c *gcli.Command, _ []string) error {
		tx := timex.Now()

		s := tx.DateFormat(dateParseOpts.Format)
		if dateParseOpts.Inline {
			fmt.Print(s, " ", tx.Unix())
		} else {
			fmt.Println(s, tx.Unix())
		}
		return nil
	},
}

// DateFormatCmd command
