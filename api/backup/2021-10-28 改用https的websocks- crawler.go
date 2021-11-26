package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"mysql_con"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/gobwas/ws" //Websocks工具
	"github.com/gobwas/ws/wsutil"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

// ***** 全局变量、常量 ***** ：
const (
	period = 30   // 30分钟自动检测一次，一个tick = 1分钟
	low    = 0.93 // 鼠标滑行速度上下限
	high   = 1.01
	// myos   = "windows"
	myos    = "linux" // win or linux 总开关 ####
	version = "v1027"
)

var SysPath string
var Conn net.Conn
var Selenium T_Selenium
var SliderCount int // 滑块遭遇次数

var timer Timer
var time_tick int
var run_count int

func tick() { // 定时器时间到了执行此函数
	time_tick++
	if time_tick >= period {
		if Selenium.wd != nil {
			Selenium.Check()
			time_tick = 0
			run_count++
			sm(strconv.Itoa(period-time_tick) + "分钟后将自动执行一次抓单")
		} else {
			sm("没有检测到Selenium")
		}
	}
}

type Timer struct {
	Interval int
	Tick     func()
	Channel  chan bool
	Running  bool
}

func (t Timer) Enabled() {
	ticker := time.NewTicker(time.Duration(t.Interval) * time.Second)
	go func() {
		for {
			select {
			case <-t.Channel:
				t.Running = false
				ticker.Stop()
				return
			case <-ticker.C:
				t.Tick()
			}
		}
	}()
}

type T_Order struct {
	id           string
	tb_id        string
	date         string
	vendor_id    string
	vendor_class string
	total_price  string
	note         string
	red_note     string
	address      string
	state        string
	Sub_Orders   []Sub_Order
	flag         int
}

type Sub_Order struct {
	sub_id       string
	order_id     string
	p_id         string
	p_name       string
	p_voltage    string
	amount       string
	price        string
	shipped_date string
	express      string
	tracking     string
	state        string
}

type T_Selenium struct {
	wd        selenium.WebDriver
	service   *selenium.Service
	orders    *[]T_Order
	automatic bool
}

func (s T_Selenium) SwitchPage(w int) bool {
	windows, _ := s.wd.WindowHandles()
	var handle string = ""
	for k, v := range windows {
		if k == w {
			handle = v
		}
	}
	if handle != "" {
		s.wd.SwitchWindow(windows[w])
		return true
	} else {
		return false
	}
}

func (s T_Selenium) SlowKeys(str string, element selenium.WebElement) {
	c := []byte(str)
	i := 0
	for { // 模拟慢速输入字符
		element.SendKeys(string(c[i]))
		i++
		if i >= len(c) {
			break
		}
		Sleep(100, 200) // 随机暂停一下
	}
}

func (s T_Selenium) CheckElement(strCond string, timeOut int, milliSecond ...bool) selenium.WebElement {
	// 根据条件判断元素是否存在，单位：秒
	var element selenium.WebElement
	var cond selenium.Condition // selenium的复杂条件模块：等待，直到条件成立或超时
	cond = func(wd selenium.WebDriver) (bool, error) {
		wd = s.wd
		element, _ = s.wd.FindElement(selenium.ByXPATH, strCond)
		if element != nil {
			return true, nil
		} else {
			return false, nil
		}
	}
	var err error
	if len(milliSecond) == 0 {
		err = s.wd.WaitWithTimeout(cond, time.Duration(timeOut)*time.Second) // 限定时间内找不到指定元素，err != nil (如不指定milliSecond，默认秒为单位)
	} else {
		err = s.wd.WaitWithTimeout(cond, time.Duration(timeOut)*time.Millisecond) // 如指定milliSecond，则以毫秒为单位)
	}
	if err != nil {
		return nil
	} else {
		return element //如果找到元素，返回这个元素
	}

}

func (s T_Selenium) CheckLogin() bool {
	// 登录状态检查
	windows, _ := s.wd.WindowHandles()
	if windows[0] != "" {
		s.wd.SwitchWindow(windows[0]) // 0: 卖家中心
	} else { // 同时销毁另外两个窗口
		if windows[1] != "" {
			s.wd.CloseWindow(windows[1]) // 1: 物流
		}
		if windows[2] != "" {
			s.wd.CloseWindow(windows[2]) // 2: 已卖出的宝贝
		}
	}
	sm("正在检查登录特征... ")
	element := s.CheckElement(`//*[@id="module-open-aside"]/div/div/div/div/ul/li[1]/div[1]/span[2]/div/div/div`, 5)
	if element == nil {
		sm("目前浏览器不在登录状态")
		return false
	} else {
		sm("目前浏览器处于登录状态")
		return true
	}
}

func (s T_Selenium) ClickSms(smscode string) {
	// 输入短信验证码并点击 （由外层WSS触发）
	Sleep(1000, 1100)
	bitmap := robotgo.OpenBitmap(SysPath + "login_sms_logo.png")
	defer robotgo.FreeBitmap(bitmap)
	basex, basey := robotgo.FindBitmap(bitmap)
	if basex != -1 {
		rx := basex - 720
		ry := basey - 115
		robotgo.MoveMouseSmooth(rx+Rnd(1, 5), ry+Rnd(1, 5), low, high)
		Sleep(500, 600)
		robotgo.MouseClick("left", true)
		Sleep(500, 600)
		robotgo.TypeStr(smscode)
		Sleep(1000, 1500)
		robotgo.KeyTap("enter")
		sm("输入完毕，正在验证短信验证码...")
		Sleep(1000, 2000)
		if s.CheckLogin() {
			sm("短信验证成功，登录正常")
		} else {
			sm("登录失败，请登录VPS查看状态")
		}
	} else {
		sm("短信验证码模块发生意外")
	}
}

