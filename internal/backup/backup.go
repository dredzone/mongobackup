package backup

import (
	"fmt"
	"github.com/codeskyblue/go-sh"
	"github.com/dredzone/mongobackup/internal/util"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"time"
)

type Status string

const (
	SuccessBackup Status = "SUCCESS"
	FailBackup Status = "FAIL"
)

type Result struct {
	Name      string        	`json:"name"`
	Plan      *Plan        		`json:"plan"`
	Duration  time.Duration 	`json:"duration"`
	Size      int64         	`json:"size"`
	Status    Status      		`json:"status"`
	Timestamp time.Time     	`json:"timestamp"`
}

func Run(config string, dir string, done func(result *Result, err error)) {
	var (
		plans []Plan
		err error
	)
	if util.DirExists(config) {
		plans, err = loadPlans(config)
		if err != nil {
			done(&Result{}, err)
		}
	} else {
		var plan, err = loadPlan(config)
		if err != nil {
			done(&Result{}, err)
		}
		plans = append(plans, *plan)
	}
	for _, plan := range plans {
		result, err := backup(dir, &plan)
		if err != nil {
			done(result, err)
			continue
		}
		notify(result)
		done(result, nil)
	}
}

func WorkDir() (string, error) {
	home, err := util.HomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, "mongobackup", ".backup")
	if err := util.MakeDirIfNotExists(dir, os.ModePerm); err != nil {
		return "", errors.Wrapf(err, "failed to create directory %v", err)
	}
	return dir, nil
}

func backup(dir string, plan *Plan) (*Result, error) {
	t1 := time.Now()

	res := &Result{
		Plan:      plan,
		Timestamp: t1.UTC(),
		Status:    FailBackup,
	}

	workDir, err := WorkDir()
	if err != nil {
		return res, err
	}

	archive := fmt.Sprintf("%v/%v-%v.gz", workDir, plan.Name, t1.Unix())
	log := fmt.Sprintf("%v/%v-%v.log", workDir, plan.Name, t1.Unix())
	err = mongodump(archive, log, plan.Target)
	if err != nil {
		return res, err
	}
	_, res.Name = filepath.Split(archive)
	backupDir := fmt.Sprintf("%v/%v", dir, plan.Name)
	err = sh.Command("mkdir", "-p", backupDir).Run()
	if err != nil {
		return res, errors.Wrapf(err, "creating dir %v in %v failed", plan.Name, dir)
	}

	fi, err := os.Stat(archive)
	if err != nil {
		return res, errors.Wrapf(err, "stat file %v failed", archive)
	}
	res.Size = fi.Size()
	files := [2]string{archive, log}
	for _, file := range files {
		err = sh.Command("mv", file, backupDir).Run()
		if err != nil {
			return res, errors.Wrapf(err, "moving file from %v to %v failed", archive, backupDir)
		}
	}

	if plan.Backup.Retention > 0 {
		err = applyRetention(backupDir, plan.Backup.Retention)
		if err != nil {
			return res, errors.Wrap(err, "retention job failed")
		}
	}

	t2 := time.Now()
	res.Status = SuccessBackup
	res.Duration = t2.Sub(t1)
	return res, nil
}

func applyRetention(path string, retention int) error {
	gz := fmt.Sprintf("cd %v && rm -f $(ls -1t *.gz | tail -n +%v)", path, retention+1)
	err := sh.Command("/bin/sh", "-c", gz).Run()
	if err != nil {
		return errors.Wrapf(err, "removing old gz files from %v failed", path)
	}

	log := fmt.Sprintf("cd %v && rm -f $(ls -1t *.log | tail -n +%v)", path, retention+1)
	err = sh.Command("/bin/sh", "-c", log).Run()
	if err != nil {
		return errors.Wrapf(err, "removing old log files from %v failed", path)
	}

	return nil
}

