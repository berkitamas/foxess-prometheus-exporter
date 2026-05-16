package foxess

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const BaseUrl = "https://www.foxesscloud.com"

type Inverter struct {
	DeviceSN    string `json:"deviceSN"`
	DeviceType  string `json:"deviceType"`
	ProductType string `json:"productType"`
}

type RealTimeData struct {
	DeviceSN string `json:"deviceSN"`
	Datas    []struct {
		Variable string      `json:"variable"`
		Value    interface{} `json:"value"`
	} `json:"datas"`
}

type ApiResponse struct {
	Errno  int             `json:"errno"`
	Result json.RawMessage `json:"result"`
}

func calculateSignature(path, apiKey string, timestamp int64) string {
	data := fmt.Sprintf("%s\\r\\n%s\\r\\n%d", path, apiKey, timestamp)
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

func post(path, apiKey string, body interface{}) (json.RawMessage, error) {
	timestamp := time.Now().UnixMilli()
	signature := calculateSignature(path, apiKey, timestamp)

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", BaseUrl+path, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("signature", signature)
	req.Header.Set("token", apiKey)
	req.Header.Set("timestamp", fmt.Sprintf("%d", timestamp))
	req.Header.Set("lang", "en")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("FoxESS API %s: HTTP %d", path, resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp ApiResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, err
	}

	if apiResp.Errno != 0 {
		return nil, fmt.Errorf("FoxESS API %s: errno %d", path, apiResp.Errno)
	}

	return apiResp.Result, nil
}

func GetDevices(apiKey string) ([]Inverter, error) {
	var results []Inverter
	page := 0
	for {
		page++
		res, err := post("/op/v0/device/list", apiKey, map[string]int{
			"currentPage": page,
			"pageSize":    100,
		})
		if err != nil {
			return nil, err
		}

		var paginated struct {
			Total int        `json:"total"`
			Data  []Inverter `json:"data"`
		}
		if err := json.Unmarshal(res, &paginated); err != nil {
			return nil, err
		}

		results = append(results, paginated.Data...)
		if len(results) >= paginated.Total {
			break
		}
	}
	return results, nil
}

func GetRealTimeData(apiKey string, deviceSNs []string, variables []string) ([]RealTimeData, error) {
	req := map[string]interface{}{}
	if len(deviceSNs) > 0 {
		req["sn"] = deviceSNs[0]
	}
	if len(variables) > 0 {
		req["variables"] = variables
	}

	res, err := post("/op/v0/device/real/query", apiKey, req)
	if err != nil {
		return nil, err
	}

	var data []RealTimeData
	if err := json.Unmarshal(res, &data); err != nil {
		return nil, err
	}
	return data, nil
}