func (s T_Selenium) Login() {
	// 登录
	SlowInput := func(str string) {
		c := []byte(str)
		i := 0
		for {
			robotgo.TypeStr(string(c[i]))
			i++
			if i >= len(c) {
				break
			}
			Sleep(100, 200) // 随机暂停一下
		}
	}
	// ***** Login 函数主体 *****：
	Sleep(500, 1000)
	sm("开始访问 https://login.taobao.com")
	s.wd.Get("https://login.taobao.com")
	sm("尝试密码登录, 等待页面稳定")
	var basex, basey, rx, ry int

	login_wait_count := 1
	for {
		Sleep(1000, 1100)
		bitmap := robotgo.OpenBitmap(SysPath + "login_base.png")
		defer robotgo.FreeBitmap(bitmap)
		basex, basey = robotgo.FindBitmap(bitmap)
		if basex != -1 {
			sm("开始输入用户名")
			rx = basex + 810 // 淘宝的宝字，最后一点为base (x+810 , y-270) = 输入框
			ry = basey + 270
			Sleep(500, 600)
			robotgo.MoveMouseSmooth(rx+Rnd(1, 5), ry+Rnd(1, 5), low, high)
			Sleep(500, 600)
			robotgo.MouseClick("left", true)
			Sleep(500, 600)
			SlowInput("sword_yang:carol")
			sm("开始输入密码")
			rx = basex + 810
			ry = basey + 330 // base > x+813, y-330 = 密码框
			Sleep(500, 600)
			robotgo.MoveMouseSmooth(rx+Rnd(1, 5), ry+Rnd(1, 5), low, high)
			Sleep(500, 600)
			robotgo.MouseClick("left", true)
			Sleep(500, 600)
			SlowInput("kkndcc110")
			Sleep(500, 600)
			robotgo.KeyTap("enter")
			break
		} else {
			login_wait_count++
			if login_wait_count >= 10 { // 10秒仍然无法进入登录界面
				sm("无法找到用户名输入框，登录失败！")
				return
			}
		}
	}

	unlock_count := 1      // 解锁次数
	unlock_wait_count := 1 // 解锁等待次数
	for {
		Sleep(1000, 1100)
		bitmap := robotgo.OpenBitmap(SysPath + "login_slider.png")
		fx, fy := robotgo.FindBitmap(bitmap)
		if fx != -1 { // 找到了滑块
			sm("遭遇登录滑块，正在尝试解锁")
			// 尝试解锁
			Sleep(500, 600)
			robotgo.MoveMouseSmooth(fx+Rnd(5, 10), fy+Rnd(5, 10), low, high)
			Sleep(50, 150)
			robotgo.DragSmooth(fx+Rnd(270, 285), fy)
			Sleep(2000, 3000) // 等待解锁完毕
			robotgo.KeyTap("enter")
			unlock_count++
			if unlock_count >= 3 { // 3次解锁失败
				sm("登陆时连续解锁失败！")
				return
			}
		} else {
			unlock_wait_count++
			if unlock_wait_count >= 3 { // 等待超过3秒，页面已稳定，没有滑块
				sm("没有滑块")
				break
			}
		}
	}

	sms_wait_count := 1 // 验证码界面等待时间
	for {
		Sleep(1000, 1100)
		bitmap := robotgo.OpenBitmap(SysPath + "login_sms_logo.png")
		fx, _ := robotgo.FindBitmap(bitmap)
		if fx != -1 { // 进入了短信验证码页面
			sm("遭遇短信验证码模块，请注意查收短信")
			rx = basex + 465
			ry = basey + 215
			Sleep(500, 600)
			robotgo.MoveMouseSmooth(rx+Rnd(5, 10), ry+Rnd(5, 10), low, high)
			Sleep(2000, 3000)
			robotgo.MouseClick("left", true) // 点击获取短信
			Sleep(2500, 3500)
			robotgo.MouseClick("left", true) // 点击获取短信
			Sleep(500, 600)
			return
		} else {
			sms_wait_count++
			if sms_wait_count >= 3 { // 等待超过3秒，页面已稳定，没有短信验证码
				sm("没有短信验证码")
				break
			}
		}
	}
	if s.CheckLogin() {
		sm("登录成功!")
	} else {
		sm("登录失败，需要登录VPS查看状态")
	}
	return
}

func (s T_Selenium) CreateWuliuPage() bool {
	var timeout bool

	tick := func() {
		sm("物流页加载严重超时，准备重启")
		timeout = true
	}

	timer := new(Timer)
	timer.Interval = 30 // 30秒仍未打开页面则加载失败
	timer.Tick = tick
	timer.Channel = make(chan bool)

	exist := func() bool {
		element, _ := s.wd.FindElement(selenium.ByXPATH, `//*[@id="J_FilterBox"]/div[1]/ul/li[1]/span/input`)
		Sleep(500, 501)
		if element != nil {
			return true
		} else {
			return false
		}
	}

	Sleep(500, 600)
	s.SwitchPage(0)
	sm("正在创建物流页")
	Sleep(1000, 2000)
	bitmap := robotgo.OpenBitmap(SysPath + "items_wuliu.png")
	defer robotgo.FreeBitmap(bitmap)
	fx, fy := robotgo.FindBitmap(bitmap)
	if fx != -1 {
		rx := fx + 5
		ry := fy + 35
		Sleep(500, 600)
		robotgo.MoveMouseSmooth(rx+Rnd(1, 5), ry+Rnd(1, 5), low, high)
		Sleep(500, 600)
		robotgo.MouseClick("left", true)
		Sleep(500, 600)
		s.SwitchPage(1)
	}
	sm("正在检查物流页面特征 ...(30S超时)")
	R := false
	timer.Enabled()
	for {
		R = exist() // 反复查询物流页特征，直至出现或超时
		switch R {
		case true:
			timer.Channel <- false // 关闭定时器
			sm("物流页面状态正常")
			s.wd.ResizeWindow("", 1380, 1024)
			return true
		case false:
			if timeout {
				timer.Channel <- false // 关闭定时器
				sm("物流页面状态错误")
				return false
			}
		}
	}
}

func (s T_Selenium) CreateSoldPage() bool {
	var timeout bool
	tick := func() {
		sm("已卖出的宝贝页加载严重超时，准备重启")
		timeout = true
	}
	timer := new(Timer)
	timer.Interval = 30 // 30秒仍未打开页面则加载失败
	timer.Tick = tick
	timer.Channel = make(chan bool)

	exist := func() bool {
		element, _ := s.wd.FindElement(selenium.ByXPATH, `//*[@id="auctionId"]`)
		if element != nil {
			return true
		} else {
			return false
		}
	}

	Sleep(500, 600)
	s.SwitchPage(0)
	sm("正在创建已卖出的宝贝页")
	Sleep(1000, 2000)
	bitmap := robotgo.OpenBitmap(SysPath + "items_sold.png")
	defer robotgo.FreeBitmap(bitmap)
	fx, fy := robotgo.FindBitmap(bitmap)
	if fx != -1 {
		rx := fx + 10
		ry := fy + 35
		Sleep(500, 600)
		robotgo.MoveMouseSmooth(rx+Rnd(5, 10), ry+Rnd(5, 10), low, high)
		Sleep(500, 600)
		robotgo.MouseClick("left", true)
		Sleep(500, 600)
		s.SwitchPage(2)
	}
	sm("正在检查已卖出的宝贝页面特征 ...(30S超时)")
	R := false
	timer.Enabled()
	for {
		R = exist()
		switch R {
		case true:
			timer.Channel <- false // 关闭定时器
			sm("已卖出的宝贝页面状态正常")
			s.wd.ResizeWindow("", 1380, 1024)
			return true
		case false:
			if timeout {
				timer.Channel <- false // 关闭定时器
				sm("已卖出的宝贝页面状态错误")
				return false
			}
		}
	}
}

