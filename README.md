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
tts := new  
client.Input("audio.mp3")  
client.Format("srt")  
client.Output("subtitle.srt")  
client.Run()  
```


