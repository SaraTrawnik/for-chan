package main

import (
	"github.com/SaraTrawnik/for-chan/read"

	"fmt"
	"os"
	"strconv"
)

func main() {
	board, thread := ParseArgs()
	if board == "" && thread == 0 { fmt.Println("didnt provide nor thread nor board"); return }
	if board == "" { fmt.Println("how can i find a thread if no board is given"); return }
	if thread == 0 { read.Catalog(board) } else { read.Thread(board, thread) }
}

func ParseArgs() (string, int64) {
	contains := func(a string, in []string) int { for place, arg := range in { if arg == a { return place } } return -1 }
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