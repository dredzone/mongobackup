package backup

import (
	"fmt"
	"github.com/dredzone/mongobackup/internal/notification/email"
	"github.com/dustin/go-humanize"
	"github.com/hako/durafmt"
	"github.com/sirupsen/logrus"
)

func notify(result *Result) {
	message := fmt.Sprintf("Backup %v finished with status %v in %v archive %v size %v",
		result.Plan.Name,
		result.Status,
		durafmt.ParseShort(result.Duration),
		result.Name,
		humanize.Bytes(uint64(result.Size)))

	if result.Plan.SMTP != nil {
		if err := email.Send(fmt.Sprintf("Backup for %s finished", result.Plan.Name),
			message, result.Plan.SMTP); err != nil {
				logrus.Errorf(fmt.Sprintf("Notification failed %v", err.Error()))
		}
	}

	logrus.Infof(message)
}

