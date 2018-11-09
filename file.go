package autodoc

import (
	"io/ioutil"
	"strings"
	"regexp"
)

var fileCache = map[string][]string{}

func readLine(filePath string, lineno int) string {
	lines, ok := fileCache[filePath]
	if !ok {
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic(err)
		}
		lines = strings.Split(string(data), "\n")
		fileCache[filePath] = lines
	}
	return lines[lineno-1]
}

// 读取调用处上方的注释(块)
// lineno 为注释块的最后一行
func readAboveComments(filePath string, lineno int) (comments []string) {
	for ; lineno > 0; lineno-- {
		line := readLine(filePath, lineno)
		comment := strings.TrimSpace(line)
		if len(comment) >= 2 && comment[:2] == "//" {
			comments = append(comments, strings.TrimSpace(comment[2:]))
		} else {
			break
		}
	}
	// reverse
	for i, j := 0, len(comments)-1; i < j; i, j = i+1, j-1 {
		comments[i], comments[j] = comments[j], comments[i]
	}
	return comments
}

var commentRegex, _ = regexp.Compile(`\) *//`)

func _readTailComment(s string) (comment string) {
	// 第一个匹配的结果
	idx := commentRegex.FindStringIndex(s)
	if idx == nil {
		return ""
	}
	return strings.TrimSpace(s[idx[1]:])
}

// 读取调用处尾部的注释
// 简单起见，只读取第一个 ) // 之后的内容（右括号和 // 之间有若干空格）
func readTailComment(filePath string, lineno int) (comment string) {
	line := readLine(filePath, lineno)
	return _readTailComment(line)
}
