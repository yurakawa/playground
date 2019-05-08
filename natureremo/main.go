package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/tenntenn/natureremo"
)

func main() {
	cli := natureremo.NewClient("")
	t := cli.LastRateLimit.Reset.Sub(time.Now())
	time.Sleep(t / time.Duration(cli.LastRateLimit.Remaining))

	ctx := context.Background()
	ds, err := cli.DeviceService.GetAll(ctx)
	if err != nil {
		fmt.Println("err: ", err)
		os.Exit(1)
	}

	for _, d := range ds {
		te := d.NewestEvents[natureremo.SensorTypeTemperature].Value
		fmt.Println("温度: ", te, "度")

		hu := d.NewestEvents[natureremo.SensorTypeHumidity].Value
		fmt.Println("湿度: ", hu, "%")

		il := d.NewestEvents[natureremo.SensortypeIllumination].Value
		fmt.Println("照度: ", il)
	}

}
