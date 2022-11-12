package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type WordsMap map[string]string

//从文件中读取成语和其释义
func getWordsFromFile() WordsMap {
	Map := make(map[string]string)
	f, err := os.Open("words.txt")
	if err != nil {
		fmt.Println("err", err)
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		words := strings.Split(line, "\t")
		Map[words[0]] = words[2]
	}
	return Map
}

//获取初始成语
func (ws WordsMap) getBeginWord() string {
	i := 0
	//map无序，随机取第10个，每次结果不同
	for word := range ws {
		if i == 9 {
			return word
		}
		i += 1
	}
	return ""
}

//根据用户回复获取接龙成语
func (ws WordsMap) getWord(word string) string {
	lastStr := string([]rune(word)[3:])
	for word, _ := range ws {
		if string([]rune(word)[:1]) == lastStr {
			return word
		}
	}
	return ""
}

//判断成语是否合法
func (ws WordsMap) isWordLegal(word string) bool {
	_, ok := ws[word]
	return len([]rune(word)) == 4 && ok == true
}

//判断是否接龙成功
func (ws WordsMap) isWordDragon(word string, preWord string) bool {
	if word != "" && preWord != "" {
		flag := string([]rune(word)[:1]) == string([]rune(preWord)[3:])
		return flag
	}
	return false
}

//获取成语释义
func (ws WordsMap) getWordMeaning(word string) string {
	return ws[word]
}
