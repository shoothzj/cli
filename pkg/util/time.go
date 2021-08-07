package util

import "time"

const (
	timeLayout = "2006-01-02 15:04:05.000000000"
)

func TimeToBytes(time time.Time) []byte {
	return []byte(time.Format(timeLayout))
}

func BytesToTime(bytes []byte) (time.Time, error) {
	s := string(bytes)
	return time.Parse("2006-01-02 15:04:05.000000000", s)
}

func BytesToTimeNoError(bytes []byte) time.Time {
	parse, err := BytesToTime(bytes)
	if err != nil {
		panic(err)
	}
	return parse
}
