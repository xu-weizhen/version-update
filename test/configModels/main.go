package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"time"
)

type DownloadStruct struct {
	download_url        string
	update_version_code string
	title               string
	update_tips         string
	md5                 string
}

const appID = "10000"

var updateVersionCodes []string

func main() {
	rand.Seed(time.Now().Unix())
	var file *os.File
	var err error
	fileName := "./versionFile.txt"
	if !checkFileIsExist(fileName) {
		if file, err = os.Create(fileName); err != nil {
			log.Fatal(err)
		}
	} else {
		if file, err = os.OpenFile(fileName, os.O_WRONLY, 0777); err != nil {
			log.Fatal(err)
		}
	}
	defer file.Close()
	keys := []string{
		"aid",
		"platform",
		"max_update_version_code",
		"min_update_version_code",
		"max_os_api",
		"min_os_api",
		"cpu_arch",
		"channel",
		"download_url",
		"update_version_code",
		"md5",
		"title",
		"update_tips",
		"white_list_name",
	}
	downloadUrlFunc := downloadUrlTitleFuncCreate(file)
	platformFunc := platformFuncCreate()
	osApiRangeFunc := osApiRangeFuncCreate()
	cpuArchFunc := cpuArchFuncCreate()
	channelFunc := channelFuncCreate()
	whiteListNameSelectFunc := whiteListNameSelectFuncCreate()
	var download DownloadStruct
	var osApiRange []int
	var updateVersionRange []string
	for {
		download = downloadUrlFunc()
		if download.download_url == "" {
			break
		}
		osApiRange = osApiRangeFunc()
		updateVersionRange = updateVersionRangeFunc()
		params := url.Values{
			"download_url":            {download.download_url},
			"update_version_code":     {download.update_version_code},
			"md5":                     {download.md5},
			"title":                   {download.title},
			"update_tips":             {download.update_tips},
			"platform":                {platformFunc()},
			"max_update_version_code": {updateVersionRange[1]},
			"min_update_version_code": {updateVersionRange[0]},
			"cpu_arch":                {cpuArchFunc()},
			"channel":                 {channelFunc()},
			"aid":                     {appID},
			"max_os_api":              {strconv.Itoa(osApiRange[1])},
			"min_os_api":              {strconv.Itoa(osApiRange[0])},
			"white_list_name":         {fmt.Sprintf("whiteList%d.txt", whiteListNameSelectFunc())},
		}
		for _, key := range keys {
			fmt.Println(params.Get(key))
		}
		fmt.Println("\n\n")
	}
}

// osApiFuncCreate 随机 os_api 对 0是min_os_api 1是max_os_api
func osApiRangeFuncCreate() func() []int {
	maxAPIRange := []int{18, 20}
	minAPIRange := []int{1, 20}
	return func() []int {
		result := make([]int, 2)
		result[0] = rand.Intn(minAPIRange[1]-minAPIRange[0]+1) + minAPIRange[0]
		result[1] = rand.Intn(maxAPIRange[1]-max(maxAPIRange[0], result[0])+1) + max(maxAPIRange[0], result[0])
		return result
	}
}

// cpuArchFuncCreate 随机 cpu_arch 字符串
func cpuArchFuncCreate() func() string {
	cpuStrings := []string{"32", "64"}
	return func() string {
		return cpuStrings[rand.Intn(2)]
	}
}

// channelFuncCreate 随机 channel
func channelFuncCreate() func() string {
	channelStrings := []string{"huawei", "xiaomi"}
	length := len(channelStrings)
	return func() string {
		return channelStrings[rand.Intn(length)]
	}
}

// platformFuncCreate 随机平台版本
func platformFuncCreate() func() string {
	platform := []string{"Android", "iOS"}
	return func() string {
		return platform[rand.Intn(2)]
	}
}

