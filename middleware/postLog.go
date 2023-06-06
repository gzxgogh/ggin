package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gzxgogh/ggin/config"
	"github.com/gzxgogh/ggin/db"
	"github.com/gzxgogh/ggin/logs"
	"github.com/gzxgogh/ggin/utils"
	"io/ioutil"
	"strings"
	"time"
)

type PostLog struct {
	Time          string            `json:"time" bson:"time"`
	ResponseTime  string            `json:"responseTime" bson:"responseTime"`
	TTL           int               `json:"ttl" bson:"ttl"`
	AppName       string            `json:"appName" bson:"appName"`
	Method        string            `json:"method" bson:"method"`
	ContentType   string            `json:"contentType" bson:"contentType"`
	Uri           string            `json:"uri" bson:"uri"`
	ClientIP      string            `json:"clientIP" bson:"clientIP"`
	RequestHeader map[string]string `json:"requestHeader" bson:"requestHeader"`
	RequestParam  any               `json:"requestParam" bson:"requestParam"`
	RequestBody   any               `json:"requestBody" bson:"requestBody"`
	ResponseStr   string            `json:"responseStr" bson:"responseStr"`
	ResponseMap   any               `json:"responseMap" bson:"responseMap"`
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

var accessChannel = make(chan string, 100)

func RequestLogger() gin.HandlerFunc {
	go handleAccessChannel()

	return func(c *gin.Context) {

		startTime := time.Now()
		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter
		data, err := c.GetRawData()
		if err != nil {
			logs.Error("GetRawData error:", err.Error())
		}
		body := string(data)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data)) // 关键点

		// 处理请求
		c.Next()
		responseBody := bodyLogWriter.body.String()
		// 日志格式
		if strings.Contains(c.Request.RequestURI, "/docs") || c.Request.RequestURI == "/" {
			return
		}

		var result any
		if responseBody != "" && responseBody[0:1] == "{" {
			err := json.Unmarshal([]byte(responseBody), &result)
			if err != nil {
				result = map[string]any{"status": -1, "msg": "解析异常:" + err.Error()}
			}
		}

		// 结束时间
		endTime := time.Now()
		// 日志格式
		var params, reqBody any
		if strings.Contains(c.ContentType(), "application/json") && body != "" {
			utils.FromJSON(body, &reqBody)
		}
		params = utils.GinParamMap(c)
		postLog := new(PostLog)
		postLog.Time = startTime.Format("2006-01-02 15:04:05")
		postLog.Uri = c.Request.RequestURI
		postLog.Method = c.Request.Method
		postLog.AppName = config.Cfg.App.Name
		postLog.ContentType = c.ContentType()
		postLog.RequestHeader = utils.GinHeaders(c)
		ip := c.GetHeader("X-Forward-For")
		if ip == "" {
			ip = c.GetHeader("X-Real-IP")
			if ip == "" {
				ip = c.ClientIP()
			}
		}
		postLog.ClientIP = ip
		postLog.RequestParam = params
		postLog.RequestBody = reqBody
		postLog.ResponseTime = endTime.Format("2006-01-02 15:04:05")
		postLog.ResponseMap = result
		postLog.TTL = int(endTime.UnixNano()/1e6 - startTime.UnixNano()/1e6)

		//dao.Log.PostLog(*postLog)
		if config.Cfg.Logger.RequestTableName != "" {
			accessChannel <- utils.ToJSON(postLog)
		}
	}
}

func handleAccessChannel() {
	for accessLog := range accessChannel {
		var postLog PostLog
		json.Unmarshal([]byte(accessLog), &postLog)
		dbName := config.Cfg.Logger.RequestTableName
		if dbName == "" {
			logs.Error("未配置请求日志的MongoDB数据库表")
			continue
		}
		conn := db.DBObj.GetMongoConn()
		if conn == nil {
			logs.Error("获取MongoDB连接错误")
		}
		err := conn.C(config.Cfg.Logger.RequestTableName).Insert(&postLog)
		if err != nil {
			logs.Error("写入请求日志到MongoDB错误")
		}
	}
	return
}
