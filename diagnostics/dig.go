package diagnostics

import (
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

// AnswerSection of diginfo
type AnswerSection struct {
	Domain          string `json:"domain"`
	TTL             int    `json:"ttl"`
	RecordClass     string `json:"recordClass"`
	RecordType      string `json:"recordType"`
	PreferenceValue string `json:"preferenceValue"`
	Value           string `json:"value"`
}

// DigInfo in response
type DigInfo struct {
	Hostname         string           `json:"hostname"`
	QuertyType       string           `json:"queryType"`
	AnswerSection    []*AnswerSection `json:"answerSection"`
	AuthoritySection string           `json:"authoritySection"`
	//Result           interface{}      `json:"result"`
}

// DigResponse with diginfo
type DigResponse struct {
	DigInfo *DigInfo `json:"digInfo"`
}

// GetDigInfo retrieves retrieves dig results for a given location and hostname
func GetDigInfo(location string, hostName string) (*DigResponse, error) {
	stat := &DigResponse{}
	hostURL := fmt.Sprintf("/diagnostic-tools/v2/ghost-locations/%s/dig-info?hostName=%s&queryType=A", location, hostName)

	req, err := client.NewRequest(
		Config,
		"GET",
		hostURL,
		nil,
	)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)
	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}
	err = client.BodyJSON(res, stat)
	if err != nil {
		return nil, err
	}
	return stat, nil
}
