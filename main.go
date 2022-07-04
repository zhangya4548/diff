package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	path := filepath.Dir(os.Args[0])
	zuoTianPath := fmt.Sprintf("%s/昨天.txt",path)
	jinTianPath := fmt.Sprintf("%s/今天.txt",path)
	newPath := fmt.Sprintf("%s/已去重处理.txt",path)
	if !IsExist(zuoTianPath) {
		fmt.Printf("昨天.txt 文件不存在 \n")
		return
	}
	if !IsExist(jinTianPath) {
		fmt.Printf("今天.txt 文件不存在 \n")
		return
	}

	if !IsExist(newPath) {
		_, err := FileCaeate(newPath, []byte{})
		if err != nil {
			fmt.Printf("创建文件异常: %s \n",err)
			return
		}
	}

	zuoTianReader, err := ioutil.ReadFile(zuoTianPath)
	if err != nil {
		fmt.Printf("读取昨天文件异常: %s \n",err)
		return
	}

	zuoTianList := SplitLines(string(zuoTianReader))
	fmt.Println("读取昨天到",len(zuoTianList),"个")

	jinTianReader, err := ioutil.ReadFile(jinTianPath)
	if err != nil {
		fmt.Printf("读取今天文件异常: %s \n",err)
		return
	}

	jinTianList := SplitLines(string(jinTianReader))
	fmt.Println("读取今天到",len(jinTianList),"个")

	tmpMap := make(map[string]struct{})
	for k, v := range zuoTianList {
		if v == "" {
			continue
		}

		if _, ok := tmpMap[v]; ok {
			continue
		}
		tmpMap[v] = struct{}{}
		zuoTianList[k] = Trim(v,"")
	}


	tmpMap2 := make(map[string]struct{})
	for k, v := range jinTianList {
		if v == "" {
			continue
		}
		if _, ok := tmpMap2[v]; ok {
			continue
		}
		tmpMap2[v] = struct{}{}
		jinTianList[k] = Trim(v,"")
	}

	listMap := make(map[string]struct{})

	_,_,tmpList := DiffStr(zuoTianList,jinTianList)
	for _, v := range tmpList {
		if _, ok := listMap[v]; ok {
			continue
		}
		listMap[v] = struct{}{}
	}

	fmt.Println("去重后",len(listMap),"个")

	//5. 打开一个已经有的文件，将原来的内容覆盖成新的内容十句"hello，Leslie"
	file, err := os.OpenFile(newPath, os.O_TRUNC|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Printf("打开文件异常: %s \n",err)
		return
	}

	writer := bufio.NewWriter(file)
	for k, _ := range listMap {
		writer.WriteString(k+`
`)
	}

	writer.Flush()
	file.Close()
	fmt.Println("处理完成")
}

func DiffStr(a, b []string) (inAAndB, inAButNotB, inBButNotA []string) {
	m := make(map[string]uint8)
	for _, k := range a {
		m[k] |= 1 << 0
	}
	for _, k := range b {
		m[k] |= 1 << 1
	}

	for k, v := range m {
		x := v&(1<<0) != 0
		y := v&(1<<1) != 0
		switch {
		case x && y:
			inAAndB = append(inAAndB, k)
		case x && !y:
			inAButNotB = append(inAButNotB, k)
		case !x && y:
			inBButNotA = append(inBButNotA, k)
		}
	}

	return
}
func FileCaeate(fileName string, opBytes []byte) (int, error) {
	file,err:=os.Create(fileName)
	if err!=nil{
		fmt.Println(err)
	}
	defer file.Close()
	//content := []byte("111")
	tag,err :=file.Write(opBytes)
	if err!=nil {
		fmt.Println(err)
	}
	return tag, err
}
func Trim(s, cutset string) string {
	if cutset == "" {
		return strings.TrimSpace(s)
	}
	return strings.Trim(s, cutset)
}
func IsExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || os.IsExist(err)
}
//按照换行切割字符串
func SplitLines(s string) []string {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}