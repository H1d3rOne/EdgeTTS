## 1、命令行运行 
```
go build EdgeTTS.go
```
## 2、导入包  
```
import "github.com/H1d3rOne/EdgeTTS" 
```
## 3、使用方法  
``` 
tts := EdgeTTS.New()
EdgeTTS.SetText(tts, "自古多情空余恨，此恨绵绵无绝期。")
EdgeTTS.SetWriteSubtitles(tts, "output.srt")
EdgeTTS.Run(tts)
EdgeTTS.ListVoices(tts)
```


