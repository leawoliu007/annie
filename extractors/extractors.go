package extractors

import (
	"github.com/leawoliu007/annie/extractors/streamtape"
	"net/url"
	"strings"

	"github.com/leawoliu007/annie/extractors/acfun"
	"github.com/leawoliu007/annie/extractors/bcy"
	"github.com/leawoliu007/annie/extractors/bilibili"
	"github.com/leawoliu007/annie/extractors/douyin"
	"github.com/leawoliu007/annie/extractors/douyu"
	"github.com/leawoliu007/annie/extractors/eporner"
	"github.com/leawoliu007/annie/extractors/facebook"
	"github.com/leawoliu007/annie/extractors/geekbang"
	"github.com/leawoliu007/annie/extractors/haokan"
	"github.com/leawoliu007/annie/extractors/instagram"
	"github.com/leawoliu007/annie/extractors/iqiyi"
	"github.com/leawoliu007/annie/extractors/mgtv"
	"github.com/leawoliu007/annie/extractors/miaopai"
	"github.com/leawoliu007/annie/extractors/netease"
	"github.com/leawoliu007/annie/extractors/pixivision"
	"github.com/leawoliu007/annie/extractors/pornhub"
	"github.com/leawoliu007/annie/extractors/qq"
	"github.com/leawoliu007/annie/extractors/tangdou"
	"github.com/leawoliu007/annie/extractors/tiktok"
	"github.com/leawoliu007/annie/extractors/tumblr"
	"github.com/leawoliu007/annie/extractors/twitter"
	"github.com/leawoliu007/annie/extractors/types"
	"github.com/leawoliu007/annie/extractors/udn"
	"github.com/leawoliu007/annie/extractors/universal"
	"github.com/leawoliu007/annie/extractors/vimeo"
	"github.com/leawoliu007/annie/extractors/weibo"
	"github.com/leawoliu007/annie/extractors/xvideos"
	"github.com/leawoliu007/annie/extractors/yinyuetai"
	"github.com/leawoliu007/annie/extractors/youku"
	"github.com/leawoliu007/annie/extractors/youtube"
	"github.com/leawoliu007/annie/utils"
)

var extractorMap map[string]types.Extractor

func init() {
	douyinExtractor := douyin.New()
	youtubeExtractor := youtube.New()
	stExtractor := streamtape.New()

	extractorMap = map[string]types.Extractor{
		"": universal.New(), // universal extractor

		"douyin":     douyinExtractor,
		"iesdouyin":  douyinExtractor,
		"bilibili":   bilibili.New(),
		"bcy":        bcy.New(),
		"pixivision": pixivision.New(),
		"youku":      youku.New(),
		"youtube":    youtubeExtractor,
		"youtu":      youtubeExtractor, // youtu.be
		"iqiyi":      iqiyi.New(iqiyi.SiteTypeIqiyi),
		"iq":         iqiyi.New(iqiyi.SiteTypeIQ),
		"mgtv":       mgtv.New(),
		"tangdou":    tangdou.New(),
		"tumblr":     tumblr.New(),
		"vimeo":      vimeo.New(),
		"facebook":   facebook.New(),
		"douyu":      douyu.New(),
		"miaopai":    miaopai.New(),
		"163":        netease.New(),
		"weibo":      weibo.New(),
		"instagram":  instagram.New(),
		"twitter":    twitter.New(),
		"qq":         qq.New(),
		"yinyuetai":  yinyuetai.New(),
		"geekbang":   geekbang.New(),
		"pornhub":    pornhub.New(),
		"xvideos":    xvideos.New(),
		"udn":        udn.New(),
		"tiktok":     tiktok.New(),
		"haokan":     haokan.New(),
		"acfun":      acfun.New(),
		"eporner":    eporner.New(),
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
