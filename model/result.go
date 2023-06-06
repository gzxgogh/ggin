package model

type Result struct {
	Status int         `json:"status" bson:"status"`
	Msg    string      `json:"msg" bson:"msg"`
	Data   interface{} `json:"data" bson:"data"`
	Page   *ResultPage `json:"page,omitempty" bson:"page,omitempty"`
}

type ResultPage struct {
	Count int `json:"count"` //总页数
	Index int `json:"index"` //页号
	Size  int `json:"size"`  //分页大小
	Total int `json:"total"` //总记录数
}

func Error(code int, msg string) Result {
	return Result{
		Status: code,
		Msg:    msg,
		Data:   nil,
	}
}

func Success(data interface{}) Result {
	return Result{
		Status: 200,
		Msg:    "success",
		Data:   data,
	}
}