func (s T_Selenium) CreateWuliuText() string { // 爬取物流页文字 : 本页反爬比较宽松，未出现过滑块
	var wuliu_text string
	var fx int
	KeyPress := func(key string, times int) {
		for i := 0; i < times; i++ {
			robotgo.KeyTap(key)
			Sleep(150, 300)
		}
	}
	// **** 函数主体 ****
	Sleep(1000, 2000)
	s.SwitchPage(1)
	s.wd.ResizeWindow("", 1380, 1024)
	Sleep(500, 600)
	sm("小眼睛扫描开始...")
	Sleep(1000, 2000)

	scan_count := 0

	for { // -> 外层翻页循环
		scan_count++
		Sleep(1000, 2000)
		sm("正在扫描已物流第" + strconv.Itoa(scan_count) + "页")
		loop_quit := false
		for { // -> 内层循环 ： 找到的条目
			bitmap := robotgo.OpenBitmap(SysPath + "smalleye.png")
			defer robotgo.FreeBitmap(bitmap)
			points := robotgo.FindEveryBitmap(bitmap)
			if points != nil {
				for _, v := range points {
					robotgo.MoveMouseSmooth(v.X+Rnd(1, 2), v.Y+Rnd(1, 5), low, high)
					robotgo.MouseClick("left", true)
					Sleep(300, 500)
				}
			}
			kd := Rnd(3, 5)
			KeyPress("down", kd) // 滚动频率，即键盘down键
			if loop_quit == true {
				break // 退出内层循环
			}

			bitmap = robotgo.OpenBitmap(SysPath + "wuliu_end_point.png")
			fx, _ = robotgo.FindBitmap(bitmap)
			if fx != -1 {
				loop_quit = true
			}
		} // -> 内层for ： 条目循环
		KeyPress("home", 1) // 回到最上
		// 页面文字:
		element_text, err := s.wd.FindElement(selenium.ByXPATH, `//*[@id='J_Express']`)
		Sleep(1000, 2000)
		if err == nil {
			str, _ := element_text.Text()
			wuliu_text += str
		} else {
			fmt.Println("error message:", err)
			fmt.Println("element_text:", element_text)
		}
		element_nextpage := s.CheckElement(`//*[@id="J_Express"]/tfoot/tr/td/div[2]/div/a[2]/span`, 5) //下一页按钮
		if element_nextpage == nil {
			sm("已无法找到[下一页]按钮，物流文字创建结束")
			break
		} else {
			// wuliu base point 到 下一页按钮 距离：40, 410
			bitmap := robotgo.OpenBitmap(SysPath + "wuliu_base.png")
			defer robotgo.FreeBitmap(bitmap)
			fx, fy := robotgo.FindBitmap(bitmap)
			if fx != -1 {
				rx := fx + 40
				ry := fy + 410
				Sleep(500, 600)
				robotgo.MoveMouseSmooth(rx+Rnd(1, 5), ry+Rnd(1, 5), low, high)
				Sleep(100, 110)
				robotgo.MouseClick("left", true)
			}
		}
	} // 外层for : 翻页循环
	ioutil.WriteFile("wuliu_text.txt", []byte(wuliu_text), 0644) // 存盘 for test
	return wuliu_text
}

// 物流页与已卖出页采用两种不同的页面布局及由此派生的两种反爬手段：
// 物流页的DOM元素没有固定名称，需要点击小眼睛才能将地址显示（与一个唯一ID绑定），这导致所有DOM元素都是由不固定的ID构成的
// 已卖出页DOM元素是固定的，但是鼠标悬浮在小红旗上才会实时生成一个新的DOM结点
// 以上特征决定了爬取方式的不同：
// 1. 在物流页因为没有固定DOM编号，只能整页读取后再正则（负责：买家地址、买家留言），处理函数返回值是页面的裸体文字
// 2. 在已卖出页由于存在实时生成的内容，只能一边生成一边读取DOM元素（负责：订单总价、产品列表、卖家留言），函数返回的是已经处理好的JSON字符串

