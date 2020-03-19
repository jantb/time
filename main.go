package main

import (
	"fmt"
	"github.com/caseymrm/menuet"
	"github.com/jantb/robotgo"
	"math/rand"

	"time"
)

var lastPos = 0
var lastTime = time.Now()

type TimeStruct struct {
	ClockIn  time.Time
	ClockOut time.Time
}

var times []TimeStruct

func tracker() {
	for {
		duration := hoursForToday()
		menuet.App().SetMenuState(&menuet.MenuState{
			Title: fmtDuration(duration),
		})
		time.Sleep(time.Second)
		if !active() {
			clockOutNow()
		}
	}
}

func hoursForToday() time.Duration {
	duration := int64(0)
	for _, timeStruct := range times {
		if timeStruct.ClockOut.IsZero() {
			duration += time.Now().Sub(timeStruct.ClockIn).Nanoseconds()
			continue
		}
		duration += timeStruct.ClockOut.Sub(timeStruct.ClockIn).Nanoseconds()
	}
	return time.Duration(duration)
}

func clockInNow() {
	if len(times) == 0 || !times[len(times)-1].ClockOut.IsZero() {
		times = append(times, TimeStruct{
			ClockIn: time.Now(),
		})
	}
}

func clockOutNow() {
	if len(times) != 0 && times[len(times)-1].ClockOut.IsZero() {
		times[len(times)-1].ClockOut = time.Now()
	}
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}

func active() bool {
	pos := rand.Int()
	if pos != lastPos {
		lastPos = pos
		lastTime = time.Now()
		return true
	}
	if time.Now().Add(-15 * time.Minute).Before(lastTime) {
		return true
	}
	x, y := robotgo.GetMousePos()
	fmt.Println("pos:", x, y)
	return false
}

func menuItems() []menuet.MenuItem {
	items := []menuet.MenuItem{
		{
			Text:    "Clock in",
			Clicked: clockInNow,
		},
		{
			Text:    "Clock out",
			Clicked: clockOutNow,
		},
	}
	return items
}

func main() {
	go tracker()
	app := menuet.App()
	app.Children = menuItems
	app.RunApplication()
}
