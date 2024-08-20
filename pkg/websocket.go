package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
)

type Result struct {
	Data  map[string]interface{}
	Error error
}

// 发送websocket请求获取音频流
func GetStream(t *TTS) <-chan Result {
	result := make(chan Result)
	var err error
	texts, err := splitTextByByteLength(
		escape(RemoveIncompatibleCharacters(t.Text)),
		calcMaxMsgSize(t.VoiceArg, t.Rate, t.Volume, t.Pitch),
	)
	finalUtterance := make(map[int]int)
	prevIdx := -1
	shiftTime := -1

	dialer := websocket.DefaultDialer
	dialer.Proxy = http.ProxyFromEnvironment

	//u, err := url.Parse(config.C.GetString("WSS_URL") + "&ConnectionId=" + ConnectID())
	u, err := url.Parse("wss://speech.platform.bing.com/consumer/speech/synthesize/readaloud/edge/v1?TrustedClientToken=6A5AA1D4EAFF4E9FB37E23D68491D6F4&ConnectionId=" + ConnectID())
	if err != nil {
		panic(err)
	}

	var downloadAudio bool
	var audioWasReceived bool

	//out := make(chan map[string]interface{})

	for idx, text := range texts {
		go func() {
			defer close(result)
			conn, _, err := dialer.Dial(u.String(), nil)
			if err != nil {
				websocketErr := &WebSocketError{
					Msg: "WebSocket连接错误",
					Err: err,
				}
				result <- Result{Error: websocketErr}
				return
			}
			defer conn.Close()

			date := dateToString()
			requestStr := fmt.Sprintf(
				"X-Timestamp:%s\r\nContent-Type:application/json; charset=utf-8\r\nPath:speech.config\r\n\r\n"+
					`{"context":{"synthesis":{"audio":{"metadataoptions":{"sentenceBoundaryEnabled":false,"wordBoundaryEnabled":true},"outputFormat":"audio-24khz-48kbitrate-mono-mp3"}}}}`,
				date,
			)
			if err := conn.WriteMessage(websocket.TextMessage, []byte(requestStr)); err != nil {
				websocketErr := &WebSocketError{
					Msg: "WebSocket发送消息错误",
					Err: err,
				}
				result <- Result{Error: websocketErr}
				return
			}

			ssml := ssmlHeadersPlusData(ConnectID(), date, Mkssml(t.VoiceArg, t.Rate, t.Volume, t.Pitch, []byte(text)))
			if err := conn.WriteMessage(websocket.TextMessage, []byte(ssml)); err != nil {
				websocketErr := &WebSocketError{
					Msg: "WebSocket发送消息错误",
					Err: err,
				}
				result <- Result{Error: websocketErr}
				return
			}
			//循环获取返回数据
		loop:
			for {
				messageType, msg, err := conn.ReadMessage()
				if err != nil {
					//out <- map[string]interface{}{"error": err.Error()}
					//return
					websocketErr := &WebSocketError{
						Msg: "WebSocket读取消息错误",
						Err: err,
					}
					result <- Result{Error: websocketErr}
					return
				}

				switch messageType {
				case websocket.TextMessage:
					parameters, data := getHeadersAndData(msg)
					path := parameters["Path"]
					switch string(path) {
					case "turn.start":
						// Download audio indicator is handled here
						downloadAudio = true
					case "turn.end":
						// End of audio data
						downloadAudio = false
						break loop
					case "audio.metadata":
						var metaData struct {
							Metadata []struct {
								Type string `json:"Type"`
								Data struct {
									Offset   int `json:"Offset"`
									Duration int `json:"Duration"`
									Text     struct {
										Text string `json:"Text"`
									} `json:"text"`
								} `json:"Data"`
							} `json:"Metadata"`
						}
						if err := json.Unmarshal(data, &metaData); err != nil {
							jsonErr := &JsonError{Msg: "JSON解析错误", Err: err}
							result <- Result{Error: jsonErr}
							return
						}
						for _, metaObj := range metaData.Metadata {
							metaType := metaObj.Type
							if idx != prevIdx {
								shiftTime = 0
								for i := 0; i < idx; i++ {
									shiftTime += finalUtterance[i]
								}
								prevIdx = idx
							}
							if metaType == "WordBoundary" {
								finalUtterance[idx] = metaObj.Data.Offset + metaObj.Data.Duration + 8750000
								//保存文本数据
								textData := map[string]interface{}{
									"type":     metaType,
									"offset":   metaObj.Data.Offset + shiftTime,
									"duration": metaObj.Data.Duration,
									"text":     metaObj.Data.Text.Text,
								}
								result <- Result{Data: textData}
							} else if metaType != "SessionEnd" {
								unkonwnResponse := &UnknownResponse{Msg: "未知的元数据类 " + metaType, Err: err}
								result <- Result{Error: unkonwnResponse}
								return
							}
						}
					case "response":
						continue
					default:
						unkonwnResponse := &UnknownResponse{Msg: "无法识别的响应", Err: err}
						result <- Result{Error: unkonwnResponse}
						return
					}
				case websocket.BinaryMessage:
					if !downloadAudio {
						unexpectedResponse := &UnexpectedResponse{Msg: "收到的不是期待的二进制消息", Err: err}
						result <- Result{Error: unexpectedResponse}
						return
					}
					if len(msg) < 2 {
						unexpectedResponse := &UnexpectedResponse{Msg: "收到的二进制消息不是请求头数据", Err: err}
						result <- Result{Error: unexpectedResponse}
						return
					}
					headerLength := int(msg[0])<<8 + int(msg[1])
					if len(msg) < headerLength+2 {
						unexpectedResponse := &UnexpectedResponse{Msg: "收到的二进制消息不是音频数据", Err: err}
						result <- Result{Error: unexpectedResponse}
						return
					}
					// 保存音频数据
					audioData := map[string]interface{}{
						"type": "audio",
						"data": msg[2+headerLength:],
					}
					result <- Result{Data: audioData}
					audioWasReceived = true
				case websocket.CloseMessage:
					unexpectedResponse := &UnexpectedResponse{Msg: "未知错误", Err: err}
					result <- Result{Error: unexpectedResponse}
					return
				default:
					unexpectedResponse := &UnexpectedResponse{Msg: "不可期待的WebSocket消息类型", Err: err}
					result <- Result{Error: unexpectedResponse}
					return
				}
			}

			if !audioWasReceived {
				otherErr := &OtherError{Msg: "没有收到音频数据，请验证你的参数是否正确", Err: err}
				result <- Result{Error: otherErr}
				return
			}
		}()

	}
	return result
}