func (s T_Selenium) CreateSoldText() string { // 读取订单的“小红旗”留言（即：卖家留言）

	KeyPress := func(key string, times int) {
		for i := 0; i < times; i++ {
			robotgo.KeyTap(key)
			Sleep(150, 300)
		}
	}
	var fx int

	type Sold_Sub_Orders struct {
		P_name string
		Price  string
		Amount string
	}

	type Sold_Order struct {
		Tb_id       string
		Total_price string
		Red_Note    string
		Sub_Orders  []Sold_Sub_Orders
	}

	var sold_orders []Sold_Order

	UnlockRedNoteSlider := func() bool {
		unlock_wait_count := 0
		for {
			Sleep(1000, 1100)
			bitmap_unlock_wait := robotgo.OpenBitmap(SysPath + "punish_wait.png")
			defer robotgo.FreeBitmap(bitmap_unlock_wait)
			fx, _ := robotgo.FindBitmap(bitmap_unlock_wait)

			if fx != -1 { // 找到了圆圈
				unlock_wait_count++
			} else {
				unlock_wait_count--
			}
			if unlock_wait_count >= 10 { // 等待圆圈出现超过10秒
				return false
			}
			if unlock_wait_count <= -3 { // 找不到等待圆圈3秒(页面已稳定)
				break
			}
		}

		bitmap_unlock := robotgo.OpenBitmap(SysPath + "punish.png")
		defer robotgo.FreeBitmap(bitmap_unlock)
		fx, fy := robotgo.FindBitmap(bitmap_unlock)
		if fx != -1 {
			rx := fx - 105 // rx, ry : 按钮坐标
			ry := fy + 125
			Sleep(500, 600)
			robotgo.MoveMouseSmooth(rx+Rnd(5, 10), ry+Rnd(5, 10), low, high)
			Sleep(50, 150)
			robotgo.DragSmooth(rx+Rnd(270, 285), ry)
			Sleep(5000, 6000) // 等待解锁完毕
			return false
		} else {
			return true
		}
	}

	check_help_window := func() {
		bitmapcs := robotgo.OpenBitmap(SysPath + "help.png")
		defer robotgo.FreeBitmap(bitmapcs)
		rx, ry := robotgo.FindBitmap(bitmapcs)
		if rx != -1 {
			sm("遭遇客服帮助弹窗，尝试关闭")
			robotgo.MoveMouseSmooth(rx, ry, low, high)
			Sleep(100, 200)
			robotgo.MouseClick("left", true)
		}
	}

	// **** 函数主体 ****
	Sleep(1000, 2000)
	s.SwitchPage(2)
	s.wd.ResizeWindow("", 1380, 1024)
	Sleep(500, 600)
	sm("点击等待发货按钮")

	KeyPress("pageup", 5)
	Sleep(3000, 5000)

	bitmap := robotgo.OpenBitmap(SysPath + "sold_base.png")
	defer robotgo.FreeBitmap(bitmap)
	basex, basey := robotgo.FindBitmap(bitmap)
	if fx != -1 {
		rx := basex - 750 // base - (75,180) = [等待发货]
		ry := basey - 170
		Sleep(1500, 1600)
		robotgo.MoveMouseSmooth(rx+Rnd(1, 5), ry+Rnd(1, 5), low, high)
		Sleep(500, 600)
		robotgo.MouseClick("left", true)
		if !UnlockRedNoteSlider() {
			SliderCount++ // 全局变量
			sm("Click: [等待发货]遭遇滑块，模块重启")
			return "unlock error"
		}
	} else {
		sm("无法点击等待发货按钮")
		return "error"
	}
	Sleep(3000, 5000)
	sm("小红旗留言扫描开始...")
	// 确定div编号
	// 解析：DIV层："sold_container"（即包含了订单文字的那个DIV层）其编号是浮动的，会根据淘宝网的公告数量而经常变化
	// 因公告的文字量不可能有订单文字那么多，因此此处枚举所有编号从1-20的div层，找出文字数量最大那个即为订单文字
	type Length struct {
		index  int
		length int
	}
	var lengths []Length

	t := 1
	for {
		e, err := s.wd.FindElement(selenium.ByXPATH, `//*[@id="sold_container"]/div/div[`+strconv.Itoa(t)+`]`)
		if err == nil {
			text, _ := e.Text()
			var one Length
			one.index = t
			one.length = len(text)
			lengths = append(lengths, one)
		}
		t++
		if t >= 20 {
			break
		}
	}
	// 求最大值：
	max := 0
	var di int // div index
	for _, v := range lengths {
		if v.length > max {
			max = v.length
			di = v.index
		}
	}
	sm("订单div层：" + strconv.Itoa(di) + ",开始创建shadow node")

	Sleep(3000, 5000)
	check := s.CheckElement(`//*[@id="sold_container"]/div/div[`+strconv.Itoa(di-2)+`]/div[1]/div[3]`, 5) // 确定是否真的在等待发货页面
	if check != nil {
		str, _ := check.GetAttribute("outerHTML")
		r := regexp.MustCompile(`selected`) // 检测是否存在selected字样
		m := r.FindStringSubmatch(str)
		if len(m) > 0 {
			sm("已确定页面位置：等待发货")
			if !UnlockRedNoteSlider() { // 有可能产生滑块的一步
				SliderCount++ // 全局变量
				sm("小红旗页遭遇滑块，模块重启")
				return "unlock error"
			}
		} else {
			sm("【严重】 不是等待发货页面")
			return "error"
		}
	}

	// SCRIPT注入：
	script := `	if (!document.getElementById("my_shadow")) {
						s=document.createElement("shadow");
						s.setAttribute("id","my_shadow");
						s.setAttribute("value","");
						document.body.append(s);	
				}; `
	s.wd.ExecuteScript(script, nil) // 创建一个DOM结点，可以利用这个结点将返回值return回来

	scan_count := 0
	for { // -> 外层翻页循环
		scan_count++
		Sleep(3000, 5000)
		sm("正在扫描已卖出的宝贝第" + strconv.Itoa(scan_count) + "页")
		// 首先截取页面中非实时生成部分：
		i := 2
		for { // 大循环：div[x] x=订单行 （从2开始，满页是16，共15行）

			var one Sold_Order
			var j int = 1
			order_e := s.CheckElement(`//*[@id="sold_container"]/div/div[`+strconv.Itoa(di)+`]/div[`+strconv.Itoa(i)+`]/table[1]/tbody/tr/td[1]/label/span[3]`, 1)
			if order_e != nil {
				ret, _ := order_e.Text()
				one.Tb_id = ret // 订单ID
				price_e := s.CheckElement(`//*[@id="sold_container"]/div/div[`+strconv.Itoa(di)+`]/div[`+strconv.Itoa(i)+`]/table[2]/tbody/tr/td[7]/div/div[1]/p/strong/span[2]`, 50, true)
				if price_e != nil {
					ret, _ := price_e.Text()
					one.Total_price = ret // 订单总价
				}
			} else {
				break // 退出大循环
			}

			for { // 小循环： tr[x] x=产品列表行 （从1开始，行数由购买的产品种类的数量决定）
				var one_sub Sold_Sub_Orders
				product_e := s.CheckElement(`//*[@id="sold_container"]/div/div[`+strconv.Itoa(di)+`]/`+
					`div[`+strconv.Itoa(i)+`]/table[2]/tbody/`+
					`tr[`+strconv.Itoa(j)+`]/td[1]/div/div[2]/p/a/span[2]`, 50, true) // 50ms短超时，因为DOM已经渲染完毕，无需等待
				if product_e != nil {
					ret, _ := product_e.Text()
					one_sub.P_name = ret
					e := s.CheckElement(`//*[@id="sold_container"]/div/div[`+strconv.Itoa(di)+`]/div[`+strconv.Itoa(i)+`]/table[2]/tbody/tr[`+strconv.Itoa(j)+`]/td[1]/div/div[2]/p[2]/span/span[3]`, 50, true)
					if e != nil {
						ret, _ = e.Text()
						one_sub.P_name += ret // 产品名拼合 = 产品名+颜色分类
					}
					e = s.CheckElement(`//*[@id="sold_container"]/div/div[`+strconv.Itoa(di)+`]/div[`+strconv.Itoa(i)+`]/table[2]/tbody/tr[`+strconv.Itoa(j)+`]/td[2]/div/p/span[2]`, 50, true)
					if e != nil {
						ret, _ = e.Text()
						one_sub.Price = ret // 产品单价
					}
					e = s.CheckElement(`//*[@id="sold_container"]/div/div[`+strconv.Itoa(di)+`]/div[`+strconv.Itoa(i)+`]/table[2]/tbody/tr[`+strconv.Itoa(j)+`]/td[3]/div/p`, 50, true)
					if e != nil {
						ret, _ = e.Text()
						one_sub.Amount = ret // 产品数量
					}
					one.Sub_Orders = append(one.Sub_Orders, one_sub)
				} else {
					break // 退出小循环
				}
				j++
			}
			sold_orders = append(sold_orders, one)
			i++
		}

		// 开始处理需要实时生成的数据（指鼠标必须移动到小红旗上即时生成的动态DOM元素）
		index_red := 5      // 小红旗元素编号从5开始递增
		break_loop := false // 退出扫描标志
		// 前置工作完毕，开始扫描：
		rnd_x := 0
		rnd_y := 0
		for { // -> 内层循环 ： 找到的条目
			bitmap := robotgo.OpenBitmap(SysPath + "redflag.png")
			points := robotgo.FindEveryBitmap(bitmap)
			check_help_window() // 是否存在客户帮助弹窗？

			if points != nil {
				for _, v := range points {
					if v.Y >= 640 && v.Y <= 740 { // 从Y轴 640-740这一段是捕获窗口，只有进入这个窗口的元素才会被捕获
						robotgo.MoveMouseSmooth(v.X+Rnd(8, 11), v.Y+Rnd(10, 13), low, high)
						rnd_x = Rnd(v.X-50, v.X+50)
						rnd_y = Rnd(v.Y-50, v.Y+50)
						Sleep(100, 200)
						// 当index_red增加了之后，即使重复扫描到了上一个红旗，但真实页面并没有生成新的动态div层,
						// 因此不会触发 index_red++ ，也不需要担心唯一性问题， index_red 从编号5开始，总是随着新的div层产生而增加
						e_red := s.CheckElement(`//*[@id="list-sold-items"]/div[`+strconv.Itoa(index_red)+`]`, 500, true) // 即时生成的小红旗div层
						if e_red != nil {
							index_red++
							px := strconv.Itoa(v.X)
							py := strconv.Itoa(v.Y)
							script := `	var oPoint = document.elementFromPoint(` + px + `,` + py + `);
										var oID = oPoint.getAttribute("data-reactid");
										document.getElementById("my_shadow").setAttribute("value", oID);`
							// my_shadow 是初始化时创建的影子node
							s.wd.ExecuteScript(script, nil)
							shadow_result, err := s.wd.FindElement(selenium.ByID, "my_shadow") // 将影子node带回主程序

							var order_id string // 从页面中返回相对应的order_id（通过JS的elementFromPoint函数）
							if err == nil {
								order_id, _ = shadow_result.GetAttribute("value")
								r := regexp.MustCompile(`\$(\d{1,})\.`)
								m := r.FindStringSubmatch(order_id) // 正则截取到订单id
								if len(m) > 0 {
									order_id = m[1]
								} else {
									order_id = "!error!"
								}
							}
							red_note, _ := e_red.Text()
							for k, v := range sold_orders {
								if order_id == v.Tb_id {
									sold_orders[k].Red_Note = red_note
								}
							} // -> 将搜索到的red_note 赋值到sold_orders序列里
						}
						Sleep(200, 300)
						robotgo.MoveMouseSmooth(rnd_x, rnd_y, low, high) // 随机落点
						Sleep(200, 300)

					} // -> if 进入窗口
				} // -> 屏幕上所有找到的小红旗标志
			} // -> 内层循环
			if break_loop {
				break
			}
			kd := Rnd(2, 3)
			KeyPress("down", kd) // 滚动频率，即键盘down键
			bitmap = robotgo.OpenBitmap(SysPath + "sold_end_point.png")
			fx, _ = robotgo.FindBitmap(bitmap)
			if fx != -1 {
				break_loop = true
			}
		} // -> 内层for ： 条目循环
		KeyPress("home", 1)
		Sleep(500, 600)
		check := s.CheckElement(`//*[@id="sold_container"]/div/div[`+strconv.Itoa(di)+`]/div[1]/div[3]/div/button[2]`, 5) // 检测下一页按钮
		if check != nil {
			str, _ := check.GetAttribute("outerHTML")
			r := regexp.MustCompile(`disabled`) // 检测是否存在disabled字样
			m := r.FindStringSubmatch(str)
			if len(m) > 0 {
				sm("下一页按钮已不能点击，小红旗扫描结束。")
				break
			} else {
				// 点击下一页按钮:
				rx := basex - 35 // base - (35,145) = [下一页]
				ry := basey - 130
				fmt.Println(" rx , ry : ", rx, " ", ry)
				Sleep(500, 600)
				robotgo.MoveMouseSmooth(rx+Rnd(1, 5), ry+Rnd(1, 5), low, high)
				Sleep(100, 110)
				robotgo.MouseClick("left", true)
				if !UnlockRedNoteSlider() { // 有可能产生滑块的一步
					SliderCount++ // 全局变量
					sm("小红旗页遭遇滑块，模块重启")
					return "unlock error"
				}
			} // -> if disable
		} // -> if button exist

	} // 外层for : 翻页循环

	b, _ := json.Marshal(sold_orders) // 将结果保存为json存盘 (for test)
	str := string(b)
	ioutil.WriteFile("sold_text.txt", []byte(str), 0644)
	Sleep(1500, 1600)
	return "success"
}

