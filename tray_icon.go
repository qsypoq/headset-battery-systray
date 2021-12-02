package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"time"

	"github.com/getlantern/systray"
)

var (
	battery_icon      string
	headset_connected bool    = false
	headset_charging  bool    = false
	headset_battery   int     = 0
	returned_battery  float64 = 0
)

func main() {
	systray.Run(onReady, onExit)
}

func Round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

func generate_tray_level(level int) []byte {
	txtimg, _ := txt_on_img(
		edit_req{
			BgImgPath: "assets/normal.png",
			FontPath:  "assets/FiraSans-SemiBold.ttf",
			FontSize:  256,
			Text:      fmt.Sprint(level),
		},
	)

	bufff := new(bytes.Buffer)
	err := convert_to_ico(bufff, img_to_nrgba(txtimg))

	if err != nil {
		fmt.Println("failed to create buffer", err)
	}

	return bufff.Bytes()
}

func onReady() {
	tooltip := "Headset is Disconnected"
	battery_icon_path := "assets/disconnected.ico"
	systray.SetIcon(getIcon(battery_icon_path))
	mQuit := systray.AddMenuItem("Quit", "Quit the app")

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	go func() {
		for {
			buf := open_headset()

			if buf[1] == 1 {
				headset_connected = true
				returned_battery = float64(buf[3])

				if buf[4] == 1 {
					headset_charging = true
					battery_icon_path = "assets/charging.ico"
				} else {
					headset_charging = false
				}

				headset_battery = convert_battery(returned_battery, headset_charging)

				if buf[4] == 0 {
					systray.SetIcon(generate_tray_level(headset_battery))
				} else {
					systray.SetIcon(getIcon(battery_icon_path))
				}
			}

			tooltip = fmt.Sprint(Round(float64(headset_battery), 1))
			systray.SetTooltip(tooltip)
			time.Sleep(5 * time.Second)
		}
	}()
}

func onExit() {
}

func getIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		fmt.Print(err)
	}
	return b
}
