package vndb

import (
	"encoding/json"
	"kurohelperservice"
	"strings"
)

// staff Response
//
// 統一字串不使用指標
type StaffSearchResponse struct {
	ID          string               `json:"id"`          // vndbid
	AID         int                  `json:"aid"`         // alias id
	IsMain      bool                 `json:"ismain"`      // 是否是主要名字
	Name        string               `json:"name"`        // 羅馬拼音名字
	Original    string               `json:"original"`    // 原文名, 可能為 null
	Lang        string               `json:"lang"`        // 主要語言
	Gender      string               `json:"gender"`      // 性別, 可能為 null
	Description string               `json:"description"` // 可能有格式化代碼
	ExtLinks    []ExtlinksResponse   `json:"extlinks"`    // 外部連結
	Aliases     []StaffAliasResponse `json:"aliases"`     // 別名清單
}

func GetStaffByFuzzy(keyword string, roleType string) (*BasicResponse[StaffSearchResponse], error) {
	req := VndbCreate()

	filters := []any{}
	if roleType != "" {
		filters = append(filters, "and")
		// 傳進來的直接就是API篩選項規格的字串
		filters = append(filters, []string{"type", "=", roleType})
		filters = append(filters, []string{"search", "=", keyword})
	} else {
		filters = []any{"search", "=", keyword}
	}

	req.Filters = filters

	basicFields := "id, aid, ismain, name, original, lang, gender, description"
	extlinksFields := "extlinks{url, label, name, id}"
	aliasesFields := "aliases{aid, name, latin, ismain}"

	allFields := []string{
		basicFields,
		extlinksFields,
		aliasesFields,
	}

	req.Fields = strings.Join(allFields, ", ")

	jsonStaff, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	r, err := sendPostRequest("/staff", jsonStaff)
	if err != nil {
		return nil, err
	}

	var res BasicResponse[StaffSearchResponse]
	err = json.Unmarshal(r, &res)
	if err != nil {
		return nil, err
	}

	if len(res.Results) == 0 {
		return nil, kurohelperservice.ErrSearchNoContent
	}

	return &res, nil
}
