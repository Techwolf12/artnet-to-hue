package hue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/techwolf12/artnet-to-hue/pkg/config"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

type HueBridgeInfo struct {
	ID                string `json:"id"`
	InternalIPAddress string `json:"internalipaddress"`
	Port              int    `json:"port"`
}

type EntertainmentLightState struct {
	On  bool       `json:"on"`
	Bri int        `json:"bri"`
	XY  [2]float64 `json:"xy"`
}

func EntertainmentSend(conf config.Config, states []EntertainmentLightState) {
	url := fmt.Sprintf("http://%s/api/%s/groups/%s/action", conf.HueBridgeIP, conf.Username, conf.EntertainmentZone)

	body, err := json.Marshal(states)
	if err != nil {
		log.Printf("Failed to marshal light states: %v", err)
		return
	}

	req, err := http.NewRequest("PUT", url, bytes.NewReader(body))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request to Hue: %v", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected response from Hue: %v", resp.Status)
	}
}

func GetHueUsername(bridgeIP net.IP, deviceType string) (string, error) {
	url := fmt.Sprintf("http://%s/api", bridgeIP)
	body := fmt.Sprintf(`{"devicetype":"%s"}`, deviceType)
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response from bridge: %s", resp.Status)
	}

	var result []map[string]map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	for _, item := range result {
		if success, ok := item["success"]; ok {
			if username, exists := success["username"]; exists {
				return fmt.Sprintf("%v", username), nil
			}
		}
		if errMsg, ok := item["error"]; ok {
			return "", fmt.Errorf("error from bridge: %v", errMsg["description"])
		}
	}
	return "", fmt.Errorf("unexpected response from bridge")
}

func DiscoverBridges() ([]HueBridgeInfo, error) {
	url := "https://discovery.meethue.com/"
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bridges: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	var bridges []HueBridgeInfo
	if err := json.NewDecoder(resp.Body).Decode(&bridges); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return bridges, nil
}
