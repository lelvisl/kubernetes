package prober

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// only unix!
// https://gist.github.com/johscheuer/dc20988895d6fddfd057e221d47587d3
func (pb *prober) getContainerNetNamespace(ctx context.Context, containerID string) (string, error) {
	u, err := url.Parse(containerID)
	if err != nil {
		return "", err
	}
	r, err := pb.runtime.ContainerStatus(ctx, u.Host, true)
	if err != nil {
		return "", err
	}
	info := r.GetInfo()

	var infop any
	err = json.Unmarshal([]byte(info["info"]), &infop)
	if err != nil {
		return "", err
	}

	namespaces := infop.(map[string]any)["runtimeSpec"].(map[string]any)["linux"].(map[string]any)["namespaces"].([]any)
	for _, ns := range namespaces {
		nss := ns.(map[string]any)
		if nss["type"] == "network" {
			return nss["path"].(string), nil
		}
	}
	return "", fmt.Errorf("not found")
}
