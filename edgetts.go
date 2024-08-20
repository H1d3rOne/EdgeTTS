package EdgeTTS

import (
	"flag"
	"fmt"
	"github.com/H1d3rOne/EdgeTTS/pkg"
)

// 新建一个TTS对象
func New() *pkg.TTS {
	return &pkg.TTS{}
}

// 设置要转换的文本
func SetText(t *pkg.TTS, text string) {
	t.Text = text
}

// 设置从文件中导入文本
func SetFile(t *pkg.TTS, file string) {
	t.File = file
}

// 设置声音
func SetVoice(t *pkg.TTS, voice string) {
	t.VoiceArg = voice
}

// 列出所有支持的声音
func ListVoices(t *pkg.TTS) {
	err := t.ListVoices()
	if err != nil {
		fmt.Println(err)
		return
	}
}

// 设置语速
func SetRate(t *pkg.TTS, rate string) {
	t.Rate = rate
}

// 设置音量
func SetVolume(t *pkg.TTS, volume string) {
	t.Volume = volume
}

// 设置音调
func SetPitch(t *pkg.TTS, pitch string) {
	t.Pitch = pitch
}

func SetWordsInCue(t *pkg.TTS, wordsInCue int) {
	t.WordsInCue = wordsInCue
}

// 设置音频输出文件
func SetWriteMedia(t *pkg.TTS, writeMedia string) {
	t.WriteMedia = writeMedia
}

// 设置字幕输出文件
func SetWriteSubtitles(t *pkg.TTS, writeSubtitles string) {
	t.WriteSubtitles = writeSubtitles
}

// 执行转换
func Run(t *pkg.TTS) {
	if t.Text == "" && t.File == "" {
		fmt.Println("请提供一个文本或文件进行转换")
		flag.Usage()
		return
	} else if t.Text == "" && t.File != "" {
		t.Text, _ = pkg.ReadFile(t.File)
	} else if t.Text != "" && t.File == "" {

	} else {
		fmt.Println("文本和文件只需输入一个就行")
		return
	}

	if t.VoiceArg == "" {
		t.VoiceArg = "zh-CN-YunjianNeural"
	}
	if t.Rate == "" {
		t.Rate = "+0%"
	}
	if t.Volume == "" {
		t.Volume = "+0%"
	}
	if t.Pitch == "" {
		t.Pitch = "+0Hz"
	}
	if t.WordsInCue == 0 {
		t.WordsInCue = 10
	}
	if t.WriteMedia == "" {
		t.WriteMedia = "output.mp3"
	}

	//pkg.DetectTerminal(t.WriteMedia)
	//验证参数
	err := t.Validator()
	if err != nil {
		fmt.Println(err)
		return
	}

	subs := pkg.NewSubtitle()
	if t.WriteSubtitles == "" {
		err := t.SaveAudio()
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		err := t.SaveAudioAndSubs(subs)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
