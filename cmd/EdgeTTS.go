package main

import (
	"flag"
	"fmt"
	"github.com/H1d3rOne/EdgeTTS/pkg"
)

var (
	text           = flag.String("t", "", "Text to be translated")
	file           = flag.String("f", "", "File to be translated")
	voice          = flag.String("v", "zh-CN-YunjianNeural", "Voice to be used,default voice:zh-CN-YunjianNeural")
	listVoices     = flag.Bool("list-voices", false, "List voices")
	rate           = flag.String("rate", "+0%", "Rate of the voice")
	volume         = flag.String("volume", "+0%", "Volume of the voice")
	pitch          = flag.String("pitch", "+0Hz", "Pitch of the voice")
	wordsInCue     = flag.Int("words-in-cue", 10, "Words in cue")
	writeMedia     = flag.String("write-media", "output.mp3", "Write media to file")
	writeSubtitles = flag.String("write-subtitles", "", "Write subtitles to file")
	proxy          = flag.String("proxy", "", "Proxy to be used")
)

//var C *config.Config

func main() {
	RunCmd()
}
func RunCmd() {

	flag.Parse()
	t := pkg.NewTTS()

	if *text == "" && *file == "" {
		if *listVoices {
			err := t.ListVoices()
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			fmt.Println("请提供一个文本或文件进行转换")
			flag.Usage()
			return
		}

	} else if *text == "" && *file != "" {
		if *file == "/dev/stdin" {
			t.Text, _ = pkg.ReadFromStdin()
		} else {
			t.File = *file
			t.Text, _ = pkg.ReadFile(t.File)
		}
	} else if *text != "" && *file == "" {
		t.Text = *text
	} else {
		fmt.Println("文本和文件只需输入一个就行")
		flag.Usage()
		return
	}

	if *listVoices {
		err := t.ListVoices()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	t.VoiceArg = *voice
	t.Rate = *rate
	t.Volume = *volume
	t.Pitch = *pitch
	t.WordsInCue = *wordsInCue
	t.WriteMedia = *writeMedia
	t.WriteSubtitles = *writeSubtitles
	RunTTS(t)

}

func RunTTS(t *pkg.TTS) {
	//验证参数
	err := t.Validator()
	if err != nil {
		fmt.Println(err)
		return
	}
	//检测终端
	pkg.DetectTerminal(t.WriteMedia)
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
