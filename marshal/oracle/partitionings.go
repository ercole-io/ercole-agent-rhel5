// Copyright (c) 2022 Sorint.lab S.p.A.
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

package oracle

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/ercole-io/ercole-agent-rhel5/marshal"
	"github.com/ercole-io/ercole-agent-rhel5/model"
)

// Partitionings returns information about database partitionings extracted
// from the partitionings fetcher command output.
func Partitionings(cmdOutput []byte) []model.OracleDatabasePartitioning {
	partitionings := []model.OracleDatabasePartitioning{}
	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))

	for scanner.Scan() {
		partitioning := model.OracleDatabasePartitioning{}
		line := scanner.Text()
		splitted := strings.Split(line, "|||")
		if len(splitted) == 5 {
			partitioning.Owner = strings.TrimSpace(splitted[0])
			partitioning.SegmentName = strings.TrimSpace(splitted[1])
			partitioning.PartitionName = strings.TrimSpace(splitted[2])
			partitioning.SegmentType = strings.TrimSpace(splitted[3])
			partitioning.Mb = marshal.TrimParseFloat64(splitted[4])

			partitionings = append(partitionings, partitioning)
		}
	}

	return partitionings
}
