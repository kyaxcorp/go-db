package dbresolver

import (
	"time"

	"github.com/gookit/color"
	"github.com/kyaxcorp/go-helper/conv"
	"github.com/rs/zerolog"
)

func printConsumedTime(info func() *zerolog.Event, funcStartTime int) {
	funcEndTime := time.Now().Nanosecond()
	totalConsumedTime := funcEndTime - funcStartTime
	info().
		Int("consumed_time", totalConsumedTime).
		Msg(color.LightBlue.Render("consumed time -> Nano:" + conv.IntToStr(totalConsumedTime)))
}
