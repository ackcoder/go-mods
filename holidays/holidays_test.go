package holidays_test

import (
	"testing"
	"time"

	"github.com/ackcoder/go-mods/holidays"
)

func TestGet(t *testing.T) {
	_, err := holidays.Get(-1)
	if err != nil {
		t.Error(err)
	}

	res, err := holidays.Get(2024)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v", res)

	time.Sleep(time.Second)

	// repeat
	res, err = holidays.Get(2024)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v", res)
}

func TestWorker_StepByStep(t *testing.T) {
	w := holidays.NewWorker(2025)

	// t.Run("测试获取公文链接", func(t *testing.T) {
	// 	url, err := w.PickPolicyDocumentUrl()
	// 	if err != nil {
	// 		t.Error(err)
	// 	}
	// 	t.Log(url)
	// })

	// t.Run("测试搜索指定节假日公文", func(t *testing.T) {
	// 	ps, err := w.SearchPolicyDocument("https://www.gov.cn/zhengce/zhengceku/202411/content_6986383.htm")
	// 	if err != nil {
	// 		t.Error(err)
	// 	}
	// 	t.Log(ps)
	// })

	t.Run("测试解析指定段落", func(t *testing.T) {
		ps := []string{
			// // 正常有效内容段落
			// "1月1日（周三）放假1天，不调休。",
			// "1月28日（农历除夕、周二）至2月4日（农历正月初七、周二）放假调休，共8天。1月26日（周日）、2月8日（周六）上班。",
			// "4月4日（周五）至6日（周日）放假，共3天。",
			// "5月1日（周四）至5日（周一）放假调休，共5天。4月27日（周日）上班。",
			// "5月31日（周六）至6月2日（周一）放假，共3天。",
			// "10月1日（周三）至8日（周三）放假调休，共8天。9月28日（周日）、10月11日（周六）上班。",

			// // 跨月内容
			// "9月29日至10月6日放假调休，共8天。10月7日（星期六）、10月8日（星期日）上班。",

			// 特殊内容
			"2012年1月1日至3日放假调休，共3天。2011年12月31日（星期六）上班。",
			"1月25日至31日放假，共7天。其中，1月25日（星期日、农历除夕）、1月26日（星期一、农历正月初一）、1月27日（星期二、农历正月初二）为法定节假日，1月31日（星期六）照常公休；1月25日（星期日）公休日调至1月28日（星期三），1月24日（星期六）、2月1日（星期日）两个公休日调至1月29日（星期四）、1月30日（星期五）。1月24日（星期六）、2月1日（星期日）上班。",
		}
		res, err := w.ParseParagraphs(ps)
		if err != nil {
			t.Error(err)
		}
		t.Logf("%#v", res)
	})
}

func TestWorker(t *testing.T) {
	w := holidays.NewWorker(2021)

	url, err := w.PickPolicyDocumentUrl()
	if err != nil {
		t.Error(err)
	}
	t.Log(url)
	ps, err := w.SearchPolicyDocument(url)
	if err != nil {
		t.Error(err)
	}
	t.Log(ps)
	res, err := w.ParseParagraphs(ps)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v", res)
}
