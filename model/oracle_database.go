// Copyright (c) 2020 Sorint.lab S.p.A.
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

package model

import (
	"strings"
)

// OracleDatabase holds information about an Oracle database.
type OracleDatabase struct {
	InstanceNumber    int                               `json:"instanceNumber" bson:"instanceNumber"`
	InstanceName      string                            `json:"instanceName" bson:"instanceName"`
	Name              string                            `json:"name" bson:"name"`
	UniqueName        string                            `json:"uniqueName" bson:"uniqueName"`
	Status            string                            `json:"status" bson:"status"`
	DbID              uint                              `json:"dbID" bson:"dbID"`
	Role              string                            `json:"role" bson:"role"`
	IsCDB             bool                              `json:"isCDB" bson:"isCDB"`
	Version           string                            `json:"version" bson:"version"`
	Platform          string                            `json:"platform" bson:"platform"`
	Archivelog        bool                              `json:"archivelog" bson:"archivelog"`
	Charset           string                            `json:"charset" bson:"charset"`
	NCharset          string                            `json:"nCharset" bson:"nCharset"`
	BlockSize         int                               `json:"blockSize" bson:"blockSize"`
	CPUCount          int                               `json:"cpuCount" bson:"cpuCount"`
	SGATarget         float64                           `json:"sgaTarget" bson:"sgaTarget"`
	PGATarget         float64                           `json:"pgaTarget" bson:"pgaTarget"`
	MemoryTarget      float64                           `json:"memoryTarget" bson:"memoryTarget"`
	SGAMaxSize        float64                           `json:"sgaMaxSize" bson:"sgaMaxSize"`
	SegmentsSize      float64                           `json:"segmentsSize" bson:"segmentsSize"`
	DatafileSize      float64                           `json:"datafileSize" bson:"datafileSize"`
	Allocable         float64                           `json:"allocable" bson:"allocable"`
	Elapsed           *float64                          `json:"elapsed" bson:"elapsed"`
	DBTime            *float64                          `json:"dbTime" bson:"dbTime"`
	DailyCPUUsage     *float64                          `json:"dailyCPUUsage" bson:"dailyCPUUsage"`
	Work              *float64                          `json:"work" bson:"work"`
	ASM               bool                              `json:"asm" bson:"asm"`
	Dataguard         bool                              `json:"dataguard" bson:"dataguard"`
	Patches           []OracleDatabasePatch             `json:"patches" bson:"patches"`
	Tablespaces       []OracleDatabaseTablespace        `json:"tablespaces" bson:"tablespaces"`
	Schemas           []OracleDatabaseSchema            `json:"schemas" bson:"schemas"`
	Licenses          []OracleDatabaseLicense           `json:"licenses" bson:"licenses"`
	ADDMs             []OracleDatabaseAddm              `json:"addms" bson:"addms"`
	SegmentAdvisors   []OracleDatabaseSegmentAdvisor    `json:"segmentAdvisors" bson:"segmentAdvisors"`
	PSUs              []OracleDatabasePSU               `json:"psus" bson:"psus"`
	Backups           []OracleDatabaseBackup            `json:"backups" bson:"backups"`
	FeatureUsageStats []OracleDatabaseFeatureUsageStat  `json:"featureUsageStats" bson:"featureUsageStats"`
	PDBs              []OracleDatabasePluggableDatabase `json:"pdbs" bson:"pdbs"`
	Services          []OracleDatabaseService           `json:"services" bson:"services"`
	OtherInfo         map[string]interface{}            `json:"-" bson:"-"`
}

var (
	OracleDatabaseStatusOpen    = "OPEN"
	OracleDatabaseStatusMounted = "MOUNTED"
)

var (
	OracleDatabaseRolePrimary         = "PRIMARY"
	OracleDatabaseRoleLogicalStandby  = "LOGICAL STANDBY"
	OracleDatabaseRolePhysicalStandby = "PHYSICAL STANDBY"
	OracleDatabaseRoleSnapshotStandby = "SNAPSHOT STANDBY"
)

var (
	OracleDatabaseEditionEnterprise = "ENT"
	OracleDatabaseEditionExtreme    = "EXE"
	OracleDatabaseEditionStandard   = "STD"
)

func (v OracleDatabase) Edition() (dbEdition string) {
	if strings.Contains(strings.ToUpper(v.Version), "ENTERPRISE") {
		dbEdition = OracleDatabaseEditionEnterprise
	} else if strings.Contains(strings.ToUpper(v.Version), "EXTREME") {
		dbEdition = OracleDatabaseEditionExtreme
	} else {
		dbEdition = OracleDatabaseEditionStandard
	}

	return
}

func (v OracleDatabase) CoreFactor(host Host) float64 {
	dbEdition := v.Edition()
	coreFactor := float64(-1)

	if host.HardwareAbstractionTechnology == HardwareAbstractionTechnologyOvm ||
		host.HardwareAbstractionTechnology == HardwareAbstractionTechnologyVmware ||
		host.HardwareAbstractionTechnology == HardwareAbstractionTechnologyVmother {

		if dbEdition == OracleDatabaseEditionExtreme || dbEdition == OracleDatabaseEditionEnterprise {
			coreFactor = 0.5
		} else if dbEdition == OracleDatabaseEditionStandard {
			coreFactor = 0
		}

	} else if host.HardwareAbstractionTechnology == HardwareAbstractionTechnologyPhysical {
		if dbEdition == OracleDatabaseEditionExtreme || dbEdition == OracleDatabaseEditionEnterprise {
			coreFactor = 0.5
		} else if dbEdition == OracleDatabaseEditionStandard {
			coreFactor = float64(host.CPUSockets)
		}
	}

	return coreFactor
}

// DatabasesArrayAsMap return the equivalent map of the database array with Database.Name as Key
func DatabasesArrayAsMap(dbs []OracleDatabase) map[string]OracleDatabase {
	out := make(map[string]OracleDatabase)
	for _, db := range dbs {
		out[db.Name] = db
	}
	return out
}

// HasEnterpriseLicense return true if the database has enterprise license.
func HasEnterpriseLicense(db OracleDatabase) bool {
	//The database may not support the "license" feature
	if db.Licenses == nil {
		return false
	}

	//Search for a enterprise license
	for _, lic := range db.Licenses {
		if (lic.Name == "Oracle ENT" || lic.Name == "oracle ENT" || lic.Name == "Oracle EXT" || lic.Name == "oracle EXT") && lic.Count > 0 {
			return true
		}
	}

	return false
}
