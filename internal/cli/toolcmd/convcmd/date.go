package convcmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gookit/color/colorp"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/timex"
	"github.com/inhere/kite-go/internal/apputil"
)

// NewTime2dateCmd instance
func NewTime2dateCmd() *gcli.Command {
	var timeRegex = regexp.MustCompile(`1\d{9}`)

	return &gcli.Command{
		Name:    "ts2date",
		Aliases: []string{"ts", "t2d", "t2date"},
		Desc:    "Quick convert all timestamp number to datetime",
		Config: func(c *gcli.Command) {
			c.AddArg("input", "want parsed input contents", true, true)
		},
		Func: func(c *gcli.Command, _ []string) (err error) {
			var txt string
			ss := c.Arg("input").Strings()
			if len(ss) > 1 {
				txt = strings.Join(ss, " ")
			} else {
				txt, err = apputil.ReadSource(ss[0])
				if err != nil {
					return err
				}
			}

			c.Infoln("Input Contents:")
			fmt.Println(txt + "\n")

			times := timeRegex.FindAllString(txt, -1)
			if len(times) == 0 {
				return errorx.Raw("not found any timestamps")
			}

			mp := make(map[string]string, len(times))
			for _, timeVal := range times {
				mp[timeVal] = timex.FormatUnix(mathutil.SafeInt64(timeVal))
			}

			show.AList("Matched timestamps", mp, func(opts *show.ListOption) {
				opts.SepChar = "  =>  "
			})
			return nil
		},
	}
}

// NewDate2tsCmd instance
func NewDate2tsCmd() *gcli.Command {
	var oneDay bool

	return &gcli.Command{
		Name:    "date2ts",
		Aliases: []string{"d2ts", "d2t"},
		Desc:    "Quick convert datetime to unix timestamp",
		Config: func(c *gcli.Command) {
			c.BoolOpt2(&oneDay, "one-day,od", "parse input date, get the timestamp of the day start and end")

			c.AddArg("input", "want parsed input contents", true, true)
		},
		Func: func(c *gcli.Command, _ []string) (err error) {
			var txt string
			ss := c.Arg("input").Strings()
			if len(ss) > 1 {
				txt = strings.Join(ss, " ")
			} else {
				txt, err = apputil.ReadSource(ss[0])
				if err != nil {
					return err
				}

				ss = strings.Split(txt, "\n")
			}

			c.Infoln("Input Contents:")
			fmt.Println(txt + "\n")

			c.Infoln("Parsed Results:")
			for _, s := range ss {
				tt, err := timex.FromDate(s)
				if err != nil {
					colorp.Warnf("parse date %q error: %s\n", s, err.Error())
					continue
				}

				if oneDay {
					fmt.Printf("%s => %d ~ %d\n", s, tt.DayStart().Unix(), tt.DayEnd().Unix())
				} else {
					fmt.Printf("%s => %d\n", s, tt.Unix())
				}
			}

			return nil
		},
	}
}
