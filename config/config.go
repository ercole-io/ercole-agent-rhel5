// Copyright (c) 2019 Sorint.lab S.p.A.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ercole-io/ercole-agent-rhel5/logger"
)

// Configuration holds the agent configuration options
type Configuration struct {
	Hostname               string
	Environment            string
	Location               string
	DataserviceURL         string
	AgentUser              string
	AgentPassword          string
	EnableServerValidation bool
	ForcePwshVersion       string
	Period                 uint
	Verbose                bool
	ParallelizeRequests    bool
	LogDirectory           string
	Features               Features
}

// Features holds features params
type Features struct {
	OracleDatabase     OracleDatabaseFeature
	Virtualization     VirtualizationFeature
	OracleExadata      OracleExadataFeature
	MicrosoftSQLServer MicrosoftSQLServerFeature
}

// OracleDatabaseFeature holds oracle database feature params
type OracleDatabaseFeature struct {
	Enabled     bool
	FetcherUser string
	Oratab      string
	AWR         int
	Forcestats  bool
}

// VirtualizationFeature holds virtualization feature params
type VirtualizationFeature struct {
	Enabled     bool
	FetcherUser string
	Hypervisors []Hypervisor
}

// Hypervisor holds the parameters used to connect to an hypervisor
type Hypervisor struct {
	Type       string
	Endpoint   string
	Username   string
	Password   string
	OvmUserKey string
	OvmControl string
}

// OracleExadataFeature holds oracle exadata feature params
type OracleExadataFeature struct {
	Enabled     bool
	FetcherUser string
}

// MicrosoftSQLServerFeature holds microsoft sql server feature params
type MicrosoftSQLServerFeature struct {
	Enabled     bool
	FetcherUser string
}

// ReadConfig reads the configuration file from the current dir
// or /opt/ercole-agent
func ReadConfig(log logger.Logger) Configuration {
	baseDir := GetBaseDir()
	configFile := ""

	configFile = baseDir + "/config.json"

	ex := exists(configFile)
	if !ex {
		configFile = "/opt/ercole-agent/config.json"
	}
	raw, err := ioutil.ReadFile(configFile)

	if err != nil {
		log.Fatal("Unable to read configuration file", err)
	}

	var conf Configuration
	err = json.Unmarshal(raw, &conf)

	if err != nil {
		log.Fatal("Unable to parse configuration file", err)
	}

	checkConfiguration(log, &conf)

	return conf
}

func exists(name string) bool {
	_, err := os.Stat(name)

	return err == nil
}

func checkConfiguration(log logger.Logger, config *Configuration) {
	checkPeriod(log, config)
	checkLogDirectory(log, config)

	if config.Features.OracleDatabase.Oratab == "" {
		config.Features.OracleDatabase.Oratab = "/etc/oratab"
	}
}

func checkPeriod(log logger.Logger, config *Configuration) {
	if config.Period == 0 {
		defaultPeriod := uint(24)
		log.Warnf("Period has invalid value [%d], set to default value [%d]", config.Period, defaultPeriod)
		config.Period = defaultPeriod
	}
}

func checkLogDirectory(log logger.Logger, config *Configuration) {
	path := config.LogDirectory
	if path == "" {
		return
	}

	isWritable, err := isDirectoryWritable(path)
	if err != nil {
		log.Fatal("LogDirectory is not valid: ", err)
	}

	if !isWritable {
		log.Fatal("LogDirectory is not writable")
	}
}

// GetBaseDir return executable base directory, os independant
func GetBaseDir() string {
	s, _ := os.Readlink("/proc/self/exe")
	s = filepath.Dir(s)

	return s
}
