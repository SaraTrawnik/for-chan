package main

import(
	"github.com/moshee/go-4chan-api/api"

	"fmt"
	"os"
	"strconv"
	"regexp"
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
	for _, oneThread := range entireCatalog[0].Threads { // fix that later
		readPost(oneThread.OP)
	}

}

func readPost(p *api.Post) {
	var file string
	if p.File != nil { file = p.File.String()}
	fmt.Printf("%v %v %v\n%v%v\n", p.Name, p.Time, p.Id, file, parseComment(p.Comment))
	fmt.Printf("---\n\n")
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
		if lenOfLine > 70 {
			if char == ' ' {
				newComm += "\n"
				lenOfLine = 0
			}
		}
		newComm += string(char)
	}
	
	return regexp.MustCompile("\n ").ReplaceAllString(newComm, "\n")
}