// downloadUrlTitleFuncCreate 创建下载链接
func downloadUrlTitleFuncCreate(file *os.File) func() DownloadStruct {
	downloadUrls := []DownloadStruct{
		{"https://www.runoob.com/html/html-links.html", "", "HTML 链接 | 菜鸟教程", "HTML 链接  HTML 使用超级链接与网络上的另一个文档相连。几乎可以在所有的网页中找到链接。点击链接可以从一张页面跳转到另一张页面。   尝试一下 - 实例   HTML 链接 如何在HTML文档中创建链接。   (可以在本页底端找到更多实例)  HTML 超链接（链接） HTML使用标签 &lt;a&gt;来设置超文本链接。 超链接可以是一个字，一个词，或者一组词，也可以是一幅图像，您可以点击这些内容来跳转到新的..", ""},
		{"https://baike.baidu.com/item/%E7%BD%91%E7%AB%99%E9%93%BE%E6%8E%A5/8916094", "", "网站链接_百度百科", "网站链接分为内链和外链。内链是指本站内部的链接，外链是指外部网站指向本站的链接。通俗的说，外链就是在你的网站（任何页面）建立一个或几个链接，点击后就可到达我的网站，同样在我的网站建立这样的链接，点击后就可以到你的网站。这样就可以让访客在两个网站之间流动。通常，要建立这样的链接需要通过相互协商，互利互惠才好。", ""},
		{"https://baike.baidu.com/item/%E8%B6%85%E9%93%BE%E6%8E%A5/97857", "", "超链接_百度百科", "超级链接简单来讲，就是指按内容链接。超级链接在本质上属于一个网页的一部分，它是一种允许我们同其他网页或站点之间进行连接的元素。各个网页链接在一起后，才能真正构成一个网站。所谓的超链接是指从一个网页指向一个目标的连接关系，这个目标可以是另一个网页，也可以是相同网页上的不同位置，还可以是一个图片，一个电子邮件地址，一个文件，甚至是一个应用程序。而在一个网页中用来超链接的对象，可以是一段文本或者是一个图片。当浏览者单击已经链接的文字或图片后，链接目标将显示在浏览器上，并且根据目标的类型来打开或运行。", ""},
		{"https://support.microsoft.com/zh-cn/office/%E5%88%9B%E5%BB%BA%E6%88%96%E7%BC%96%E8%BE%91%E8%B6%85%E9%93%BE%E6%8E%A5-5d8c0804-f998-4143-86b1-1199735e07bf", "", "创建或编辑超链接", "在文档中添加或编辑指向网站、本地文件、电子邮件或定位点的超链接。 ", ""},
		{"https://www.zhihu.com/question/267339149", "", "(14 条消息) 怎么批量提取一个网页里面的链接？ - 知乎", "最近在学习淘宝技术，想知道怎么批量提取出这种页面里所有宝贝的网页链接。如果一个一个复制的话，真的是…", ""},
		{"https://www.huaweicloud.com/theme/400105-1-R", "", "404页面-华为云", "", ""},
		{"https://m.300.cn/wenda/c777.html", "", "制作网页链接【互联网知识】– 中企动力移动端", "中企动力移动版（300.cn）互联网知识库 - 为您提供制作网页链接相关问答，帮助您解决各类制作网页链接疑问，便于解决您的疑虑。", ""},
		{"http://www.divcss5.com/html/h405.shtml", "", "从A网页链接转到B网页 - DIVCSS5", "利用html a锚文本，从一个页面链接另外页面。在html网页中使用 超链接标签。", ""},
		{"http://www.woshipm.com/pd/310637.html", "", "技术贴：一篇文章看懂链接（超链接）设计 | 人人都是产品经理", "产品设计时细节是产品经理最头疼的问题，一个button，一个链接都要考虑太多的细节问题。作者整理了常见的一些功能设计问题，一篇文章看懂这些功能设计。来学习吧。 定义 链接也称为超链接，所谓的超链接是指从一个网页指向一个目标的连接关系，这个目标可以是另一个网页，", ""},
		{"https://books.google.co.jp/books?id=HAnaA52Jl9MC&pg=PA218&lpg=PA218&dq=%E7%BD%91%E9%A1%B5%E9%93%BE%E6%8E%A5&source=bl&ots=vpmu3R2rqC&sig=ACfU3U0JXkjdVtoPje0u0yZEketla7iQPg&hl=zh-CN&sa=X&ved=2ahUKEwjBy-SEqJD0AhXUJKYKHeFqABoQ6AF6BAglEAM", "", "信息技术基础教程 - 王仲轩, 罗廷礼, 徐贤军 - Google 图书", "全书先对信息技术、计算机工作原理等基础理论进行了系统的介绍;然后对信息技术、计算机工作原理等基础理论进行了系统的介绍;最后,又对计算机网络、Internet、信息安全和多媒体分别进行了细致的讲述.", ""},
		{"https://support.apple.com/kb/PH23632?viewlocale=zh_CN&locale=zh_CN", "", "用于 Mac 的Pages: 在 Pages 文稿中链接到网页、电子邮件或页面 (中国)", "", ""},
		{"http://help.163.com/09/1217/16/5QOF2C7I00753VB8.html", "", "写信时如何插入网页链接 ？-163邮箱常见问题", "网易帮助中心(help.163.com),为网易用户提供在线的产品帮助,您还可以与网易助手小易即时对话解决使用网易产品时遇到的问题", ""},
		{"https://finance.eastmoney.com/a/202111112177963174.html", "", "或受20亿元定增计划影响 恒顺醋业股价大跌 _ 东方财富网", "【或受20亿元定增计划影响 恒顺醋业股价大跌】昨日宣布20亿元定增计划，今日股价跳水。11月11日早间开盘，恒顺醋业股价大幅下跌，跌幅一度超7%。截至午间休盘，其股价下跌5.17% ，报15.79元/股。", ""},
		{"https://support.google.com/webmasters/answer/9012289?hl=zh-Hans", "", "网址检查工具 - Search Console帮助", "网址检查工具简介网址检查工具可提供与特定网页被 Google 编入索引的版本有关的信息，其中包括 AMP 错误、结构化数据错误和索引编制问题。常见任", ""},
		{"https://help.fanruan.com/finebi/doc-view-1596.html", "", "跳转到网页链接- FineBI帮助文档 FineBI帮助文档", "1. 概述1.1 应用场景&nbsp;FineBI 支持点击组件后跳转到普通网页。场景一：用户想要点击不同字段值，跳转到不同的网络页面。如下图我在表格中点击不同的标题，就能跳转到对应的帮助文档页面。场景", ""},
		{"https://finance.sina.com.cn/tech/2021-11-10/doc-iktzqtyu6435416.shtml", "", "法官命令App Store在12月9日之前允许开发者添加外部支付链接|app store|开发者|苹果_手机_新浪科技_新浪网", "在今天早些时候的听证会后，YvonneGonzalezRogers法官拒绝了苹果公司关于推迟执行永久禁令的请求，该禁令将要求苹果公司对AppStore做出重大改变。作为苹果诉Epic诉讼案判决的一部分，Rogers要求苹果允许开发者添加外部", ""},
	}
	count := len(downloadUrls)
	updateVersionCodes = make([]string, 1, count+1)
	updateVersionCodes[0] = "0.0.0.0"
	updateVersionCodeFunc := updateVersionCodeFuncCreate(file)
	return func() DownloadStruct {
		count--
		if count <= 0 {
			return DownloadStruct{}
		}
		randPos := rand.Intn(count)
		downloadUrls[randPos], downloadUrls[count] = downloadUrls[count], downloadUrls[randPos]
		result := downloadUrls[count]
		result.update_version_code = updateVersionCodeFunc()
		t := result.download_url + result.update_tips + result.update_version_code + result.title
		result.md5 = fmt.Sprintf("%x", md5.Sum([]byte(t)))
		return result
	}
}

