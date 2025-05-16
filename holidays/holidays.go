package holidays

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// TODO:
// 1. 改为结构体方法、通过New()创建实例调用下面三个公开方法；或是额外暴露一个公开便捷获取holidays与workdays方法
// 2. 添加单元测试、基准测试
// 3. 提升正则变量实例到顶层、避免每次调用时都创建
// 4. 错误检查和边界校验。例如年份不能是未来的、或是过旧的（比如2000年以前）
// 5. 更详细的文档（入参/返回/示例等；要么写外面README.md、要么写到方法头部注释里）
// 6. ParseParagraphs() 中对多个段落的处理、可以使用 goroutine + waitgroup
// 7. 节假日、调休日的结果缓存（应该在完成第1点前提下，写到结构体里）
// 8. http请求方法、需要超时控制、User-Agent、context支持等

const GovApi = "https://sousuo.www.gov.cn/search-gov/data"

// 获取政府发布的节假日公文URL
//   - {year} 指定年份
func GetHolidayPolicyDocumentURL(year int) (policyUrl string, err error) {
	u, err := url.Parse(GovApi)
	if err != nil {
		return
	}
	qry := u.Query()
	qry.Set("t", "zhengcelibrary_gw")
	qry.Set("p", "1")
	qry.Set("n", "5")
	qry.Set("sort", "score")
	qry.Set("sortType", "1")
	qry.Set("searchfield", "title")
	qry.Set("pcodeJiguan", "国办发明电")
	qry.Set("q", fmt.Sprintf("%d 节假日", year))
	u.RawQuery = qry.Encode()

	origRes, err := httpGet(u.String())
	if err != nil {
		return
	}
	// fmt.Println("原始结果数据", string(origRes))

	var result struct {
		Code     int    `json:"code"`
		Msg      string `json:"msg"`
		SearchVO struct {
			ListVO []struct {
				Pcode   string `json:"pcode"`      //文件号
				Title   string `json:"title"`      //文件标题
				Pubtime string `json:"pubtimeStr"` //发布时间
				Url     string `json:"url"`        //文件访问地址
			} `json:"listVO"`
		} `json:"searchVO"`
	}
	if err = json.Unmarshal(origRes, &result); err != nil {
		return
	}
	if result.Code != 200 {
		err = fmt.Errorf("接口请求异常: [%d] %s", result.Code, result.Msg)
		return
	}
	policyUrl = result.SearchVO.ListVO[0].Url

	return
}

// 搜索政策文件、得到有效内容段落
//   - {policyUrl} 通知公文访问地址
func SearchPolicyDocument(policyUrl string) (paragraphs []string, err error) {
	htmlData, err := httpGet(policyUrl)
	if err != nil {
		return
	}

	doc, err := html.Parse(bytes.NewReader(htmlData))
	if err != nil {
		return
	}

	var traverseFunc func(n *html.Node, p *[]string)
	traverseFunc = func(n *html.Node, p *[]string) {
		if n.Type == html.ElementNode && n.Data == "p" {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type != html.TextNode {
					continue
				}
				pd := strings.TrimSpace(c.Data)
				pd = strings.ReplaceAll(pd, `\n`, "")
				if pd != "" {
					*p = append(*p, pd) //确保添加的段落都是有内容的
				}
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			traverseFunc(child, p)
		}
	}
	var originalParagraphs []string
	traverseFunc(doc, &originalParagraphs) //遍历文档获取段落数据

	var capturing bool
	for _, p := range originalParagraphs {
		if !capturing {
			if strings.Contains(p, "通知如下") {
				capturing = true
			}
			continue
		}
		if strings.Contains(p, "节假日期间") {
			break
		}
		paragraphs = append(paragraphs, p) //只过滤获取有效的段落内容
	}
	return
}

