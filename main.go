package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

var Cho = []string{"ㄱ", "ㄲ", "ㄴ", "ㄷ", "ㄸ", "ㄹ", "ㅁ", "ㅂ", "ㅃ", "ㅅ", "ㅆ", "ㅇ", "ㅈ", "ㅉ", "ㅊ", "ㅋ", "ㅌ", "ㅍ", "ㅎ"}
var Joong = []string{"ㅏ", "ㅐ", "ㅑ", "ㅒ", "ㅓ", "ㅔ", "ㅕ", "ㅖ", "ㅗ", "ㅘ", "ㅙ", "ㅚ", "ㅛ", "ㅜ", "ㅝ", "ㅞ", "ㅟ", "ㅠ", "ㅡ", "ㅢ", "ㅣ"}
var Jong = []string{"", "ㄱ", "ㄲ", "ㄳ", "ㄴ", "ㄵ", "ㄶ", "ㄷ", "ㄹ", "ㄺ", "ㄻ", "ㄼ", "ㄽ", "ㄾ", "ㄿ", "ㅀ", "ㅁ", "ㅂ", "ㅄ", "ㅅ", "ㅆ", "ㅇ", "ㅈ", "ㅊ", "ㅋ", "ㅌ", "ㅍ", "ㅎ"}

var DestCho = []string{"ㄱ", "ㄴ", "ㄷ", "ㄹ", "ㅁ", "ㅂ", "ㅅ", "ㅇ", "ㅈ", "ㅊ", "ㅋ", "ㅌ", "ㅍ", "ㅎ", "ㄲ", "ㄸ", "ㅃ", "ㅆ", "ㅉ"}
var DestJoong = []string{"ㅏ", "ㅑ", "ㅓ", "ㅕ", "ㅗ", "ㅛ", "ㅜ", "ㅠ", "ㅡ", "ㅣ", "ㅐ", "ㅔ", "ㅒ", "ㅖ", "ㅘ", "ㅙ", "ㅚ", "ㅝ", "ㅞ", "ㅟ", "ㅢ"}
var DestJong = []string{"ㄱ", "ㄴ", "ㄷ", "ㄹ", "ㅁ", "ㅂ", "ㅅ", "ㅇ", "ㅈ", "ㅊ", "ㅋ", "ㅌ", "ㅍ", "ㅎ", "ㄲ", "ㅆ", "ㄳ", "ㄵ", "ㄶ", "ㄺ", "ㄻ", "ㄼ", "ㄽ", "ㄾ", "ㄿ", "ㅀ", "ㅄ"}

var sources = [][]string{Cho, Joong, Jong}
var destinations = [][]string{DestCho, DestJoong, DestJong}

func main() {

	fmt.Println("\n정보 입력으 시작합니다.")

	scanner := bufio.NewScanner(os.Stdin)

	res := "\n===== 다음과 같이 동작해주세요 =====\n\n" +
		"1. PGM 키를 모드 스위치에 꼽고 모드를 PGM 으로 할당하세요\n" +
		"2. 아래에 맞추어 버튼을 입력하세요.\n(5, 6번째 줄은 공백일 경우 무시됩니다.)\n\n" +
		"2 - 소계 (화면 좌측에 PGM2가 나타납니다.)\n"

	for i := 1; i < 7; i++ {
		fmt.Printf("%d번째 줄 입력: ", i)
		scanner.Scan()
		inputText := scanner.Text()

		if i >= 5 && len(inputText) == 0 {
			continue
		}

		t := strconv.Itoa(i) + "32 - 소계 - "
		charSet := ""

		for _, runeChar := range inputText {
			charSet += translate(runeChar) + "*"
		}

		if len(charSet) > 1 && charSet[0:] == "*" {
			charSet = charSet[1:]
		}

		for j := range charSet {
			if j > 0 && charSet[j] == '*' && charSet[j-1] == '*' {
				charSet = charSet[0:j] + charSet[j+1:]
			}
		}

		t += charSet + " - 마감\n"
		res += t
	}

	res += "소계\n\n===== END ======"

	fmt.Println(res)
}

func translate(r rune) string {
	if isKorean(r) {
		return translateKorean(r)
	}

	if r == 34 || r < 32 || r > 127 {
		return ""
	}

	return buildIndexNum(int(r))
}

func translateKorean(r rune) string {
	// if syllable
	if r >= 0xAC00 && r <= 0xD7A3 {
		return resolveSyllableKor(r)
	}
	// all other cases
	return resolveSingleKor(r)
}

func resolveSingleKor(r rune) string {
	baseCode, err := getKoreanBaseCode(r)

	if err != nil || baseCode != 0x3130 {
		return ""
	}

	s := string(r)

	for _, val := range []int{findIndex(s, DestCho), findIndex(s, DestJoong), findIndex(s, DestJong)} {
		if val > 0 {
			return buildIndexNum(val + 1)
		}
	}

	log.Fatalf("failed parsing... %c\n", r)
	return ""
}

func buildIndexNum(i int) string {
	if i == 0 {
		return ""
	} else if i < 10 {
		return "0" + strconv.Itoa(i)
	} else {
		return strconv.Itoa(i)
	}
}

func resolveSyllableKor(r rune) string {
	offset := (int)(r - 0xAC00)

	c, j1, j2 := offset/28/21, offset/28%21, offset%28

	c = findIndex(Cho[c], DestCho)
	j1 = findIndex(Joong[j1], DestJoong)
	j2 = findIndex(Jong[j2], DestJong)

	s := ""

	for _, v := range []int{c, j1, j2} {
		s += buildIndexNum(v + 1)
	}

	return s
}

func findIndex(letter string, arr []string) int {
	for i, v := range arr {
		if letter == v {
			return i
		}
	}
	return -1
}

func isKorean(r rune) bool {
	return r >= 0xAC00 && r < 0xD7A4 ||
		r >= 0x1100 && r < 0x1200 ||
		r >= 0x3130 && r < 0x3190 ||
		r >= 0xA960 && r < 0xA980 ||
		r >= 0xD7B0 && r < 0xD800
}

func getKoreanBaseCode(r rune) (rune, error) {
	if !isKorean(r) {
		return -1, fmt.Errorf("given rune %c is not a korean letter\n", r)
	}
	if r >= 0xAC00 && r <= 0xD7A3 {
		return 0xAC00, nil
	} else if r >= 0x1100 && r <= 0x11FF {
		return 0x1100, nil
	} else if r >= 0x3130 && r <= 0x318F {
		return 0x3130, nil
	} else if r >= 0xA960 && r <= 0xA97F {
		return 0xA960, nil
	} else {
		return 0xD7B0, nil
	}
}