func (s T_Selenium) MakeOrders(wuliu_text string, sold_text string) []T_Order {
	// 形参说明： wuliu_text = 裸文字 ，sold_text = json字符串
	var SQL string

	type Sold_Sub_Orders struct {
		P_name string
		Price  string
		Amount string
	}

	type Sold_Order struct {
		Tb_id       string
		Total_price string
		Red_Note    string
		Sub_Orders  []Sold_Sub_Orders
	}

	var sold_orders []Sold_Order

	type Product_List struct { // 预读产品库
		p_id      string
		p_name    string
		tb_name   string
		isvoltage string
	}
	SQL = "select id, name, tb_name, isvoltage from ld_products"
	rst, _ := mysql_con.Query(SQL)
	var ProductList = make([]Product_List, len(rst))
	for k, v := range rst {
		ProductList[k].p_id = v["id"]
		ProductList[k].p_name = v["name"]
		ProductList[k].tb_name = v["tb_name"]
		ProductList[k].isvoltage = v["isvoltage"]
	}

	type Vendor_List struct { // 预读经销商库
		vendor_id      string
		vendor_tb_name string
	}

	SQL = "select id, taobao from ld_vendor where class = 2"
	rst, _ = mysql_con.Query(SQL)
	var VendorList = make([]Vendor_List, len(rst))
	for k, v := range rst {
		VendorList[k].vendor_id = v["id"]
		VendorList[k].vendor_tb_name = v["taobao"]
	}

	var orders []T_Order

	reg := regexp.MustCompile(`(订单编号：[\s\S]*?)###333###666###999###`) // 切割物流字符串（大段）
	match := reg.FindAllString(wuliu_text, -1)                        // 订单与订单之间，返回组

	orders = make([]T_Order, len(match))
	json.Unmarshal([]byte(sold_text), &sold_orders) // 将sold_text恢复为变量

	sm("开始正则物流页文字")
	for k, v := range match { // 最外层for -> 将物流页切割成很多个主订单

		r := regexp.MustCompile(`订单编号：(\d{1,})[\r\n]创建时间：(.+?)[\r\n]`)
		m := r.FindStringSubmatch(v)

		if len(m) > 0 {
			orders[k].tb_id = m[1]
			orders[k].date = m[2]
			orders[k].vendor_id = "84" // 默认经销商：自营
			orders[k].vendor_class = "1"
			orders[k].state = "2" //待发货
		}

		r = regexp.MustCompile(`收货信息：[\r\n](.+?)[\r\n]`)
		m = r.FindStringSubmatch(v)
		// 合成淘宝风格收货地址：
		r_add := regexp.MustCompile(`^(.+?)\s{3}(.+?),.+?,\s{1}(.+?),\s{1}(.+?)$`)
		m_add := r_add.FindStringSubmatch(m[1])
		if len(m_add) > 0 {
			r = regexp.MustCompile(` `)
			trim_space := r.ReplaceAllString(m_add[2], "") //去掉原文里的空格
			address := "{" + m_add[1] + "}" + " " + m_add[3] + "，" + m_add[4] + "，" + trim_space
			orders[k].address = address
			// 处理经销商
			tb_name := m_add[1]
			for _, v := range VendorList {
				if tb_name == v.vendor_tb_name {
					orders[k].vendor_id = v.vendor_id
					orders[k].vendor_class = "2" // 搜出来的默认就是经销商
					break
				}
			}
		} else {
			orders[k].address = "地址抓取错误，订单编号：" + orders[k].tb_id
		}

		r = regexp.MustCompile(`买家留言：[\r\n](.+?)[\r\n]`)
		m = r.FindStringSubmatch(v)
		if len(m) > 0 {
			orders[k].note = "[" + m[1] + "]"
		}

		// 统一 wuliu_text 和 sold_text ：
		for _, v_sold := range sold_orders { // 第二层for : 引入小红旗页，比对小红旗页主订单和物流页主订单
			if orders[k].tb_id == v_sold.Tb_id {
				if v_sold.Red_Note != "" {
					orders[k].note += "[" + v_sold.Red_Note + "]" // 整合买家留言 和 卖家留言
				}
				orders[k].total_price = v_sold.Total_price // 实际价格
				set_price := 0                             // 预设价格
				var suborders []Sub_Order
				suborders = make([]Sub_Order, len(v_sold.Sub_Orders)) // 在产品列表已知的情况下，不需要用append

				for j, v_sold_sub := range v_sold.Sub_Orders { // 第三层for : 处理子订单
					suborders[j].p_name = v_sold_sub.P_name
					suborders[j].p_id = "38"                // 默认ID: 自定义
					suborders[j].p_voltage = "1"            // 默认电压：国规
					suborders[j].amount = v_sold_sub.Amount // 数量

					// 预设价格处理：
					float_num, _ := strconv.ParseFloat(v_sold_sub.Price, 16) // 原字符串是带小数点的浮点数，因此要用ParseFloat
					P := int(float_num)
					A, _ := strconv.Atoi(v_sold_sub.Amount)
					T := P * A
					suborders[j].price = strconv.Itoa(T) // 子订单价格 = 单价*数量
					set_price += T                       // 子订单价格总计

					suborders[j].shipped_date = ""
					suborders[j].express = ""
					suborders[j].tracking = ""
					suborders[j].state = "2"

					//此处处理p_id 和 p_name 的对应关系
					for _, p := range ProductList {

						if v_sold_sub.P_name == p.tb_name { // 如产品能在数据库中找到，用数据库的
							suborders[j].p_id = p.p_id
							suborders[j].p_name = p.p_name
							if p.isvoltage == "1" {
								// PS: 1国 2美 3英 4欧 5未知 6空 处理经销商电压，从留言中大致分析，以防止经销商拍下后不标电压的情况
								if orders[k].vendor_class == "2" {
									suborders[j].p_voltage = "5" // 如果接下来的正则全部不通过 或 没有留言，则默认为"未知"以提醒发货的时候注意
									if orders[k].note != "" {
										r := regexp.MustCompile(`(美|110|US|Us|us)`)
										m := r.FindStringSubmatch(orders[k].note) //返回：值
										if len(m) > 0 {
											suborders[j].p_voltage = "2"
										}
										r = regexp.MustCompile(`(英|UK|Uk|uk|港|HK|Hk|hk)`)
										m = r.FindStringSubmatch(orders[k].note) //返回：值
										if len(m) > 0 {
											suborders[j].p_voltage = "3"
										}
										r = regexp.MustCompile(`(欧|EU|Eu|eu|歐|法|德)`)
										m = r.FindStringSubmatch(orders[k].note) //返回：值
										if len(m) > 0 {
											suborders[j].p_voltage = "4"
										}
									}
								}
							} else {
								suborders[j].p_voltage = "6" // 如产品本身不插电，则电压选项为空
							}
							break
						} // -> 如果在产品库中找到对应的tb_name
					} // -> 结束p_id 与 p_name 对应关系的处理

				} // -> for 子订单
				orders[k].Sub_Orders = suborders // 子订单附加到主订单
				// 订单价格处理： 从物流页抓取的价格是商品的预设价格，往往和最终成交价格不一致
				// 公式： 差值 = 实际总价（小红旗价格） - 预设价格（物流页价格）
				// 将差值附加到 suborders[] 的第一个产品类目上
				float_real_price, _ := strconv.ParseFloat(v_sold.Total_price, 16)
				real_price := int(float_real_price)
				diff := real_price - set_price
				float_P, _ := strconv.ParseFloat(orders[k].Sub_Orders[0].price, 16)
				P := int(float_P)
				result := P + diff
				orders[k].Sub_Orders[0].price = strconv.Itoa(result)
				//fmt.Println("orders[k].Sub_Orders[0].price:  ", orders[k].Sub_Orders[0].price)

				//suborders_total_price // 从物流页抓取来的产品售价
				// 小红旗过来的总价
			}

		}
	} // -> for 主订单

	return orders
}

