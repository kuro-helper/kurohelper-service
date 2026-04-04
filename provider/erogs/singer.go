package erogs

import (
	"encoding/json"
	"fmt"
	"strings"

	"kurohelperservice"
)

/*
 * 獲取指定歌手的相關情報
 * creator的特定Shubetu處理(Shubetu = 6)
 */

type Singer struct {
	SingerID   int    `json:"singer_id"`
	SingerName string `json:"singer_name"`
	Twitter    string `json:"twitter"`
	Blog       string `json:"blog"`
	Pixiv      string `json:"pixiv"`
	MusicInfo  []struct {
		MusicID        int     `json:"music_id"`
		MusicName      string  `json:"musicname"`
		ReleaseDate    string  `json:"releasedate"`
		MusicAvgScore  float64 `json:"music_avg_score"`
		MusicVoteCount int     `json:"music_vote_count"`
		GameName       string  `json:"game_name"`
		DMM            string  `json:"dmm"` // 遊戲的 DMM 圖片網址
	} `json:"music_info"` // 音樂作品清單
}

// Use kewords search singer list data
//
// search table is the same as Creator (use createrlist)
func SearchSingerListByKeyword(keywords []string) ([]CreatorList, error) {
	if len(keywords) == 0 {
		return nil, kurohelperservice.ErrSearchNoContent
	}

	// pre-build keySQL
	keySQL := "WHERE "
	var keywordSQLList []string
	for _, k := range keywords {
		formatK := buildSearchStringSQL(k)
		if strings.TrimSpace(formatK) != "" {
			keywordSQLList = append(keywordSQLList, fmt.Sprintf("cr.name ILIKE '%s'", formatK))
		}
	}

	keySQL += strings.Join(keywordSQLList, " OR ")
	keySQL += " AND EXISTS ( SELECT 1 FROM shokushu s WHERE s.creater = cr.id AND s.shubetu = 6)"

	sql := buildCreatorListSQL(keySQL)

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}

	var res []CreatorList
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Use erogs id search single singer data
func SearchSingerByKeyword(id int) (*Singer, error) {
	sql := buildSingerSQL(fmt.Sprintf("WHERE cr.id = '%d'", id))

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}

	var res Singer
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// build search singer sql
// Arguments:
//   - keySQL: A pre-constructed SQL WHERE-clause fragment.
func buildSingerSQL(keySQL string) string {
	return fmt.Sprintf(`
SELECT
  row_to_json(t)
FROM (
  SELECT
    cr.id AS singer_id,
    cr.name AS singer_name,
    cr.twitter_username AS twitter,
    cr.blog,
    cr.pixiv,
    (
      SELECT
        json_agg(m_data)
      FROM (
        SELECT
          m.id AS music_id,
          m.name AS musicname,
          m.releasedate,
          score_data.avg_tokuten AS music_avg_score,
          score_data.tokuten_count AS music_vote_count,
          STRING_AGG(DISTINCT g.gamename, ', ') AS game_name,
          (ARRAY_AGG(g.dmm ORDER BY g.id))[1] AS dmm
        FROM musiclist m
        INNER JOIN singer s 
          ON s.music = m.id 
          AND s.creater = cr.id
        INNER JOIN (
          SELECT
            music,
            ROUND(AVG(LEAST(tokuten, 100))::numeric, 2) AS avg_tokuten,
            COUNT(tokuten) AS tokuten_count
          FROM usermusic_tokuten
          GROUP BY music
        ) AS score_data 
          ON score_data.music = m.id
        INNER JOIN game_music gm ON gm.music = m.id
        INNER JOIN gamelist g ON g.id = gm.game
        GROUP BY 
          m.id, 
          m.name, 
          m.releasedate, 
          score_data.avg_tokuten, 
          score_data.tokuten_count
        ORDER BY score_data.avg_tokuten DESC, m.releasedate DESC
      ) AS m_data
    ) AS music_info
  FROM createrlist cr
  %s
) t;
`, keySQL)
}
