package color

import "github.com/logrusorgru/aurora"

var (
	fgcolors = []aurora.Color{aurora.BlackFg, aurora.RedFg, aurora.GreenFg, aurora.BrownFg, aurora.BlueFg, aurora.MagentaFg, aurora.CyanFg, aurora.GrayFg}
	bgcolors = []aurora.Color{aurora.BlackBg, aurora.RedBg, aurora.GreenBg, aurora.BrownBg, aurora.BlueBg, aurora.MagentaBg, aurora.CyanBg, aurora.GrayBg}
)

func Get(id int64) string {
	var tempbgcolors []aurora.Color
	fgVal := hash(id)%8
	bgVal := hash(id)%7
	if fgVal == 0 { 
		tempbgcolors = bgcolors[(fgVal+1):] 
	}
	if fgVal == int64((len(bgcolors) - 1)) { 
		tempbgcolors = bgcolors[:fgVal]
	} else {
		tempbgcolors = append(bgcolors[(fgVal+1):], bgcolors[:fgVal]...)
	}
	return aurora.Sprintf(aurora.Colorize(id, fgcolors[fgVal] | tempbgcolors[bgVal]))
}

func hash(id int64) int64 {
	return (((((id + 1) ^ (id >> 3)) * 7) ^ (id >> 3)) * 11)
}