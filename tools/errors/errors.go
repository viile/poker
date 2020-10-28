package errors

import (
	"errors"
	"fmt"
)

type Err struct {
	e error
	p error
}

func NewErr(text string) Err {
	return Err{e: errors.New(text)}
}

func (e Err) Error() string {
	var str string
	if e.p != nil {
		str = fmt.Sprintf("parent error : %s ", e.p.Error())
	}
	if e.e != nil {
		str += e.e.Error()
	}
	return str
}

func (e Err) With(err error) Err {
	return Err{e: err, p: e}
}

func (e Err) WithMsg(msg string) Err {
	return Err{e: NewErr(msg), p: e}
}

var (
	ErrObjectNotFound = NewErr("数据不存在")
	ErrParser         = NewErr("命令解析失败")
	ErrUserExist      = NewErr("用户已存在")
	ErrUserNotExist   = NewErr("用户不存在")
	ErrPassword       = NewErr("密码错误")
	ErrType           = NewErr("数据类型不符合预期")
	ErrRoomNotExist   = NewErr("房间不存在")
)
