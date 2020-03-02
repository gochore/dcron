package dcron

import "github.com/robfig/cron/v3"

//go:generate mockgen -source=entry_getter.go -destination mock_dcron/entry_getter.go

type entryGetter interface {
	Entry(id cron.EntryID) cron.Entry
}
