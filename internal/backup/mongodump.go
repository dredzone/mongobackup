package backup

import (
	"fmt"
	"github.com/codeskyblue/go-sh"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
	"time"
)

type Config struct {
	Archive string 	`yaml:"archive"`
	Log     string	`yaml:"log"`
	Target  *Target	`yaml:"target"`
}

type Target struct {
	Database string `yaml:"database"`
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Params   string `yaml:"params"`
	Timeout  int	`yaml:"timeout"`
}

func mongodump(archive string, log string, cfg *Target) error {
	output, err := sh.Command("/bin/sh", "-c", "mongodump --version").CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "mongodump failed %v", err.Error())
	}

	dump := fmt.Sprintf("mongodump --archive=%v --gzip --host %v --port %v ",
		archive,
		cfg.Host,
		cfg.Port,
	)

	if cfg.Database != "" {
		dump += fmt.Sprintf("--db %v ", cfg.Database)
	}
	if cfg.Username != "" && cfg.Password != "" {
		dump += fmt.Sprintf("-u %v -p %v ", cfg.Username, cfg.Password)
	}
	if cfg.Params != "" {
		dump += fmt.Sprintf("%v", cfg.Params)
	}

	output, err = sh.Command("/bin/sh", "-c", dump).SetTimeout(time.Duration(cfg.Timeout) * time.Second).CombinedOutput()
	if err != nil {
		ex := ""
		if len(output) > 0 {
			ex = strings.Replace(string(output), "\n", " ", -1)
		}
		return errors.Wrapf(err, "mongodump log %v", ex)
	}

	return logToFile(log, output)
}

func logToFile(file string, data []byte) error {
	if len(data) > 0 {
		err := ioutil.WriteFile(file, data, 0644)
		if err != nil {
			return errors.Wrapf(err, "writing log %v failed", file)
		}
	}

	return nil
}

