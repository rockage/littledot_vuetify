package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"mysql_con"
	"strconv"
	"time"
)

// sessions.go
var (
	cookieNameForSessionID = "mycookiesessionnameid"
	sess                   = sessions.New(sessions.Config{Cookie: cookieNameForSessionID})
)

func CheckLogin(ctx iris.Context) bool {
	session := sess.Start(ctx)
	auth, _ := session.GetBoolean("authenticated")
	if !auth {
		return false
	}
	return true
}
func secret(ctx iris.Context) {
	// 检查用户是否已通过身份验证
	session := sess.Start(ctx)
	auth, _ := session.GetBoolean("authenticated")

	if !auth {
		//ctx.StatusCode(iris.StatusForbidden)
		_, _ = ctx.WriteString("token check fail")
		return
	}

	_, _ = ctx.WriteString(session.GetString("username") + "|" + session.GetString("uid"))
}

// login
func login(ctx iris.Context) {
	session := sess.Start(ctx)
	// 在此处进行身份验证
	username := ctx.FormValue("username")
	password := ctx.FormValue("password")
	var rst []map[string]string
	rst, _ = mysql_con.Query("select uid from pre_members where username = '" + username + "' and `password` = '" + password + "'")
	if rst != nil {
		uid := rst[0]["uid"]
		session.Set("authenticated", true)
		session.Set("username", username)
		session.Set("uid", uid)
		unixTime := strconv.FormatInt(time.Now().Unix(), 10)
		mysql_con.Exec("update pre_members set lastvisited = " + unixTime + " where uid =" + rst[0]["uid"]) //写入最后登录时间
		_, _ = ctx.WriteString(uid)

	} else {
		_, _ = ctx.Text("not found")
	}
}

func userZhuce(ctx iris.Context) {
	session := sess.Start(ctx)
	username := ctx.FormValue("username")
	password := ctx.FormValue("password")
	unixTime := strconv.FormatInt(time.Now().Unix(), 10)

	sql := "INSERT INTO pre_members ( username, password, regdate, lastvisited) VALUES ('" + username + "','" + password + "','" + unixTime + "','" + unixTime + "')"
	uid := mysql_con.Exec(sql) // 新增用户
	session.Set("authenticated", true)
	session.Set("username", username)
	session.Set("uid", uid)
	_, _ = ctx.WriteString(uid)
}

func userExist(ctx iris.Context) {
	username := ctx.FormValue("username")
	var rst []map[string]string
	rst, _ = mysql_con.Query("select uid from pre_members where username = '" + username + "'")
	if rst != nil {
		_, _ = ctx.WriteString("sucess")
	} else {
		_, _ = ctx.Text("fail")
	}
}

func logout(ctx iris.Context) {
	session := sess.Start(ctx)
	session.Set("authenticated", false)
	_, _ = ctx.WriteString("退出登录。")
}
