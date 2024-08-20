package pkg

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/term"
	"io"
	"math"
	"os"
	"os/signal"
	"strings"
	"time"
)

// 从文件读取
func ReadFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var text string
	for scanner.Scan() {
		text += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return text, nil
}

// 从标准输入读取
func ReadFromStdin() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	var text string
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", err
		}
		text += line
		if err == io.EOF {
			break
		}
	}
	return text, nil
}

// 是否为终端
//func isInteractiveTerminal() bool {
//	return os.Stdin.Fd() && os.Stdout && os.Stderr &&
//		syscall.SYS_IOC
//	os.Stdin.
//}

// isTerminal 检查给定的文件描述符是否是终端
func isTerminal(fd *os.File) bool {
	return fd != nil && fd.Fd() >= 0 && term.IsTerminal(int(fd.Fd()))
}

func DetectTerminal(writeMedia string) {
	if isTerminal(os.Stdin) && isTerminal(os.Stdout) && writeMedia == "" {
		_, err := fmt.Fprintln(os.Stderr, "Warning: TTS output will be written to the terminal. Use --write-media to write to a file.\nPress Ctrl+C to cancel the operation.\nPress Enter to continue.")
		if err != nil {
			panic("Error reading input:")
		}
		// 等待用户输入
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			panic("Error reading input:")
		}
	}
	// 设置中断信号处理
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Println(os.Stderr, "\nOperation canceled.")
			os.Exit(1)
		}
	}()
}

func splitTextByByteLength(text interface{}, byteLength int) ([]string, error) {
	var b []byte

	switch v := text.(type) {
	case string:
		b = []byte(v)
	case []byte:
		b = v
	default:
		return nil, fmt.Errorf("text must be string or []byte")
	}

	if byteLength <= 0 {
		return nil, fmt.Errorf("byteLength must be greater than 0")
	}

	var result []string

	for len(b) > byteLength {
		// Find the last space in the string
		splitAt := bytes.LastIndex(b[:byteLength], []byte(" "))

		// If no space found, splitAt is byteLength
		if splitAt == -1 {
			splitAt = byteLength
		}

		// Verify all & are terminated with a ;
		for {
			ampersandIndex := bytes.LastIndex(b[:splitAt], []byte("&"))
			if ampersandIndex == -1 {
				break
			}

			if bytes.IndexByte(b[ampersandIndex:], ';') != -1 {
				break
			}

			splitAt = ampersandIndex - 1
			if splitAt < 0 {
				return nil, fmt.Errorf("maximum byte length is too small or invalid text")
			}
			if splitAt == 0 {
				break
			}
		}

		// Append the string to the list
		newText := strings.Trim(string(b[:splitAt]), " ")
		if newText != "" {
			result = append(result, newText)
		}
		if splitAt == 0 {
			splitAt = 1
		}
		b = b[splitAt:]
	}

	newText := strings.Trim(string(b), " ")
	if newText != "" {
		result = append(result, newText)
	}

	return result, nil
}

func RemoveIncompatibleCharacters(s interface{}) string {
	var str string
	switch v := s.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return ""
	}

	var result strings.Builder
	for _, char := range str {
		code := int(char)
		if (0 <= code && code <= 8) || (11 <= code && code <= 12) || (14 <= code && code <= 31) {
			result.WriteRune(' ')
		} else {
			result.WriteRune(char)
		}
	}

	return result.String()
}

// Returns the headers and data to be used in the request
func ssmlHeadersPlusData(requestId string, timestamp string, ssml string) string {
	return fmt.Sprintf(
		"X-RequestId:%s\r\n"+
			"Content-Type:application/ssml+xml\r\n"+
			"X-Timestamp:%sZ\r\n"+ // This is not a mistake, Microsoft Edge bug.
			"Path:ssml\r\n\r\n"+
			"%s",
		requestId,
		timestamp,
		ssml,
	)
}

func getHeadersAndData(data interface{}) (map[string][]byte, []byte) {
	var byteData []byte
	switch v := data.(type) {
	case string:
		byteData = []byte(v)
	case []byte:
		byteData = v
	default:
		return nil, nil
	}

	headerEnd := bytes.Index(byteData, []byte("\r\n\r\n"))
	if headerEnd == -1 {
		return nil, nil
	}

	headers := make(map[string][]byte)
	headerLines := bytes.Split(byteData[:headerEnd], []byte("\r\n"))
	for _, line := range headerLines {
		parts := bytes.SplitN(line, []byte(":"), 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(string(parts[0]))
			value := bytes.TrimSpace(parts[1])
			headers[key] = value
		}
	}

	return headers, byteData[headerEnd+4:]
}

func calcMaxMsgSize(voice, rate, volume, pitch string) int {
	websocketMaxSize := int(math.Pow(2, 16))
	overheadPerMessage := len(ssmlHeadersPlusData(ConnectID(), dateToString(), Mkssml(voice, rate, volume, pitch, []byte("")))) + 50 // margin of error
	return websocketMaxSize - overheadPerMessage
}

// Returns a UUID without dashes
func ConnectID() string {
	uuidWithDashes := uuid.New().String()
	return strings.ReplaceAll(uuidWithDashes, "-", "")
}

// dateToString returns a Javascript-style date string.
func dateToString() string {
	return time.Now().UTC().Format("Mon Jan _2 2006 15:04:05 GMT+0000 (Coordinated Universal Time)")
}

func Mkssml(voice string, rate string, volume string, pitch string, escapedText []byte) string {
	text := string(escapedText)

	ssml := strings.Builder{}
	ssml.WriteString("<speak version='1.0' xmlns='http://www.w3.org/2001/10/synthesis' xml:lang='en-US'>")
	ssml.WriteString(fmt.Sprintf("<voice name='%s'>", voice))
	ssml.WriteString(fmt.Sprintf("<prosody pitch='%s' rate='%s' volume='%s'>", pitch, rate, volume))
	ssml.WriteString(text)
	ssml.WriteString("</prosody>")
	ssml.WriteString("</voice>")
	ssml.WriteString("</speak>")

	return ssml.String()
}

// escape escapes characters in the text.
// escapeHTML escapes HTML entities in the input string.
func escape(data string) string {
	data = strings.ReplaceAll(data, "&", "&amp;")
	data = strings.ReplaceAll(data, ">", "&gt;")
	data = strings.ReplaceAll(data, "<", "&lt;")
	return data
}

// unescapeHTML unescapes HTML entities in the input string.
func unescape(data string) string {
	data = strings.ReplaceAll(data, "&lt", "<")
	data = strings.ReplaceAll(data, "&gt", ">")
	data = strings.ReplaceAll(data, "&amp", "&")
	return data
}
