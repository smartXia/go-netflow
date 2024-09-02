package utils

import "github.com/dustin/go-humanize"

func HumanBytes(n int64) string {
	return humanize.Bytes(uint64(n))
}
