package hue

import (
	"crypto/tls"
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

type BridgeInfo struct {
	ID                string `json:"id"`
	InternalIPAddress string `json:"internalipaddress"`
	Port              int    `json:"port"`
}

func GetHueApplicationID(config config.Config) (string, error) {
	url := fmt.Sprintf("https://%s/auth/v1", config.HueBridgeIP)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("hue-application-key", config.Username)

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	appID := resp.Header.Get("hue-application-id")
	if appID == "" {
		return "", fmt.Errorf("hue-application-id not found in response headers")
	}
	return appID, nil
}

func GetHueUsername(bridgeIP net.IP, deviceType string) (string, string, error) {
	url := fmt.Sprintf("https://%s/api", bridgeIP)
	body := fmt.Sprintf(`{"devicetype":"%s","generateclientkey":true}`, deviceType)
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to send request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("unexpected response from bridge: %s", resp.Status)
	}

	var result []map[string]map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", fmt.Errorf("failed to decode response: %w", err)
	}

	for _, item := range result {
		if success, ok := item["success"]; ok {
			if username, exists := success["username"]; exists {
				if clientKey, exists := success["clientkey"]; exists {
					return fmt.Sprintf("%v", username), fmt.Sprintf("%v", clientKey), nil
				}
				return "", "", fmt.Errorf("clientkey not found in response")
			}
		}
		if errMsg, ok := item["error"]; ok {
			return "", "", fmt.Errorf("error from bridge: %v", errMsg["description"])
		}
	}
	return "", "", fmt.Errorf("unexpected response from bridge")
}

type EntertainmentConfig struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func GetEntertainmentConfigID(config config.Config) ([]EntertainmentConfig, error) {
	url := fmt.Sprintf("https://%s/clip/v2/resource/entertainment_configuration", config.HueBridgeIP)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("hue-application-key", config.Username)
	client := &http.Client{Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	var configResponse EntertainmentConfigResponse
	if err := json.NewDecoder(resp.Body).Decode(&configResponse); err != nil {
		return nil, err
	}
	if len(configResponse.Data) == 0 {
		return nil, fmt.Errorf("no entertainment configuration found")
	}
	var entertainmentConfigs []EntertainmentConfig
	for _, data := range configResponse.Data {
		entertainmentConfigs = append(entertainmentConfigs, EntertainmentConfig{
			ID:     data.ID,
			Name:   data.Name,
			Status: data.Status, // Assuming all returned configs are active
		})
	}
	return entertainmentConfigs, nil
}

func DiscoverBridges() ([]BridgeInfo, error) {
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

	var bridges []BridgeInfo
	if err := json.NewDecoder(resp.Body).Decode(&bridges); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return bridges, nil
}
