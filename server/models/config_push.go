package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"server/modules/setting"
	"strings"
)

func PushConfig(env *Environment, config string, serverId string) string {
	downloadUrl := fmt.Sprintf("%s/v1/config/%v/json/%s", setting.AppURL, env.Id.Hex(), config)

	server, err := GetServerById(serverId)
	if err != nil {
		return "Server not found"
	}

	var results []string

	urls := strings.Split(server.Url, "\n")

	for _, url := range urls {
		url = strings.TrimSpace(url)
		if url == "" {
			continue
		}

		if err := pushServer(url+"/config/update", env.Name, downloadUrl); err != nil {
			logrus.Warningf("Push to %v:%v fail: %v", server.Name, url, err)
			results = append(results, fmt.Sprintf("%v > FAIL", server.Name))
		} else {
			results = append(results, fmt.Sprintf("%v > SUCCESS", server.Name))
		}
	}

	return strings.Join(results, ";")
}

func pushServer(serverUrl, environment, downloadUrl string) error {
	logrus.Infof("pushServer to %v, env: %v, downloadUrl: %v", serverUrl, environment, downloadUrl)

	var req struct {
		Environment string `json:"environment,omitempty"`
		DownloadUrl string `json:"download_url,omitempty"`
	}

	req.Environment = environment
	req.DownloadUrl = downloadUrl

	body := &bytes.Buffer{}
	reader := json.NewEncoder(body)

	if err := reader.Encode(req); err != nil {
		return err
	}

	res, err := http.Post(serverUrl, "application/json", body)
	if err != nil {
		return err
	}

	buf, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Request StatusCode: %v, ErrorInfo: %s", res.Status, buf)
	}

	return nil
}