// 解析段落内容、得到 节假日,调休工作日 列表
//   - {year} 指定年份、若段落内容中有写年份则以其为优先
//   - {paragraphs} 段落内容列表
func ParseParagraphs(year int, paragraphs []string) (holidays, workdays []string, err error) {
	// 段落内容拆分
	contentRe := regexp.MustCompile(`^(.*?)放假(.*?)$`)
	// 匹配调休工作安排
	workPreRe := regexp.MustCompile(`([^。]*)(上班|补休).*?。`)
	// 匹配 "x年y月z日" 或 "x月y日" 或 "x日"
	dateRe := regexp.MustCompile(`(?:(\d{4})年)?(?:(\d{1,2})月)?(\d{1,2})日`)

	for _, paragraph := range paragraphs {
		paragraphMatches := contentRe.FindStringSubmatch(paragraph)
		if len(paragraphMatches) < 3 {
			continue
		}

		// 处理节假日
		holidayMatches := dateRe.FindAllStringSubmatch(paragraphMatches[1], -1)
		if len(holidayMatches) == 2 {
			// 展开日期范围
			sdStr := holidayMatches[0]
			sdy := year
			if len(sdStr[1]) > 0 {
				sdy, err = strconv.Atoi(sdStr[1])
				if err != nil {
					return
				}
			}
			sdm := sdStr[2]
			if len(sdm) == 1 {
				sdm = "0" + sdm
			}
			sdd := sdStr[3]
			if len(sdd) == 1 {
				sdd = "0" + sdd
			}
			sd := fmt.Sprintf("%04d-%s-%s", sdy, sdm, sdd)

			edStr := holidayMatches[1]
			edy := sdy
			if len(edStr[1]) > 0 {
				edy, err = strconv.Atoi(edStr[1])
				if err != nil {
					return
				}
			}
			edm := edStr[2]
			if len(edm) == 0 {
				edm = sdm
			} else if len(edm) == 1 {
				edm = "0" + edm
			}
			edd := edStr[3]
			if len(edd) == 1 {
				edd = "0" + edd
			}
			ed := fmt.Sprintf("%04d-%s-%s", edy, edm, edd)

			var sdT, edT time.Time
			// sdT, err = time.Parse("2006-01-02", sd)
			sdT, err = fastTimeParse(sd)
			if err != nil {
				return
			}
			// edT, err = time.Parse("2006-01-02", ed)
			edT, err = fastTimeParse(ed)
			if err != nil {
				return
			}
			for !sdT.After(edT) {
				holidays = append(holidays, sdT.Format("2006-01-02"))
				sdT = sdT.AddDate(0, 0, 1)
			}
		} else {
			// 逐个格式化日期
			for _, holidayMatche := range holidayMatches {
				dy := year
				if len(holidayMatche[1]) > 0 {
					dy, err = strconv.Atoi(holidayMatche[1])
					if err != nil {
						return
					}
				}
				dm := holidayMatche[2]
				if len(dm) == 0 {
					continue //不满足"x月x日"的不加入
				}
				if len(dm) == 1 {
					dm = "0" + dm
				}
				dd := holidayMatche[3]
				if len(dd) == 1 {
					dd = "0" + dd
				}
				holidays = append(holidays, fmt.Sprintf("%04d-%s-%s", dy, dm, dd))
			}
		}

		// 处理调休/工作日
		workdayPreMatches := workPreRe.FindAllStringSubmatch(paragraphMatches[2], -1)
		if len(workdayPreMatches) != 1 {
			continue
		}
		workdayMatches := dateRe.FindAllStringSubmatch(workdayPreMatches[0][0], -1)
		for _, workdayMatche := range workdayMatches {
			dy := year
			if len(workdayMatche[1]) > 0 {
				dy, err = strconv.Atoi(workdayMatche[1])
				if err != nil {
					return
				}
			}
			dm := workdayMatche[2]
			if len(dm) == 0 {
				continue //不满足"x月x日"的不加入
			}
			if len(dm) == 1 {
				dm = "0" + dm
			}
			dd := workdayMatche[3]
			if len(dd) == 1 {
				dd = "0" + dd
			}
			workdays = append(workdays, fmt.Sprintf("%04d-%s-%s", dy, dm, dd))
		}
	}
	return
}

func fastTimeParse(s string) (time.Time, error) {
	if len(s) != 10 || s[4] != '-' || s[7] != '-' {
		return time.Time{}, fmt.Errorf("invalid date format: %s", s)
	}

	year, err := strconv.Atoi(s[0:4])
	if err != nil {
		return time.Time{}, err
	}

	month, err := strconv.Atoi(s[5:7])
	if err != nil {
		return time.Time{}, err
	}

	day, err := strconv.Atoi(s[8:10])
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local), nil
}

// TODO: 后续考虑用包内的 http-req 组件替代
func httpGet(url string) (resp []byte, err error) {
	respData, err := http.Get(url)
	if err != nil {
		return
	}
	defer respData.Body.Close()
	resp, err = io.ReadAll(respData.Body)
	return
}
