package pkg

import "fmt"

// 处理输入参数的错误
type InvalidInputError struct {
	Message string
}

func (e *InvalidInputError) Error() string {
	return e.Message
}

// 处理websocket的错误
type WebSocketError struct {
	Msg string
	Err error
}

func (e *WebSocketError) Error() string {
	//return fmt.Sprintf("WebSocket请求错误: %v", e.Err)
	return fmt.Sprintf("%s: %v", e.Msg, e.Err)
}

type AudioWriteError struct {
	Msg string
	Err error
}

func (e *AudioWriteError) Error() string {
	return fmt.Sprintf("%s: %v", e.Msg, e.Err)
}

type SubtitleError struct {
	Msg string
	Err error
}

func (e *SubtitleError) Error() string {
	return fmt.Sprintf("%s: %v", e.Msg, e.Err)
}

type JsonError struct {
	Msg string
	Err error
}

func (e *JsonError) Error() string {
	return fmt.Sprintf("s: %v", e.Msg, e.Err)
}

type UnknownResponse struct {
	Msg string
	Err error
}

func (e *UnknownResponse) Error() string {
	return fmt.Sprintf("%s: %v", e.Msg, e.Err)
}

type UnexpectedResponse struct {
	Msg string
	Err error
}

func (e *UnexpectedResponse) Error() string {
	return fmt.Sprintf("%s: %v", e.Msg, e.Err)
}

type OtherError struct {
	Msg string
	Err error
}

func (e *OtherError) Error() string {
	return fmt.Sprintf("%s: %v", e.Msg, e.Err)
}
