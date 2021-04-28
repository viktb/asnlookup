package database

import (
	"encoding/binary"
	"errors"
	"github.com/kaorimatz/go-mrt"
	"net"
)

func mrtRIBToMapping(rib *mrt.TableDumpV2RIB) (*net.IPNet, uint32, error) {
	for _, entry := range rib.RIBEntries {
		for _, attr := range entry.BGPAttributes {
			path, ok := attr.Value.(mrt.BGPPathAttributeASPath)
			if !ok {
				continue
			}

			for _, segment := range path {
				return rib.Prefix, mrtASToUint32(segment.Value[len(segment.Value)-1]), nil
			}
		}
	}

	return nil, 0, errors.New("RIB record does not contain AS path")
}

func mrtASToUint32(b mrt.AS) uint32 {
	switch len(b) {
	case 2:
		return uint32(binary.BigEndian.Uint16(b))
	case 4:
		return binary.BigEndian.Uint32(b)
	}
	return 0
}