func (s T_Selenium) Compare(pageOrders []T_Order) []T_Order {
	var SQL string
	var compared_orders []T_Order
	// 生产订单列表（mysql端）
	SQL = "SELECT id, tb_id, note FROM ld_order where (state_id = 2 or state_id = 4) and tb_id <> ''"
	rst, _ := mysql_con.Query(SQL)
	mysql_orders := make([]T_Order, len(rst)) // 预读mysql中 state_id = 2、4的所有订单
	for k, v := range rst {
		mysql_orders[k].id = v["id"]
		mysql_orders[k].tb_id = v["tb_id"]
		mysql_orders[k].note = v["note"]
	}
	// 比对1：页面 -> mySQL (页面有而数据库没有，则新增这条订单到数据库，如留言发生变化，则更新在数据库中更新这条记录)
	for k, page := range pageOrders {
		add_new := 1 // 默认每条page.note都需要更新
		for _, mysql := range mysql_orders {
			if page.tb_id == mysql.tb_id { // 如mysql中存在相同的记录：
				if len(page.note) > len(mysql.note) { // 比较page侧和mysql侧留言长度，哪边的留言更长，就以那边为准，如以mysql为准则无需动作
					add_new = 2
					var one T_Order
					one = pageOrders[k]
					one.note = page.note
					one.flag = 2
					compared_orders = append(compared_orders, one) // 将需要做改变的记录续接在compared_orders上
				} else {
					add_new = 0
				}
				break // 退出内层循环
			}
		}

		var one T_Order
		if add_new == 1 {
			one = pageOrders[k]
			one.flag = 1
			compared_orders = append(compared_orders, one)
		}
	}
	// 比对2：mySQL -> 页面 （数据库有而页面没有，则在数据库中屏蔽这条订单）
	for _, mysql := range mysql_orders {
		disable := true
		for _, page := range pageOrders {
			if mysql.tb_id == page.tb_id {
				disable = false
				break
			}
		}
		if disable { // 需要屏蔽
			var one T_Order
			one.tb_id = mysql.tb_id
			one.flag = 3
			compared_orders = append(compared_orders, one)
		}
	}
	sm("订单分析结束")
	return compared_orders
}

