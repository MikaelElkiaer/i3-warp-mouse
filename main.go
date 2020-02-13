package main

import (
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"go.i3wm.org/i3/v4"
)

var lastUpdate int64 = 0

func main() {
	recv := i3.Subscribe(i3.WindowEventType)
	for recv.Next() {
		if e, ok := recv.Event().(*i3.WindowEvent); ok && e.Change == "focus" {
			log.Printf("Change: %s", e.Change)
			if time.Now().UnixNano()-lastUpdate > 10000000 {
				warpMouse(e.Container)
				lastUpdate = time.Now().UnixNano()
			}
		}
	}
}

func warpMouse(container i3.Node) {
	log.Printf("Window position: %dx%d (%d)", container.WindowRect.Width, container.WindowRect.Y, container.Window)

	mouseOutput, _ := exec.Command("xdotool", "getmouselocation", "--shell").Output()
	r := regexp.MustCompile(`(?s)X=(\d+).*Y=(\d+).*WINDOW=(\d+)`)
	c := r.FindStringSubmatch(string(mouseOutput))
	mouseX, _ := strconv.ParseInt(c[1], 10, 64)
	mouseY, _ := strconv.ParseInt(c[2], 10, 64)
	mouseWindow, _ := strconv.ParseInt(c[3], 10, 64)

	log.Printf("Mouse position: %dx%d (%d)", mouseX, mouseY, mouseWindow)

	if mouseWindow != container.Window {
		window := strconv.FormatInt(container.Window, 10)
		newMouseX := strconv.FormatInt(container.WindowRect.Width/2, 10)
		newMouseY := strconv.FormatInt(container.WindowRect.Height/2, 10)

		log.Printf("New mouse position: %sx%s (%s)", newMouseX, newMouseY, window)

		err := exec.Command("xdotool", "mousemove", "-window", window, newMouseX, newMouseY).Run()
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
