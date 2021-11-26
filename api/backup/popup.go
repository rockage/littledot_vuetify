package main

import (
	"encoding/base64"
	"encoding/json"
	"github.com/kataras/iris/v12"
	"io/ioutil"
	"mysql_con"
	"os"
	"regexp"
	"strings"
)





// new password
func SetNewPasswd(ctx iris.Context) {
	uid := ctx.FormValue("uid")
	username := ctx.FormValue("username")
	oldpasswd := ctx.FormValue("oldpasswd")
	newpasswd := ctx.FormValue("newpasswd")
	var rst []map[string]string
	rst, _ = mysql_con.Query("select uid from pre_members where uid = " + uid + " and username = '" + username + "' and password = '" + oldpasswd + "'")
	if rst != nil {
		mysql_con.Exec("UPDATE pre_members set password = '" + newpasswd + "' where uid = " + uid)
	} else {
		ctx.Text("error")
	}
}

// getprofile
func GetProfile(ctx iris.Context) {
	uid := ctx.FormValue("uid")
	var rst []map[string]string
	rst, _ = mysql_con.Query("select email,gender,location,born,mobile,signature,avatar from pre_members where uid = " + uid)
	if rst != nil {
		b, err := json.Marshal(rst)
		if err == nil {
			_, _ = ctx.JSON(string(b))
		}
	} else {
		ctx.Text("error")
	}

}
//setProfile
func SetProfile(ctx iris.Context) {
	if CheckLogin(ctx) == false {
		ctx.Text("error")
		return
	}

	uid := ctx.FormValue("uid")
	password := ctx.FormValue("password")
	avatar := ctx.FormValue("avatar")
	gender := ctx.FormValue("gender")
	location := ctx.FormValue("location")
	born := ctx.FormValue("born")
	mobilePhone := ctx.FormValue("mobilePhone")
	signature := ctx.FormValue("signature")

	var fn string
	re, _ := regexp.Compile(`^(data:.+?;base64,)`)
	match := re.FindAllStringSubmatch(avatar, -1)

	if match == nil {
		fn = avatar //如果data不是base64，直接存储图片路径
	} else {
		avatar = re.ReplaceAllString(avatar, "")
		data, _ := base64.StdEncoding.DecodeString(avatar)
		fn = saveAvatar(data, uid)
	}
	var rst []map[string]string
	rst, _ = mysql_con.Query("select uid from pre_members where uid = " + uid + " and password = '" + password + "'")
	if rst != nil {
		mysql_con.Exec("UPDATE pre_members set avatar = '" + fn + "', " +
			"gender = " + gender + "," +
			"location = " + location + "," +
			"born = '" + born + "'," +
			"mobile = '" + mobilePhone + "'," +
			"signature = '" + signature + "' " +
			"where uid = " + uid)
		ctx.Text("个人资料修改成功。")
	} else {
		ctx.Text("error")
	}
}
func saveAvatar(data []byte, uid string) string {
	dir := "../front/static/avatar/"
	_, err := os.Stat(dir)
	if err != nil {
		err = os.Mkdir(dir, 0777)
	}
	fileName := dir + uid + ".jpg"
	err = ioutil.WriteFile(fileName, data, 0666)
	fileName = strings.Replace(fileName, "../front", "", -1) //去掉../front，否则前端无法读取
	return fileName
}

//getUserProfile
func GetUserProfile(ctx iris.Context) {
	if CheckLogin(ctx) == false {
		ctx.Text("error")
		return
	}
	uname := ctx.FormValue("uname")
	var rst []map[string]string
	rst, _ = mysql_con.Query("select uid,posts,threads,lastvisited,regdate,email,gender,location,born,mobile,signature,avatar from pre_members where username = '" + uname + "'")
	if rst != nil {
		b, err := json.Marshal(rst)
		if err == nil {
			_, _ = ctx.JSON(string(b))
		}
	} else {
		ctx.Text("error")
	}
}
