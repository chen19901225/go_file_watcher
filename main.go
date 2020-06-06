package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/radovskyb/watcher"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Required: true,
				Usage:    "Load configuration from  from `File`",
			},
		},
		Name:  "go_file_watcher",
		Usage: "go_file_watcher -c **.conf",
		Action: func(c *cli.Context) error {
			config := c.String("config")
			fmt.Printf("config:%s\n", config)

			_, err := os.Stat(config)

			if err != nil {

				log.Fatal(fmt.Sprintf("Fail to go config :%s", config))
			}

			var dat map[string]interface{}
			byt, err := ioutil.ReadFile(config)
			if err != nil {
				log.Fatal(fmt.Sprintf("Fail to read config content:%s", config))
			}

			if err := json.Unmarshal(byt, &dat); err != nil {
				log.Fatal(fmt.Sprintf("Fail to parse config content:%s", byt))
			}

			directory, ok := dat["directory"]

			if ok == false {
				log.Fatal(fmt.Sprintf("missing direcotry in config:%+v", dat))
			}

			commandList, ok := dat["command_list"]
			if ok == false {
				log.Fatal(fmt.Sprintf("missing command_list in config:%+v", dat))
			}
			directory = directory.(string)
			fmt.Println(fmt.Sprintf("Directory:%s", directory))

			w := watcher.New()

			// SetMaxEvents to 1 to allow at most 1 event's to be received
			// on the Event channel per watching cycle.
			//
			// If SetMaxEvents is not set, the default is to send all events.
			w.SetMaxEvents(1)

			// Only notify rename and move events.
			w.FilterOps(watcher.Rename, watcher.Move, watcher.Create,
				watcher.Remove, watcher.Write)

			// Only files that match the regular expression during file listings
			// will be watched.
			//r := regexp.MustCompile(".*")
			// w.AddFilterHook(watcher.RegexFilterHook(r, false))

			completeChannel := make(chan int, 0)

			go func() {
				for {
					select {
					case event := <-w.Event:
						//fmt.Println(fmt.Sprintf("path:%s, oldPath:%s", event.Path, event.OldPath))
						handleEvent(event, commandList.([]interface{}), directory.(string), w)
						//fmt.Println(event) // Print the event's info.
					case err := <-w.Error:
						log.Fatalln(err)
						completeChannel <- 1
					case <-w.Closed:
						completeChannel <- 1
						return
					case <-time.After(5 * time.Second):
						continue

					}
				}
			}()

			// Watch test_folder recursively for changes.
			if err := w.AddRecursive(directory.(string)); err != nil {
				log.Fatalln(err)
			}

			if err := w.Start(time.Millisecond * 200); err != nil {
				log.Fatalln(err)
			}
			<-completeChannel
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
