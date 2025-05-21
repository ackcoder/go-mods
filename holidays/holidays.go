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

const SearchApi = "https://sousuo.www.gov.cn/search-gov/data"

var (
	// 段落内容拆分
	contentRegexp = regexp.MustCompile(`^(.*?)放假(.*?)$`)
	// 匹配调休工作安排
	workPreRegexp = regexp.MustCompile(`([^。]*)(上班|补休).*?。`)
	// 匹配 "x年y月z日" 或 "x月y日" 或 "x日"
	dateRegexp = regexp.MustCompile(`(?:(\d{4})年)?(?:(\d{1,2})月)?(\d{1,2})日`)
)

type Worker struct {
	year int
	// 结果缓存
	cacheData map[int]Result
	// 自定义公文url获取函数
	pickFunc func(year int) (policyUrl string, err error)
}

// 创建工作实例
//   - {year} 指定年份, 例如"2023"
func NewWorker(year int) *Worker {
	w := new(Worker)
	w.cacheData = make(map[int]Result)
	if err := w.SetYear(year); err != nil {
		panic(err)
	}
	return w
}

// SetYear 设置年份
//   - {year} 指定年份, 例如"2023"
func (w *Worker) SetYear(year int) (err error) {
	if year < 2000 || year > time.Now().Year() {
		err = fmt.Errorf("无效年份: %d, 过旧(2000)或过新", year)
		return
	}
	w.year = year
	return
}

// SetPickPolicyDocumentUrlFunc 设置节假日政策公文Url获取的函数
//
// 注: 因政策查询api接口可能变化、所以留出可自定义的设置
func (w *Worker) SetPickPolicyDocumentUrlFunc(f func(year int) (policyUrl string, err error)) {
	w.pickFunc = f
}

// PickPolicyDocumentUrl 节假日政策公文Url获取
func (w *Worker) PickPolicyDocumentUrl() (policyUrl string, err error) {
	if w.pickFunc != nil {
		return w.pickFunc(w.year)
	}

	u, err := url.Parse(SearchApi)
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
	qry.Set("q", fmt.Sprintf("%d 节假日", w.year))
	u.RawQuery = qry.Encode()

	resultData, err := httpGet(u.String())
	if err != nil {
		return
	}
	// fmt.Println("原始结果数据", string(resultData))

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("接口数据解析异常、或无有效数据: %v", r)
		}
	}()

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
	if err = json.Unmarshal(resultData, &result); err != nil {
		return
	}
	if result.Code != 200 {
		err = fmt.Errorf("接口请求异常: [%d] %s", result.Code, result.Msg)
		return
	}
	policyUrl = result.SearchVO.ListVO[0].Url
	return
}

// SearchPolicyDocument 搜索政策公文、得到有效内容段落
//   - {policyUrl} 节假日公文Url
func (w *Worker) SearchPolicyDocument(policyUrl string) (paragraphs []string, err error) {
	htmlData, err := httpGet(policyUrl)
	if err != nil {
		return
	}

	doc, err := html.Parse(bytes.NewReader(htmlData))
	if err != nil {
		return
	}

	var originalParagraphs []string
	traverseParagraphs(doc, &originalParagraphs)

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

// ParseParagraphs 解析段落内容、得到最终结果
//   - {paragraphs} 段落内容列表
func (w *Worker) ParseParagraphs(paragraphs []string) (res Result, err error) {
	res = Result{}
	for _, paragraph := range paragraphs {
		paragraphMatches := contentRegexp.FindStringSubmatch(paragraph)
		if len(paragraphMatches) != 3 || len(paragraphMatches[1]) == 0 || len(paragraphMatches[2]) == 0 {
			continue
		}

		// 处理节假日
		holidayMatches := dateRegexp.FindAllStringSubmatch(paragraphMatches[1], -1)
		if len(holidayMatches) == 2 {
			// 展开日期范围
			sdStr := holidayMatches[0]
			sdy := w.year
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
				res.Holidays = append(res.Holidays, sdT.Format("2006-01-02"))
				sdT = sdT.AddDate(0, 0, 1)
			}
		} else {
			// 逐个格式化日期
			for _, holidayMatche := range holidayMatches {
				dy := w.year
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
				res.Holidays = append(res.Holidays, fmt.Sprintf("%04d-%s-%s", dy, dm, dd))
			}
		}

		// 处理调休工作日
		workdayPreMatches := workPreRegexp.FindAllStringSubmatch(paragraphMatches[2], -1)
		if len(workdayPreMatches) != 1 {
			continue
		}
		workdayMatches := dateRegexp.FindAllStringSubmatch(workdayPreMatches[0][0], -1)
		for _, workdayMatche := range workdayMatches {
			dy := w.year
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
			res.Workdays = append(res.Workdays, fmt.Sprintf("%04d-%s-%s", dy, dm, dd))
		}
	}
	if len(res.Holidays) != 0 {
		w.cacheData[w.year] = res
	}
	return
}

// QueryCache 查询缓存
func (w *Worker) QueryCache() (res Result, ok bool) {
	res, ok = w.cacheData[w.year]
	return
}

// 快速时间解析、替代time.Parse()
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

// 递归文档得到所有段落数据
func traverseParagraphs(n *html.Node, p *[]string) {
	if n.Type == html.ElementNode && n.Data == "p" {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type != html.TextNode {
				continue
			}
			pd := strings.TrimSpace(c.Data)
			pd = strings.ReplaceAll(pd, `\n`, "")
			if pd != "" {
				*p = append(*p, pd)
			}
		}
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		traverseParagraphs(child, p)
	}
}

func httpGet(url string) (resp []byte, err error) {
	respData, err := http.Get(url)
	if err != nil {
		return
	}
	defer respData.Body.Close()
	resp, err = io.ReadAll(respData.Body)
	return
}
