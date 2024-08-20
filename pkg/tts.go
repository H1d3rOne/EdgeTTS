package pkg

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

func NewTTS() *TTS {
	return &TTS{}
}

func (t *TTS) Validator() error {

	// Validate and normalize voice
	voiceMatch, err := regexp.MatchString(`^([a-z]{2,})-([A-Z]{2,})-(.+Neural)$`, t.VoiceArg)
	if err != nil {
		return errors.New("正则匹配出错")
	}

	if !voiceMatch {
		return &InvalidInputError{Message: "无效的声音格式，格式如：zh-CN-XiaoxiaoNeural"}
	}

	rateMatch, err := regexp.MatchString(`^[+-]\d+%$`, t.Rate)
	if err != nil {
		return errors.New("正则匹配出错")
	}
	if !rateMatch {
		return &InvalidInputError{Message: "无效的速率格式，格式如：+10%"}
	}

	volumeMatch, err := regexp.MatchString(`^[+-]\d+%$`, t.Volume)
	if err != nil {
		return errors.New("正则匹配出错")
	}
	if !volumeMatch {
		return &InvalidInputError{Message: "无效的音量格式，格式如：+10%"}
	}

	pitchMatch, err := regexp.MatchString(`^[+-]\d+Hz$`, t.Pitch)
	if err != nil {
		return errors.New("正则匹配出错")
	}
	if !pitchMatch {
		return &InvalidInputError{Message: "无效的声音pitch格式，格式如：+10Hz"}
	}
	return nil

}

// 列出所用的声音
func (t *TTS) ListVoices() error {
	//url := config.C.GetString("VOICE_LIST")
	url := "https://speech.platform.bing.com/consumer/speech/synthesize/readaloud/voices/list?trustedclienttoken=6A5AA1D4EAFF4E9FB37E23D68491D6F4"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errors.New("创建http请求失败")
	}
	req.Header.Set("Authority", "speech.platform.bing.com")
	req.Header.Set("Sec-CH-UA", `" Not;A Brand";v="99", "Microsoft Edge";v="91", "Chromium";v="91"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36 Edg/91.0.864.41")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{},
		Proxy:           http.ProxyFromEnvironment, // or http.ProxyURL(http.URL{Scheme: "http", Host: proxy})
	}
	client := &http.Client{Transport: transport}
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("请求列出所有声音失败")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("读取响应体失败")
	}

	var voices []Voice
	err = json.Unmarshal(body, &voices)
	if err != nil {
		return &JsonError{Msg: "JSON解析错误", Err: err}
	}
	for _, voice := range voices {
		fmt.Printf("Name：%s\n", voice.ShortName)
		fmt.Printf("Gender：%s\n\n", voice.Gender)
	}
	return nil
}

// 保存音频文件
func (t *TTS) SaveAudio() error {
	if t.WriteMedia != "" {
		audioFile, err := os.Create(t.WriteMedia)
		if err != nil {
			return err
		}

		for result := range GetStream(t) {
			err := result.Error
			if err != nil {
				return err
			}
			chunk := result.Data
			if chunk["type"] == "audio" {
				_, err := audioFile.Write(chunk["data"].([]byte))
				if err != nil {
					return &AudioWriteError{Msg: "音频写入失败", Err: err}
				}
			}
		}
		println("音频成功生成")
	}
	return nil
}

// 保存音频文件和字幕文件
func (t *TTS) SaveAudioAndSubs(subs *Subtitle) error {
	if t.WriteMedia != "" {
		audioFile, err := os.Create(t.WriteMedia)
		if err != nil {
			return err
		}

		for result := range GetStream(t) {
			err := result.Error
			if err != nil {
				return err
			}
			chunk := result.Data
			if chunk["type"] == "audio" {
				_, err := audioFile.Write(chunk["data"].([]byte))
				if err != nil {
					return &AudioWriteError{Msg: "音频写入失败", Err: err}
				}
			} else if chunk["type"] == "WordBoundary" {

				subs.CreateSub([2]float64{float64(chunk["offset"].(int)), float64(chunk["duration"].(int))}, chunk["text"].(string))
			}
		}
		println("音频成功生成")
	}
	if t.WriteSubtitles != "" {
		//subsFile, err := os.Create(t.WriteSubtitles)
		subsFile, err := os.OpenFile(t.WriteSubtitles, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			panic(err)
		}
		defer subsFile.Close()
		subString, err := subs.GenerateSubs(t.WordsInCue)
		if err != nil {
			return err
		}
		// 设置文件编码格式，不然写入文件为ANSI编码而出现中文乱码
		utf8bom := []byte{0xEF, 0xBB, 0xBF}
		_, err = subsFile.Write(utf8bom)
		if err != nil {
			return err
		}
		_, err = subsFile.WriteString(subString)
		if err != nil {
			return err
		}
		err = subsFile.Close()
		if err != nil {
			return err
		}
		//println(subString)
		println("字幕成功生成")
	}
	return nil
}
