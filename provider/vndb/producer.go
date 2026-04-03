package vndb

import (
	"encoding/json"
	"kurohelperservice"
	"strings"
)

// 查詢品牌API
type (
	// producer Response
	ProducerSearchResponse struct {
		Producer BasicResponse[ProducerSearchProducerResponse]
		Vn       BasicResponse[ProducerSearchVnResponse]
	}

	// 品牌結構
	ProducerSearchProducerResponse struct {
		ID          string             `json:"id"`
		Name        string             `json:"name"`
		Original    string             `json:"original"` // *string
		Aliases     []string           `json:"aliases"`
		Lang        string             `json:"lang"`
		Type        string             `json:"type"`
		Description string             `json:"description"` // *string
		Extlinks    []ExtlinksResponse `json:"extlinks"`
	}

	// 遊戲結構
	ProducerSearchVnResponse struct {
		ID            string  `json:"id"` // vndbid
		Title         string  `json:"title"`
		Alttitle      string  `json:"alttitle"`
		Released      *string `json:"released"` // 發售日期，因為vndb不是回傳標準格式，用字串儲存
		Average       float64 `json:"average"`
		Rating        float64 `json:"rating"`
		Votecount     int     `json:"votecount"`
		LengthMinutes int     `json:"length_minutes"`
		LengthVotes   int     `json:"length_votes"`
		Image         Image   `json:"image"`
	}

	// 遊戲中的圖片結構(只取需要的)
	//
	// Sexual跟Violence官方文檔說明是整數，但實測有浮點數出現可能
	Image struct {
		Thumbnail string  `json:"thumbnail"`
		Sexual    float64 `json:"sexual"`
		Violence  float64 `json:"violence"`
	}
)

func GetProducerByFuzzy(keyword string, companyType string) (*ProducerSearchResponse, error) {
	reqProducer := VndbCreate()

	filtersProducer := []any{}
	if companyType != "" {
		filtersProducer = append(filtersProducer, "and")
		switch companyType {
		case "company":
			filtersProducer = append(filtersProducer, []string{"type", "=", "co"})
		case "individual":
			filtersProducer = append(filtersProducer, []string{"type", "=", "in"})
		case "amateur":
			filtersProducer = append(filtersProducer, []string{"type", "=", "ng"})
		}
		filtersProducer = append(filtersProducer, []string{"search", "=", keyword})
	} else {
		filtersProducer = []any{"search", "=", keyword}
	}

	reqProducer.Filters = filtersProducer

	basicFields := "id, name, original, aliases, lang, type, description"
	extlinksFields := "extlinks.url, extlinks.label, extlinks.name, extlinks.id"

	allFields := []string{
		basicFields,
		extlinksFields,
	}

	reqProducer.Fields = strings.Join(allFields, ", ")

	jsonProducer, err := json.Marshal(reqProducer)
	if err != nil {
		return nil, err
	}

	r, err := sendPostRequest("/producer", jsonProducer)
	if err != nil {
		return nil, err
	}

	var resProducer BasicResponse[ProducerSearchProducerResponse]
	err = json.Unmarshal(r, &resProducer)
	if err != nil {
		return nil, err
	}

	if len(resProducer.Results) == 0 {
		return nil, kurohelperservice.ErrSearchNoContent
	}

	// 等到查詢解析完後才能去查詢遊戲的資料
	reqVn := VndbCreate()

	reqVn.Filters = []any{
		"developer", "=", []any{"id", "=", resProducer.Results[0].ID},
	}

	reqVn.Fields = "id, title, alttitle, released, length_minutes, length_votes, average, rating, votecount, image.sexual, image.violence, image.votecount, image.thumbnail"

	reverse := true
	reqVn.Reverse = &reverse

	jsonVn, err := json.Marshal(reqVn)
	if err != nil {
		return nil, err
	}

	r, err = sendPostRequest("/vn", jsonVn)
	if err != nil {
		return nil, err
	}

	var resVn BasicResponse[ProducerSearchVnResponse]
	err = json.Unmarshal(r, &resVn)
	if err != nil {
		return nil, err
	}

	if len(resVn.Results) == 0 {
		return nil, kurohelperservice.ErrSearchNoContent
	}

	return &ProducerSearchResponse{
		Producer: resProducer,
		Vn:       resVn,
	}, nil
}
