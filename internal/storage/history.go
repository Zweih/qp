package storage

import (
	"fmt"
	"os"
	pb "qp/internal/protobuf"

	"google.golang.org/protobuf/proto"
)

func SaveInstallHistory(cacheRoot string, history map[string]int64, latestLogTime int64) error {
	historyPath := cacheRoot + dotHistory
	installHistory := &pb.InstallHistory{
		InstallTimestamps:  history,
		Version:            historyVersion,
		LatestLogTimestamp: latestLogTime,
	}

	byteData, err := proto.Marshal(installHistory)
	if err != nil {
		return fmt.Errorf("failed to marshal history: %v", err)
	}

	return os.WriteFile(historyPath, byteData, 0644)
}

func LoadInstallHistory(cacheRoot string) (map[string]int64, int64, error) {
	historyPath := cacheRoot + dotHistory
	if _, err := os.Stat(historyPath); os.IsNotExist(err) {
		return make(map[string]int64), 0, nil
	}

	byteData, err := os.ReadFile(historyPath)
	if err != nil {
		return nil, 0, err
	}

	var installHistory pb.InstallHistory
	err = proto.Unmarshal(byteData, &installHistory)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to unmarshal install history: %v", err)
	}

	if installHistory.Version != historyVersion {
		return make(map[string]int64), 0, nil
	}

	return installHistory.InstallTimestamps, installHistory.LatestLogTimestamp, nil
}
