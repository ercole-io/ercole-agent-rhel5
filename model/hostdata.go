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

// SchemaVersion contains the version of the schema
const SchemaVersion int = 1

// HostData holds all informations about a host & services
type HostData struct {
	Hostname                string                  `json:"hostname"`
	Location                string                  `json:"location"`
	Environment             string                  `json:"environment"`
	Tags                    []string                `json:"tags"`
	AgentVersion            string                  `json:"agentVersion"`
	SchemaVersion           int                     `json:"schemaVersion"`
	Period                  uint                    `json:"period"`
	Info                    Host                    `json:"info"`
	ClusterMembershipStatus ClusterMembershipStatus `json:"clusterMembershipStatus"`
	Features                Features                `json:"features"`
	Filesystems             []Filesystem            `json:"filesystems"`
	Clusters                []ClusterInfo           `json:"clusters"`
	OtherInfo               map[string]interface{}  `json:"-"`
}
