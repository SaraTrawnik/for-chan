package main

import(
	"github.com/moshee/go-4chan-api/api"
	"github.com/logrusorgru/aurora"

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
									regexp.MustCompile("<wbr>")}
	stuffEndup = []string{"\n", "", "", ">", "", "", "'", "-", "-", "\"", ""}
	postColorizer = regexp.MustCompile(">>[0-9]+")
	fgcolors = []aurora.Color{aurora.BlackFg, aurora.RedFg, aurora.GreenFg, aurora.BrownFg, aurora.BlueFg, aurora.MagentaFg, aurora.CyanFg, aurora.GrayFg}
	bgcolors = []aurora.Color{aurora.BlackBg, aurora.RedBg, aurora.GreenBg, aurora.BrownBg, aurora.BlueBg, aurora.MagentaBg, aurora.CyanBg, aurora.GrayBg}
)

func main() {
	board, thread := ParseArgs()
	if board == "" && thread == 0 { fmt.Println("didnt provide nor thread nor board"); return }
	if board == "" { fmt.Println("how can i find a thread if no board is given"); return }
	if thread == 0 { ReadCatalog(board) } else { ReadThread(board, thread) }
}

func ParseArgs() (string, int64) {
	arguments := os.Args[1:]
	var chosenBoard string
	var chosenThread int64
	if x := contains("-b", arguments); x != -1 {
		if x < len(arguments) {
			if board, _ := api.LookupBoard(arguments[x+1]); board.Board == arguments[x+1] {
				chosenBoard = board.Board
			} else {
				fmt.Println("no such board")
				return "", 0
			}
		} else {
			fmt.Println("you didn't provide any board")
			return "", 0
		}
	} else {
		fmt.Println("no board given")
		return "", 0
	}

	if x := contains("-t", arguments); x != -1 {
		if x < len(arguments) {
			t, err := strconv.ParseInt(arguments[x+1], 10, 64)
			if err != nil { fmt.Println("invalid thread number"); return "", 0}
			chosenThread = t
		} else {
			fmt.Println("this isn't a number dummy")
			return "", 0
		}
	}
	return chosenBoard, chosenThread
}
func contains(a string, in []string) int {
	for place, arg := range in {
		if arg == a {
			return place
		}
	}
	return -1
}

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
	var file string
	if p.File != nil { file = p.File.String()}
	spacers := boxPost(len(p.Name)) + boxPost(len(p.Time.String())) + boxPost(int(math.Log10(float64(p.Id)))) + "---"
	var additional string
	if len(file) > 0 {
		additional = boxPost(len(file)-len(spacers))
	} else { additional = "" }
	fmt.Printf("%v %v %v\n%v%v\n", p.Name, p.Time, getColor(p.Id), file, parseComment(p.Comment))
	fmt.Printf("%v\n", spacers+additional)
}
func boxPost(amount int) string { return strings.Repeat("-", amount) }

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

	newComm = postColorizer.ReplaceAllStringFunc(newComm, helperGetColor)
	return regexp.MustCompile("\n ").ReplaceAllString(newComm, "\n")
}

func hash(id int64) int64 {
	return (((((id + 1) ^ (id >> 3)) * 7) ^ (id >> 3)) * 11)
}
func getColor(id int64) string {
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

func helperGetColor(id string) string {
	val, _ := strconv.ParseInt(id[2:], 10, 64)
	return ">>"+getColor(val)
}