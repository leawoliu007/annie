package extractors

import (
	"github.com/leawoliu007/annie/extractors/streamtape"
	"net/url"
	"strings"

	"github.com/leawoliu007/annie/extractors/bilibili"
	"github.com/leawoliu007/annie/extractors/types"
	"github.com/leawoliu007/annie/extractors/udn"
	"github.com/leawoliu007/annie/extractors/universal"
	"github.com/leawoliu007/annie/utils"
)

var extractorMap map[string]types.Extractor

func init() {
	stExtractor := streamtape.New()

	extractorMap = map[string]types.Extractor{
		"": universal.New(), // universal extractor

		"bilibili":   bilibili.New(),
		"udn":        udn.New(),
		"streamtape": stExtractor,
		"streamta":   stExtractor, // streamta.pe
	}
}

// Extract is the main function to extract the data.
func Extract(u string, option types.Options) ([]*types.Data, error) {
	u = strings.TrimSpace(u)
	var domain string

	bilibiliShortLink := utils.MatchOneOf(u, `^(av|BV|ep)\w+`)
	if len(bilibiliShortLink) > 1 {
		bilibiliURL := map[string]string{
			"av": "https://www.bilibili.com/video/",
			"BV": "https://www.bilibili.com/video/",
			"ep": "https://www.bilibili.com/bangumi/play/",
		}
		domain = "bilibili"
		u = bilibiliURL[bilibiliShortLink[1]] + u
	} else {
		u, err := url.ParseRequestURI(u)
		if err != nil {
			return nil, err
		}
		if u.Host == "haokan.baidu.com" {
			domain = "haokan"
		} else {
			domain = utils.Domain(u.Host)
		}
	}
	extractor := extractorMap[domain]
	if extractor == nil {
		extractor = extractorMap[""]
	}
	videos, err := extractor.Extract(u, option)
	if err != nil {
		return nil, err
	}
	for _, v := range videos {
		v.FillUpStreamsData()
	}
	return videos, nil
}
