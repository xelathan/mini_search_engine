package utils

import (
	"fmt"

	"github.com/robfig/cron"
	"github.com/xelathan/mini_search_engine/search"
)

func StartCronJobs() {
	c := cron.New()
	c.AddFunc("0 * * * *", search.RunEngine)
	c.AddFunc("15 * * * *", search.RunIndex)
	c.Start()
	cronCount := len(c.Entries())
	fmt.Printf("Cron Jobs Started: %d", cronCount)
}
