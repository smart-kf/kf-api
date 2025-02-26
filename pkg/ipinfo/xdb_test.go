package ipinfo

import (
	"testing"
)

func TestGetLocation(t *testing.T) {
	GetLocation("/Users/zmj/Desktop/projects/std-api/ip2region.xdb", "125.70.175.75")
}
