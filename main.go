package main

import(
	"github.com/moshee/go-4chan-api/api"

	"fmt"
	"os"
	"strconv"
)

func main() {
	arguments = os.Args[1:]
	board := arguments[0]
	thread, err := strconv.ParseInt(arguments[1], 10, 64)
	if err != nil { fmt.Println("invalid thread number"); return }
	threadPosts := api.GetThread(board, thread)
	for _, x := range thread.Posts {
		fmt.Println(x.Id, x.Time, x.Subject, "\n",
					x.Name, "\n",
					x.Comment, "\n"
				)
	}
}