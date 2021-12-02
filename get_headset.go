package main

import (
	"fmt"
	"math"

	"github.com/qsypoq/hid"
)

var (
	vid uint = 4152
	pid uint = 4802
)

func convert_battery(value float64, headset_charging bool) int {
	if headset_charging != false {
		value = value - 21
	}
	return int(float64((-0.005442*(math.Pow(value, float64(2))) + (float64(3.196) * value)) - 264.9))
}

func open_headset() []byte {
	buf := make([]byte, 9)
	buf[0] = 0x0
	buf[1] = 0x20

	devices := hid.Enumerate(uint16(vid), uint16(pid))
	if len(devices) > 2 {
		headset_info := devices[2]

		headset, err := headset_info.Open()

		if err != nil {
			fmt.Println(err)
		}

		defer headset.Close()

		_, err = headset.Write(buf)
		if err != nil {
			fmt.Printf("Write fail: %v", err)
		}

		_, err = headset.Read(buf)
		if err != nil {
			fmt.Printf("Read fail: %v", err)
		}
	}

	return buf
}