// updateVersionCodeFuncCreate 随机更新版本号的函数 按照概率
func updateVersionCodeFuncCreate(file *os.File) func() string {
	bits := make([]int, 4)
	bits[0] = 1
	const firstVersionUpdate = 0.025
	const secondVersionUpdate = 0.1
	const thirdVersionUpdate = 0.3
	if _, err := file.WriteString("0.0.0.0 "); err != nil {
		log.Fatal(err)
	}
	return func() string {
		result := fmt.Sprintf("%d.%d.%d.%d ", bits[0], bits[1], bits[2], bits[3])
		if _, err := file.WriteString(result); err != nil {
			log.Fatal(err)
		}
		updateVersionCodes = append(updateVersionCodes, result)
		randFloat := rand.Float32()
		// 下个版本号
		switch {
		case randFloat < firstVersionUpdate:
			bits[0]++
			bits[1], bits[2], bits[3] = 0, 0, 0
			break
		case randFloat < secondVersionUpdate:
			bits[1]++
			bits[2], bits[3] = 0, 0
			break
		case randFloat < thirdVersionUpdate:
			bits[2]++
			bits[3] = 0
			break
		default:
			bits[3]++
		}
		return result
	}
}

// updateVersionRangeFuncCreate 随机可更新版本 0 是最小可更新版本 1是最大可更新
func updateVersionRangeFunc() []string {
	minVersion := rand.Intn(len(updateVersionCodes) - 1)
	maxVersion := rand.Intn(len(updateVersionCodes)-minVersion) + minVersion
	return []string{updateVersionCodes[minVersion], updateVersionCodes[maxVersion]}
}

func whiteListNameSelectFuncCreate() func() int {
	fileCodes := make([]int, 16)
	for i := range fileCodes {
		fileCodes[i] = i
	}
	length := len(fileCodes)
	return func() int {
		randPos := rand.Intn(length)
		length--
		fileCodes[randPos], fileCodes[length] = fileCodes[length], fileCodes[randPos]
		return fileCodes[length]
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
