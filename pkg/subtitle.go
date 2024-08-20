package pkg

import (
	"fmt"
	"html"
	"math"
	"strings"
)

// Formatter formats the subtitle entry.
func formatter(startTime, endTime float64, subdata string) string {
	return fmt.Sprintf("%s --> %s\n%s\n\n", mktimestamp(startTime), mktimestamp(endTime), html.EscapeString(subdata))
}

// Mktimestamp converts a time unit to the subtitle timecode format.
func mktimestamp(timeUnit float64) string {
	hour := int(math.Floor(timeUnit / 1e7 / 3600))
	minute := int(math.Floor((timeUnit / 1e7 / 60) - float64(hour*60)))
	seconds := timeUnit / 1e7
	return fmt.Sprintf("%02d:%02d:%06.3f", hour, minute, seconds)
}

// SubMaker holds the subtitle data.
type Subtitle struct {
	offset []struct{ start, end float64 }
	subs   []string
}

// NewSubMaker creates a new SubMaker instance.
func NewSubtitle() *Subtitle {
	return &Subtitle{}
}

// 创建一个字幕结构体
func (sb *Subtitle) CreateSub(timestamp [2]float64, text string) {
	//fmt.Println(timestamp)
	//fmt.Println(text)
	sb.offset = append(sb.offset, struct{ start, end float64 }{timestamp[0], timestamp[0] + timestamp[1]})
	sb.subs = append(sb.subs, text)
}

// 产生完整的字幕格式文件
func (sb *Subtitle) GenerateSubs(wordsInCue int) (string, error) {
	if len(sb.subs) != len(sb.offset) {
		//fmt.Println("subs and offset are not of the same length:%v，%v", len(sb.subs), len(sb.offset))
		return "", &SubtitleError{Msg: "时间戳长度与字幕内容长度不一致"}
	}

	if wordsInCue <= 0 {
		return "", &SubtitleError{Msg: "wordsInCue参数必须大于0"}
	}

	var data strings.Builder
	//data.WriteString("WEBVTT\n\n")

	subStateCount := 0
	subStateStart := -1.0
	subStateSubs := ""

	for idx, offset := range sb.offset {
		startTime, endTime := offset.start, offset.end
		//subs := html.UnescapeString(sb.subs[idx])
		subs := sb.subs[idx]

		if len(subStateSubs) > 0 {
			subStateSubs += " "
		}
		subStateSubs += subs

		if subStateStart == -1.0 {
			subStateStart = float64(startTime)
		}
		subStateCount++

		if subStateCount == wordsInCue || idx == len(sb.offset)-1 {
			subsText := subStateSubs
			var splitSubs []string
			//for i := 0; i < len(subsText); i += 79 {
			//	end := i + 79
			//	if end > len(subsText) {
			//		end = len(subsText)
			//	}
			//
			//	fmt.Println(subsText[i:end])
			//	fmt.Println("---")
			//	splitSubs = append(splitSubs, subsText[i:end])
			//}
			splitSubs = append(splitSubs, subsText)

			for i := 0; i < len(splitSubs)-1; i++ {
				sub := splitSubs[i]
				splitAtWord := true
				if strings.HasSuffix(sub, " ") {
					splitSubs[i] = sub[:len(sub)-1]
					splitAtWord = false
				}

				if strings.HasPrefix(sub, " ") {
					splitSubs[i] = sub[1:]
					splitAtWord = false
				}

				if splitAtWord {
					splitSubs[i] += "-"
				}
			}
			srt := formatter(subStateStart, endTime, strings.Join(splitSubs, "\n"))
			data.WriteString(srt)
			subStateCount = 0
			subStateStart = -1
			subStateSubs = ""
		}
	}

	return data.String(), nil
}
