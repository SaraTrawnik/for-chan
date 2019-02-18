package main

import(
	"github.com/moshee/go-4chan-api/api"

	"fmt"
	"os"
	"strconv"
)

func main() {
	board, thread := ParseArgs()
	if board == "" && thread == "" { fmt.Println("didnt provide nor thread nor board"); return }
	if board == "" { fmt.Pritnln("how can i find a thread if no board is given"); return }
	if thread == "" { ReadCatalog(board) } else { ReadThread(board, thread) }
}

func ParseArgs() (string, string) {
	arguments = os.Args[1:]
	var chosenBoard string
	var chosenThread string
	if x := contains("-b", arguments); x != -1 {
		if x < len(arguments) {
			if board, _ := api.LookupBoards(); board.Board == arguments[x+1] {
				chosenBoard = board.Board
			} else {
				fmt.Println("no such board")
				return "", ""
			}
		} else {
			fmt.Println("you didn't provide any board")
			return "", ""
		}
	} else {
		fmt.Println("no board given")
		return "", ""
	}

	if x := contains("-t", arguments); x != -1 {
		if x < len(arguments) {
			chosenThread, err := strconv.ParseInt(arguments[x+1], 10, 64)
			if err != nil { fmt.Println("invalid thread number"); return "", ""}
		} else {
			fmt.Println("this isn't a number dummy")
			return "", ""
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

func ReadThread(b string, t string) {
	threadPosts, err := api.GetThread(b, t)
	if err != nil {
		fmt.Pritnln("could not read thread", t, "from board", b)
		return
	}

	for _, x := range threadPosts.Posts {
		readPost(x)
	}
}

func ReadCatalog(b string) {
	entireCatalog, err := api.GetCatalog(b)
	if err != nil {
		fmt.Pritnln("could not get catalog from", b, "board")
		return
	}
	for _, oneThread := range entireCatalog.Threads { // fix that later
		readPost(oneThread.OP)
	}

}

func readPost(p api.Post) {
	fmt.Println(p.Id, p.Time, p.Subject, "\n",
				p.Name, "\n",
				p.Comment, "\n"
			)
}