package snap

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"qp/internal/pkgdata"
	"time"
)

func fetchPackages() ([]*pkgdata.PkgInfo, error) {
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return net.Dial(networkUnix, snapdSocket)
			},
		},
		Timeout: 5 * time.Second,
	}

	snapResp, err := fetchSnaps(client)
	if err != nil {
		return nil, err
	}

	connsResp, err := fetchConnections(client)
	if err != nil {
		return nil, err
	}

	deps := parseConnections(connsResp)

	return parseSnaps(snapResp, deps)
}

func fetchSnaps(client *http.Client) (SnapdResponse, error) {
	resp, err := client.Get(snapLocalHost)
	if err != nil {
		return SnapdResponse{}, fmt.Errorf("failed to connect to snapd: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SnapdResponse{}, fmt.Errorf("snapd API error: %d", resp.StatusCode)
	}

	var snapRep SnapdResponse
	if err := json.NewDecoder(resp.Body).Decode(&snapRep); err != nil {
		return SnapdResponse{}, fmt.Errorf("failed to parse snapd response: %w", err)
	}

	return snapRep, nil
}

func fetchConnections(client *http.Client) (ConnectionsResponse, error) {
	resp, err := client.Get("http://localhost/v2/connections")
	if err != nil {
		return ConnectionsResponse{}, err
	}
	defer resp.Body.Close()

	var connsResp ConnectionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&connsResp); err != nil {
		return ConnectionsResponse{}, err
	}

	return connsResp, nil
}
