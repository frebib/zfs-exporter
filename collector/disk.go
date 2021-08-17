package collector

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func diskFromPartition(path string) (string, error) {
	// Resolve dev/disk/by-* paths to dev/? first
	dev, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}
	parts := strings.SplitN(dev, "/", 3)
	// expecting []string{"", "dev", "<the bit we want>"}
	if len(parts) != 3 {
		return "", fmt.Errorf("unexpected resolved device path")
	}
	dev = "/sys/class/block/" + parts[2]

	// sys/class/block/?/.. is the dev for the disk if `path` is a partition
	uevent, err := os.Open(dev + "/../uevent")
	if errors.Is(err, os.ErrNotExist) {
		// Must be a top-level disk itself (not a partition of a disk) so go
		// straight for uevent
		uevent, err = os.Open(dev + "/uevent")
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}

	// Now find the DEVNAME=? line and return dev/?
	rd := bufio.NewScanner(uevent)
	for rd.Scan() {
		parts := strings.Split(rd.Text(), "=")
		if parts[0] == "DEVNAME" {
			path = "/dev/" + parts[1]
		}
	}

	return path, nil
}
