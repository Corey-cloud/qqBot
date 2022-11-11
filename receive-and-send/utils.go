package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

//从文件中读取成语和其释义
func getWordsFromFile() map[string]string {
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
func getBeginWord(words map[string]string) string {
	i := 0
	//map无序，随机取第10个，每次结果不同
	for word := range words {
		if i == 9 {
			return word
		}
		i += 1
	}
	return ""
}

//根据用户回复获取接龙成语
func getWord(word string, words map[string]string) string {
	lastStr := string([]rune(word)[3:])
	for word, _ := range words {
		if string([]rune(word)[:1]) == lastStr {
			return word
		}
	}
	return ""
}

//判断成语是否合法
func isWordLegal(word string, words map[string]string) bool {
	_, ok := words[word]
	return len([]rune(word)) == 4 && ok == true
}

//判断是否接龙成功
func isWordDragon(word string, preWord string) bool {
	if word != "" && preWord != "" {
		flag := string([]rune(word)[:1]) == string([]rune(preWord)[3:])
		return flag
	}
	return false
}

//获取成语释义
func getWordMeaning(word string, words map[string]string) string {
	return words[word]
}
