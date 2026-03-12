package monitor

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetNetworkBytes() (uint64, uint64, error) {
	data, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return 0, 0, err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "enp7s0:") {
			fields := strings.Fields(line)
			recv, err1 := strconv.ParseUint(fields[1], 10, 64)
			send, err2 := strconv.ParseUint(fields[9], 10, 64)
			if err1 != nil || err2 != nil {
				return 0, 0, fmt.Errorf("parse error")
			}
			return recv, send, nil
		}
	}
	return 0, 0, fmt.Errorf("interface not found")
}
func CalcuilateSpeed(oldRecv, oldSend, newRecv, newSend uint64, delta float64) (float64, float64) {
	down := float64(newRecv-oldRecv) * 8 / 1_000_000 / delta // Mbps
	up := float64(newSend-oldSend) * 8 / 1_000_000 / delta   // Mbps
	return down, up
}
func MonitorNetwork() {
	var oldRecv, oldSend uint64
	oldRecv, oldSend, _ = GetNetworkBytes()
	for {
		time.Sleep(1 * time.Second)
		newRecv, newSend, err := GetNetworkBytes()
		if err != nil {
			fmt.Printf("Error:", err)
			continue
		}

		down, up := CalcuilateSpeed(oldRecv, oldSend, newRecv, newSend, 1.0)
		fmt.Printf("↓ %.2f Mbps, ↑ %.2f Mbps\n", down, up)

		oldRecv, oldSend = newRecv, newSend
	}
}
