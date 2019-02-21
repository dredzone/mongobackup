package backup

import (
	"github.com/dredzone/mongobackup/internal/notification/email"
	"github.com/dredzone/mongobackup/internal/storage/aws"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Plan struct {
	Name   	string          `yaml:"name"`
	Target 	*Target 		`yaml:"target"`
	Backup	*Backup         `yaml:"backup"`
	S3      *aws.S3    		`yaml:"s3"`
	SMTP    *email.SMTP     `yaml:"smtp"`
}

type Backup struct {
	Retention int    `yaml:"retention"`
}

func loadPlan(config string) (*Plan, error) {
	plan := &Plan{}
	abs, _ := filepath.Abs(config)
	data, err := ioutil.ReadFile(abs)
	if err != nil {
		return plan, errors.Wrapf(err, "Reading %v failed", config)
	}

	if err := yaml.Unmarshal(data, &plan); err != nil {
		return plan, errors.Wrapf(err, "Parsing %v failed", config)
	}

	_, filename := filepath.Split(config)
	plan.Name = strings.TrimSuffix(filename, filepath.Ext(filename))

	return plan, nil
}

func loadPlans(dir string) ([]Plan, error) {
	files := make([]string, 0)

	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, "yml") || strings.Contains(path, "yaml") {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, errors.Wrapf(err, "Reading from %v failed", dir)
	}

	plans := make([]Plan, 0)

	for _, path := range files {
		plan, err := loadPlan(path)
		if err != nil {
			return nil, err
		}
		duplicate := false
		for _, p := range plans {
			if p.Name == plan.Name {
				duplicate = true
				break
			}
		}
		if duplicate {
			continue
		}

		plans = append(plans, *plan)

	}
	if len(plans) < 1 {
		return nil, errors.Errorf("No backup plans found in %v", dir)
	}

	return plans, nil
}

