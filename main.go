package main

import(
	"github.com/moshee/go-4chan-api/api"

	"fmt"
	"os"
	"strconv"
)

func main() {
	board, thread := parseArgs()
	if board == "" && thread == "" { fmt.Println("didnt provide nor thread nor board"); return }
	if board == "" { fmt.Pritnln("how can i find a thread if no board is given"); return }
	if thread == "" { readCatalog() } else { readThread(board, thread) }
}

func parseArgs() (string, string) {
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

func readThread(b string, t string) {
	threadPosts := api.GetThread(b, t)
	for _, x := range threadPosts.Posts {
		fmt.Println(x.Id, x.Time, x.Subject, "\n",
					x.Name, "\n",
					x.Comment, "\n"
				)
	}
}