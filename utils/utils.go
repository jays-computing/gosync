package utils

import (
	"time"
	"github.com/beevik/ntp"
)

func ToMilli(nano int64) int64 {
	return nano/int64(time.Millisecond)
}

func GetAverageNtpRtt(ntpHost string, times int32) (int32,error) {
	var totalRtt int32
	for i  := 0; i < int(times); i++ {
		t, err := ntp.Query(ntpHost, 4)
		if err != nil {
			return 0, err
		}
		totalRtt = totalRtt + int32(ToMilli(t.RTT.Nanoseconds()))
	}
	return totalRtt/times, nil
}

