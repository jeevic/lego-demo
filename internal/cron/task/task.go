package task

import (
	"fmt"

	"github.com/jeevic/lego/components/crontab"
)

func Cron(scheduler crontab.Scheduler) {
	//每一百秒执行一次
	scheduler.Every(100).Second().Do(func() {
		fmt.Println("cron Every 100s")
	})

	//每天执行一次
	scheduler.Every(1).Day().Do(func() {
		fmt.Println("cron Every days")
	})

	//十点零5秒执行
	scheduler.At("10:00:05").Do(func() {
		fmt.Println("cron Every days")
	})

}