func (s T_Selenium) Analyse() []T_Order {
	wuliu_text := ReadAndReg("wuliu_text.txt", `订单编号：\d{1,}`) // 读text 并做一些正则预处理
	b, _ := ioutil.ReadFile("sold_text.txt")                  // 读Json
	sold_text := string(b)
	var orders []T_Order
	orders = s.MakeOrders(wuliu_text, sold_text) // 将裸文字转换为页面端订单列表
	orders = s.Compare(orders)                   // 比较页面订单与数据库订单

	if len(orders) != 0 {
		state_str := map[int]string{1: "新增", 2: "修改留言", 3: "屏蔽"}
		voltage_str := map[int]string{1: "国规220V", 2: "美规110V", 3: "英规220V", 4: "欧规220V", 5: "未知000V", 6: ""}

		for k, v := range orders {
			sm("#ORDERS#" + "【No.（" + strconv.Itoa(k) + ")   " + state_str[v.flag] + "】")
			sm("#ORDERS#" + "下单时间:" + v.date + " 订单编号:" + v.tb_id + " 经销商ID:" + v.vendor_id + " 订单总价：" + v.total_price)
			sm("#ORDERS#" + "地址:" + v.address)
			sm("#ORDERS#" + "留言:" + v.note)
			for k2, v2 := range v.Sub_Orders {
				vid, _ := strconv.Atoi(v2.p_voltage)
				sm("#ORDERS#" + "产品" + strconv.Itoa(k2) + ":" + v2.p_name + " | " + voltage_str[vid] + "x" + v2.amount + " ¥" + v2.price)
			}
			sm("#ORDERS#" + "- - - - - - - - - - - - - - - - - - - -")
		}
	} else {
		sm("#ORDERS#" + "没有需要更新的订单。")
	}
	return orders
}

//m+i-
func (s T_Selenium) Save(finalOrders []T_Order) {
	if finalOrders != nil {
		var SQL string
		for _, v := range finalOrders {
			switch v.flag {
			case 0:
				sm("订单内容不变：" + v.tb_id)
				break
			case 1:
				// 假设在mySQL端刚刚点完发货按钮后(state=3)，还未来得及到淘宝页面点击发货，而
				// 此时恰好爬虫更新，就会错误地新增一条记录，此处仍需再做一次唯一性检测。（见图）
				SQL = `SELECT state_id from ld_order WHERE tb_id = '` + v.tb_id + `'`
				rst, _ := mysql_con.Query(SQL)
				if len(rst) > 0 { // mySQL中存在相同的tb_id，无论状态如何，都不添加
					sm("订单：" + v.tb_id + "(state:" + rst[0]["state_id"] + ")已存在!")
				} else {
					sm("正在新增订单：" + v.tb_id)
					SQL = "INSERT INTO ld_order " +
						"(date, vendor_id, tb_id, price, state_id, note, address)" +
						" VALUES " +
						"('" + v.date + "'," + v.vendor_id + ", " + "'" + v.tb_id + "', " + v.total_price + ", " +
						"2, '" + v.note + "', '" + v.address + "')"
					newOrderID := mysql_con.Exec(SQL) // 返回一个新的OrderID
					for _, v2 := range v.Sub_Orders {
						SQL = "INSERT INTO ld_order_suborder " +
							"(order_id, product_id, vendor_id, item_describe, voltage, amount, price, shiped_date, state)" +
							" VALUES " +
							"(" + newOrderID + ", " + v2.p_id + ", " + v.vendor_id + ", '" + v2.p_name + "', " +
							v2.p_voltage + ", " + v2.amount + ", " + v2.price + ", '1900-01-01 00:00:00', " + v2.state + ")"
						mysql_con.Exec(SQL)
					}

				}
				break
			case 2:
				sm("订单留言修改：" + v.tb_id)
				SQL = "UPDATE ld_order SET note = '" + v.note + "' WHERE tb_id = '" + v.tb_id + "'" // 主订单修改
				mysql_con.Exec(SQL)
				break
			case 3:
				sm("正在屏蔽订单：" + v.tb_id)
				SQL = "UPDATE ld_order SET state_id = 12 WHERE tb_id = '" + v.tb_id + "'" // 主订单修改
				mysql_con.Exec(SQL)
				SQL = "UPDATE ld_order_suborder SET state = 12 WHERE order_id = '" + v.id + "'" // 子订单修改, v.id = 主订单id (order_id)
				mysql_con.Exec(SQL)
				break
			}
		}
		sm("数据库写入结束，本次更新完毕!")
	} else {
		fmt.Println("finalOrders:", finalOrders)
	}
}

func (s T_Selenium) Check() bool {

	var try_count int
	var loop_quit bool

	windows, _ := s.wd.WindowHandles()
	for k, v := range windows {
		if k != 0 {
			s.wd.SwitchWindow(v)
			Sleep(500, 600)
			s.wd.Close()
		}
	}
	// 准备工作：
	sm("cls") // 客户端清屏
	// 1: 检测登录
	if s.CheckLogin() != true {
		sm("ERR: 登录验证失败。")
		return false
	}

	// 2 : 创建物流页面
	try_count = 0
	for {
		if s.CreateWuliuPage() != true {
			windows, _ := s.wd.WindowHandles() // 关闭窗口重启
			s.wd.CloseWindow(windows[1])
			try_count++
			if try_count > 2 {
				sm("ERR: 物流页异常，程序退出！")
				return false
			}
		} else {
			break
		}
	}

	// 3 : 创建已卖出的宝贝页面 + 抓取小红旗留言（返回JSON）
	try_count = 0
	loop_quit = false
	for {

		try_create_sold_count := 0
		for {
			if s.CreateSoldPage() != true {
				windows, _ := s.wd.WindowHandles() // 关闭窗口重启
				s.wd.CloseWindow(windows[2])
				try_create_sold_count++
				if try_create_sold_count > 2 {
					sm("ERR: 已卖出的宝贝页面异常，程序退出！")
					return false
				}
			} else {
				break
			}
		}

		red_note := s.CreateSoldText()
		switch red_note {
		case "error":
			sm("已卖出宝贝异常退出！")
			return false
		case "unlock error":
			try_count++
			sm("遭遇滑块，模块第 " + strconv.Itoa(try_count) + "次重启！")
			if try_count > 2 {
				sm("多次解锁失败，已卖出宝贝异常退出！")
				return false
			} else {
				windows, _ := s.wd.WindowHandles()
				s.wd.CloseWindow(windows[2])
			}
		case "success":
			sm("已卖出宝贝正常抓取完毕！")
			loop_quit = true
		}
		if loop_quit == true { // 引用loop_quit这个变量的原因：在switch case里break仅退出switch，而不是for
			break
		}
	}

	// 4: 抓取物流页（裸文本）
	PageWuliuText := s.CreateWuliuText()
	if PageWuliuText == "error" {
		sm("截取物流页文字异常退出！")
		return false
	}

	// 5: 如处于自动状态，则附加执行Analyse和Save
	if s.automatic == true {
		var temp_orders []T_Order
		temp_orders = s.Analyse()
		s.Save(temp_orders)
	}
	sm("恭喜！订单抓取全流程结束。")
	return true
}

