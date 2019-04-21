package read

import (
	github.com/SaraTrawnik/for-chan/color"
	"fmt"
	"os"
	"strconv"
	"regexp"
	"strings"
	"math"
)

var(
	stuffReplace = []*regexp.Regexp{regexp.MustCompile("<br>"),
					regexp.MustCompile("<a href=\"#p[0-9]+\" class=\"quotelink\">"),
					regexp.MustCompile("</a>"),
					regexp.MustCompile("&gt;"),
					regexp.MustCompile("<span class=\"quote\">"),
					regexp.MustCompile("</span>"),
					regexp.MustCompile("&#039;"),
					regexp.MustCompile("<s>"),
					regexp.MustCompile("</s>"),
					regexp.MustCompile("&quot;"),
					regexp.MustCompile("<wbr>"),
					}
	postColorizer = regexp.MustCompile(">>[0-9]+")
	stuffEndup = []string{"\n", "", "", ">", "", "", "'", "-", "-", "\"", ""}
)

func ReadThread(b string, t int64) {
	threadPosts, err := api.GetThread(b, t)
	if err != nil {
		fmt.Println("could not read thread", t, "from board", b)
		return
	}

	for _, x := range threadPosts.Posts {
		readPost(x)
	}
}

func ReadCatalog(b string) {
	entireCatalog, err := api.GetCatalog(b)
	if err != nil {
		fmt.Println("could not get catalog from", b, "board")
		return
	}
	for _, catalogPage := range entireCatalog {
		for _, oneThread := range catalogPage.Threads {
			fmt.Printf("R:%v I:%v ", oneThread.Replies(), oneThread.Images())
			readPost(oneThread.OP)
		}
	}

}

func readPost(p *api.Post) {
	boxPost := func(amount int) string { return strings.Repeat("-", amount) }
	var file string
	if p.File != nil { file = p.File.String()}
	spacers := boxPost(len(p.Name)) + boxPost(len(p.Time.String())) + boxPost(int(math.Log10(float64(p.Id)))) + "---"
	var additional string
	if len(file) > 0 {
		additional = boxPost(len(file)-len(spacers))
	} else { additional = "" }
	fmt.Printf("%v %v %v\n%v%v\n", p.Name, p.Time, color.Get(p.Id), file, parseComment(p.Comment))
	fmt.Printf("%v\n", spacers+additional)
}

func parseComment(comm string) string {
	var newComm string
	for y, x := range stuffReplace {
		comm = x.ReplaceAllString(comm, stuffEndup[y])
	}

	lenOfLine := 0
	for _, char := range comm {
		lenOfLine += 1
		if char == '\n' { lenOfLine = 0 }
		if lenOfLine > 80 {
			if char == ' ' {
				newComm += "\n"
				lenOfLine = 0
			}
		}
		newComm += string(char)
	}

	newComm = postColorizer.ReplaceAllStringFunc(newComm, func(id string) string { val, _ := strconv.ParseInt(id[2:], 10, 64); return ">>" + color.Get(val) })
	return regexp.MustCompile("\n ").ReplaceAllString(newComm, "\n")
}