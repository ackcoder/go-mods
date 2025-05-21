package holidays

type Result struct {
	Holidays []string `json:"holidays" dc:"节假日"`
	Workdays []string `json:"workdays" dc:"调休工作日"`
}

var localWorker *Worker

func init() {
	localWorker = NewWorker(2023)
}

// Get 获取指定年份的节假日、调休工作日
//   - {year} 指定年份, 如"2023"
func Get(year int) (res Result, err error) {
	if err = localWorker.SetYear(year); err != nil {
		return
	}
	if res, ok := localWorker.QueryCache(); ok {
		return res, nil
	}
	url, err := localWorker.PickPolicyDocumentUrl()
	if err != nil {
		return
	}
	ps, err := localWorker.SearchPolicyDocument(url)
	if err != nil {
		return
	}
	return localWorker.ParseParagraphs(ps)
}
