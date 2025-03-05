package congrat

import (
	"database/sql"
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"net"
	"testing"
)

func TestSimple(t *testing.T) {
	geoDB, err := geoip2.Open("../resources/GeoLite2-ASN.mmdb")
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		_ = geoDB.Close()
	}()

	parsedIP := net.ParseIP("8.8.8.8")
	asn, err := geoDB.ASN(parsedIP)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(asn.AutonomousSystemNumber)
}

func TestUpdateASN(t *testing.T) {
	geoDB, err := geoip2.Open("../resources/GeoLite2-ASN.mmdb")
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		_ = geoDB.Close()
	}()

	UseDBConnection(func(db *sql.DB) error {
		return UpdateASN(db, geoDB)
	})
}
