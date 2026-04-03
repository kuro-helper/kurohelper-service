package vndb

import (
	"encoding/json"
)

type Stats struct {
	Chars     int `json:"chars"`
	Producers int `json:"producers"`
	Releases  int `json:"releases"`
	Staff     int `json:"staff"`
	Tags      int `json:"tags"`
	Traits    int `json:"traits"`
	VN        int `json:"vn"`
}

// 取得VNDB統計資料
func GetStats() (*Stats, error) {
	r, err := sendGetRequest("/stats")
	if err != nil {
		return nil, err
	}

	var res Stats
	err = json.Unmarshal(r, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
