package cmd

import (
	"fmt"
	"github.com/dredzone/mongobackup/internal/backup"
	"github.com/dredzone/mongobackup/internal/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	backupDir	string
	configFile 	string
	RootCmd = &cobra.Command{
		Use:   "mongobackup",
		Short: "A convenient wrapper around mongodump to perform database backups",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initLogFile()
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := util.MakeDirIfNotExists(backupDir, os.ModePerm); err != nil {
				logrus.Fatal("Cannot create %v directory.", backupDir)
			}
			backup.Run(configFile, backupDir, onBackup)
		},
	}
)

func init() {
	f := RootCmd.PersistentFlags()
	f.StringVarP(
		&backupDir,
		"backupDir",
		"d",
		filepath.Join(homeDir(), "mongobackup","data", "backup"),
		"backup directory",
	)

	f.StringVarP(
		&configFile,
		"config",
		"c",
		filepath.Join(homeDir(), "mongobackup","data", "config"),
		"config path to a file or directory",
	)

	logrus.SetFormatter(new(logrus.TextFormatter))
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout)
}

func initLogFile() error {
	workDir, err := backup.WorkDir()
	if err != nil {
		return err
	}
	file, err := os.OpenFile(filepath.Join(workDir, "backup.log"), os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	logrus.SetOutput(file)
	return nil
}

func onFail(result *backup.Result, err error) {
	if result.Name == "" {
		logrus.Errorf(fmt.Sprintf("Backup failed %v", err.Error()))
		return
	}
	logrus.Errorf(fmt.Sprintf("Backup %v failed %v", result.Plan.Name, err.Error()))

}

func onBackup(result *backup.Result, err error) {
	if err != nil {
		if result.Name == "" {
			logrus.Errorf(fmt.Sprintf("Backup failed %v", err.Error()))
		} else {
			logrus.Errorf(fmt.Sprintf("Backup %v failed %v", result.Plan.Name, err.Error()))
		}
		return
	}
	logrus.Infof(fmt.Sprintf("Backup for %v completed", result.Plan.Name))
}

func homeDir() string {
	dir, err := util.HomeDir()
	if err != nil {
		logrus.Fatal(err)
	}
	return dir
}
