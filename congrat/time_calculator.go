package congrat

import "encoding/binary"

type IPTimestamps struct {
	ID        int64
	IPAddress string
	AeadID    int64
	T1        []byte
	RealT1    []byte
	T2        []byte
	T3        []byte
	T4        []byte
}

func GetServerExtensionFieldsCost(t1, t2, t3, t4, ntsT3, ntsT4 []byte) (delay, offset, cost int64) {
	t164 := int64(binary.BigEndian.Uint64(t1))
	t264 := int64(binary.BigEndian.Uint64(t2))
	t364 := int64(binary.BigEndian.Uint64(t3))
	t464 := int64(binary.BigEndian.Uint64(t4))
	ntsT364 := int64(binary.BigEndian.Uint64(ntsT3))
	ntsT464 := int64(binary.BigEndian.Uint64(ntsT4))
	delay = (t264 - t164) + (t464-t364)/2
	offset = (t264 - t164) + (t364-t464)/2
	cost = ntsT464 - ntsT364 - (t464 - t364)
	return
}
