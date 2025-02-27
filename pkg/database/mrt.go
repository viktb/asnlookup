package database

import (
	"encoding/binary"
	"errors"
	"net"

	"github.com/viktb/go-mrt"
)

func mrtRIBToMapping(rib *mrt.TableDumpV2RIB) (*net.IPNet, uint32, error) {
	for _, entry := range rib.RIBEntries {
		var asPath mrt.BGPPathAttributeASPath
		for _, attr := range entry.BGPAttributes {
			path, ok := attr.Value.(mrt.BGPPathAttributeASPath)
			if !ok {
				continue
			}
			asPath = path
			break
		}
		if len(asPath) == 0 {
			continue
		}

		var segment *mrt.BGPASPathSegment
		for _, seg := range asPath {
			if seg.Type != mrt.BGPASPathSegmentTypeASSequence || len(seg.Value) == 0 {
				continue
			}
			segment = seg
			break
		}
		if segment == nil {
			continue
		}

		asn, ok := mrtASToUint32(segment.Value[len(segment.Value)-1])
		if !ok {
			continue
		}

		return rib.Prefix, asn, nil
	}

	return nil, 0, errors.New("RIB record does not contain a valid AS path")
}

func mrtASToUint32(b mrt.AS) (uint32, bool) {
	switch len(b) {
	case 2:
		return uint32(binary.BigEndian.Uint16(b)), true
	case 4:
		return binary.BigEndian.Uint32(b), true
	}
	return 0, false
}