// ####
func WebSocket() {
	// 初始化 Websocks、Timer

	timer.Interval = 60 // 60秒为1个tick
	timer.Tick = tick
	timer.Channel = make(chan bool)
	timer.Running = false

	var httpSrv *http.Server
	httpSrv = &http.Server{Addr: ":2052"}
	http.HandleFunc("/", WSS)
	httpSrv.ListenAndServe() //代码置于此行之后将不被执行，因为协程已经被Listen起来了
}

func startWsServer(w http.ResponseWriter, r *http.Request) net.Conn {
	Conn, _, _, _ = ws.UpgradeHTTP(r, w)
	return Conn
}

func WSS(w http.ResponseWriter, r *http.Request) {
	var temp_orders []T_Order

	// Websocks程序主体处理
	if Conn != nil {
		Conn.Close()
	}

	Conn := startWsServer(w, r)
	var msg_receive []byte
	str_conn := fmt.Sprintf("%d", Conn)
	sm("New Conn Created:" + str_conn + version)

	for {

		msg_receive, _, _ = wsutil.ReadClientData(Conn)
		msg_str := string(msg_receive)

		if msg_str == "" {
			continue
		}

		reg := regexp.MustCompile(`#rockage_cmds#(.+?)#rockage_datas#(.+)`)
		match := reg.FindStringSubmatch(string(msg_receive))

		if match == nil {
			continue
		}

		cmds := match[1]
		datas := match[2]
		fmt.Println(cmds)

		switch cmds {
		case "status":
			sm("- - - - -")
			sm("合计抓单次数：" + strconv.Itoa(run_count))
			sm("遭遇滑块次数：" + strconv.Itoa(SliderCount))
			sm("自动抓单阈值：" + strconv.Itoa(period) + "分钟/次")
			sm("当前Tick值：" + strconv.Itoa(time_tick))
			s := map[bool]string{true: "运行中", false: "已关闭"}
			sm("定时器自动状态：" + s[timer.Running])
			sm("- - - - -")

		case "restart":
			if Selenium.service != nil {
				Selenium.service.Stop() //停止之前的selenium实例
			}
			startChrome() // 初始化Selenium

		case "check":
			if Selenium.wd != nil {
				if Selenium.CheckLogin() == true {
				} else {
					go Selenium.Login()
				}
			} else {
				sm("没有检测到Selenium")
			}
		case "start":
			if Selenium.wd != nil {
				Selenium.Check()
			} else {
				sm("没有检测到Selenium")
			}

		case "automatic": // 自动运行Check
			switch datas {
			case "open":
				if timer.Running == false {
					time_tick = 0
					timer.Running = true // 只允许运行一个定时器
					Selenium.automatic = true
					timer.Enabled() // 定时器开始
				}
				sm(strconv.Itoa(period-time_tick) + "分钟后将自动执行一次抓单")
			case "close":
				if timer.Running == true {
					timer.Channel <- true // 关闭定时器
					timer.Running = false
					Selenium.automatic = false
					sm("定时器已关闭")
				} else {
					sm("定时器并未开启")
				}
			}

		case "analyse":
			temp_orders = Selenium.Analyse()

		case "save":
			Selenium.Save(temp_orders) // 上一步：Analyse() ，已将全局变量Selenium.orders生成

		case "sms":
			if Selenium.wd != nil {
				sm("收到验证码：" + datas)
				Selenium.ClickSms(datas)
			} else {
				sm("没有检测到Selenium")
			}

		default:

		}

	}
}

func startChrome() {
	opts := []selenium.ServiceOption{}
	caps := selenium.Capabilities{
		"browserName":      "chrome",
		"pageLoadStrategy": "eager",
	}
	imagCaps := map[string]interface{}{
		//		"profile.managed_default_content_settings.images": 2,
	}

	chromeCaps := chrome.Capabilities{
		Prefs:           imagCaps,
		Path:            "",
		ExcludeSwitches: []string{"enable-automation"}, // 重点
		Args: []string{
			//"--headless",
			`--window-position=-7,0`,
			`--window-size=1380,1024`,
			"--no-sandbox",
			"--disable-blink-features=AutomationControlled", // 重点
			"--user-agent=Mozilla/ 5.0(Windows NT 10.0; WOW64) Chrome/55.0.2883.87 Safari/537.36",
		},
	}
	caps.AddChrome(chromeCaps)

	if myos == "windows" {
		Selenium.service, _ = selenium.NewChromeDriverService("./windows/chromedriver.exe", 9515, opts...)
		SysPath = "./windows/"
	} else {
		Selenium.service, _ = selenium.NewChromeDriverService("./linux/chromedriver", 9515, opts...)
		SysPath = "./linux/"
	}

	Selenium.wd, _ = selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", 9515))
	sm("正在打开淘宝网首页")
	Selenium.wd.Get("https://taobao.com")
	Selenium.wd.ResizeWindow("", 1380, 1024)

	var check_str string
	if myos == "linux" {
		check_str = `//*[@id="root"]/div[1]/div[1]/div[2]/div[1]/input` // 淘宝国际站
	} else {
		check_str = `//*[@id="q"]` // 淘宝大陆站
	}
	ok := Selenium.CheckElement(check_str, 10) // 进入首页标志, 国际站和国内站有所区别
	if ok != nil {
		sm("正确访问淘宝网首页，selenium初始化成功")
	} else {
		sm("无法访问淘宝网首页，selenium初始化失败！")
	}

}

func sm(msg string) {
	var msg_send []byte = []byte(msg)
	wsutil.WriteServerMessage(Conn, ws.OpText, msg_send)
}

func Rnd(min int, max int) int { // 随机数产生器
	rand.Seed(time.Now().Unix())
	return min + rand.Intn(max-min)
}

func Sleep(min int, max int) { // 随机延时产生器
	rand.Seed(time.Now().Unix())
	time.Sleep(time.Duration(min+rand.Intn(max-min)) * time.Millisecond) // 毫秒
}

func ReadAndReg(fileName string, args ...string) string {
	var regMark string = ""
	if len(args) > 0 {
		regMark = args[0]
	}

	file, _ := os.Open(fileName)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	pre_line := ""
	all := ""

	for scanner.Scan() {
		cur_line := scanner.Text()

		if regMark != "" {
			rr := regexp.MustCompile(`\r`) // 去掉原文里的\r ， 否则无法正则
			mr := rr.ReplaceAllString(cur_line, "")
			cur_line = mr
			r := regexp.MustCompile(regMark)
			m := r.FindStringSubmatch(cur_line) //返回：值
			if len(m) > 0 {
				pre_line += "\n###333###666###999###\n" // 自定义切割标记
			}
		}
		all += pre_line + "\n"
		pre_line = cur_line

	}
	if regMark != "" {
		all += "\n###333###666###999###\n" // 文件结尾再加一行标记
	}
	return all
}
