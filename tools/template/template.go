package template

import (
	"text/template"
)

var (
	Welcome ,_ = template.New("Welcome").Parse(`
=====================================================
欢迎进入在线游戏平台
版本v0.1.0

输入[list]查看当前游戏房间列表
输入[create 游戏类型]创建游戏房间 (例如: "create landlord" 创建斗地主房间)
输入[join 房间序号]加入游戏房间 (例如: "join 112")
输入[exit]退出游戏平台
`)
	RoomList, _ = template.New("RoomList").Parse(`
=====================================================
序号  |   房间名  |   游戏类型  | 当前在线人数  
{{range $s := . }}
{{.Id}} | {{.Name}} | {{.TypeName}}  |  {{.OnlineCounts}}
{{end}}
=====================================================
`)
	RoomCreateSuccess,_ = template.New("RoomCreateSuccess").Parse(`
房价创建成功!序号:{{.Id}}
`)
)
