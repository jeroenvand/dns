package dns

// A structure for handling zone data

import (
	"github.com/miekg/radix"
)

type Zone struct {
	Name         string // Name of the zone
	*radix.Radix        // Zone data
}

type ZoneData struct {
	Name string          // Domain name for this node
	RR   map[uint16][]RR // Map of the RR type to the RR
	// DNSSEC signatures for the RRsets
	Signatures []*RR_RRSIG
	// Almost always true, except for non-origin NS records (and accompanying glue)
	Authoritatve bool
}

// New ...
func NewZone(origin string) *Zone {
	z := new(Zone)
	z.Origin = origin
	z.Radix = radix.New()
	return z
}

func (z *Zone) Insert(r RR) {
	zd := z.Radix.Find(r.Header().Name)
	if zd == nil {
		zd := new(ZoneData)
		zd.Name = r.Header().Name
		zd.RR = make(map[uint16][]RR)
		zd.Signatures = make([]*RR_RRSIG, 0)
		switch t := r.Header().Rrtype; t {
		case TypeRRSIG:
			zd.Signatures = append(zd.Signatures, r.(*RR_RRSIG))
		default:
			zd.RR[t] = append(zd.RR[t], r)
		}
		z.Radix.Insert(r.Header().Name, zd)
		return
	}
	switch t := r.Header().Rrtype; t {
	case TypeRRSIG:
		zd.(*ZoneData).Signatures = append(zd.(*ZoneData).Signatures, r.(*RR_RRSIG))
	default:
		zd.(*ZoneData).RR[t] = append(zd.(*ZoneData).RR[t], r)
	}
	// TODO(mg): Glue
	return
}

func (z *Zone) Remove(r RR) {

}