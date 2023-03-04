package utils

import (
	"strconv"
	"strings"
)

func GetNodeId(ns string) (int64, error) {
	// n0, n1, n2 -> 0, 1, 2
	nids := strings.TrimLeft(ns, "n")
	nid, err := strconv.Atoi(nids)
	if err != nil {
		return 0, err
	}
	return int64(nid), nil
}
