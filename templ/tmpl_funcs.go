package templ

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
)

func TimeAgo(t time.Time) string {
	d := time.Since(t)

	if d < time.Minute {
		seconds := int(d.Seconds())
		if seconds == 0 {
			return "just now"
		}
		return fmt.Sprintf("%d second(s) ago", seconds)
	} else if d < time.Hour {
		return fmt.Sprintf("%d minute(s) ago", int(d.Minutes()))
	} else if d < 24*time.Hour {
		return fmt.Sprintf("%d hour(s) ago", int(d.Hours()))
	} else if d < 30*24*time.Hour {
		return fmt.Sprintf("%d day(s) ago", int(d.Hours()/24))
	} else if d < 12*30*24*time.Hour {
		return fmt.Sprintf("%d month(s) ago", int(d.Hours()/(24*30)))
	}
	return fmt.Sprintf("%d year(s) ago", int(d.Hours()/(24*365)))
}

func DerefFloat64(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

func FormatFloat(f float64) string {
	//return strconv.FormatFloat(f, 'f', -1, 64)
	return humanize.FormatFloat("", f)
}
