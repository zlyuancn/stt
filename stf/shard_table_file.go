package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	TableRe = regexp.MustCompile(`[\\/]?(\w+)(.alter)?\.sql`) // 根据文件名提取表名
)

func main() {
	shardNums := flag.Int("c", 2, "分表数量")
	startIndex := flag.Int("i", 0, "第一个分表的编号")
	outFileNameSuffix := flag.String("o", ".out.sql", "输出文件名后缀")
	tables := flag.String("t", "", "要生成的表模板文件用半角逗号分隔")
	flag.Parse()

	if *tables == "" {
		fmt.Println("使用-t指定要生成的表")
		return
	}
	fmt.Println("加载模板")
	tableTemplate := loadShardTableTemplate(strings.Split(*tables, ","))

	// 遍历模板
	for tableFile, template := range tableTemplate {
		fileName := tableFile
		// 去掉扩展名
		if ext := filepath.Ext(tableFile); ext != "" {
			fileName = fileName[:len(fileName)-len(ext)]
		}
		fileName += *outFileNameSuffix
		fmt.Println("开始生成 ", fileName)

		outFile, err := os.Create(fileName)
		if err != nil {
			panic(fmt.Errorf("创建输出文件失败, fileName: %v, err: %v", fileName, err))
		}
		defer outFile.Close()

		for shard := 0; shard < *shardNums; shard++ {
			renderTemplate(tableFile, template, shard+(*startIndex), outFile)
		}
	}
	fmt.Println("完成")
}

// 渲染模板
func renderTemplate(tableFile string, template []byte, shard int, outFile io.Writer) {
	ss := TableRe.FindStringSubmatch(tableFile)
	if len(ss) == 0 {
		panic(fmt.Errorf("无法根据模板文件名解析出表名, %s", tableFile))
	}
	oldText := ss[1]
	newText := oldText + strconv.Itoa(shard)
	fmt.Println(newText)

	out := bytes.ReplaceAll(template, []byte(oldText), []byte(newText))
	_, err := outFile.Write(out)
	if err != nil {
		panic(fmt.Errorf("写入模板数据时失败, tableFile: %v, shard: %v, err: %v", tableFile, shard, err))
	}
	_, err = outFile.Write([]byte("\n\n"))
	if err != nil {
		panic(fmt.Errorf("写入模板数据时失败, tableFile: %v, shard: %v, err: %v", tableFile, shard, err))
	}
}

// 加载模板文件
func loadShardTableTemplate(tables []string) map[string][]byte {
	ret := make(map[string][]byte, len(tables))
	for _, tableFile := range tables {
		data, err := os.ReadFile(tableFile)
		if err != nil {
			panic(fmt.Errorf("加载模板文件数据失败, tableFile: %v, err: %v", tableFile, err))
		}
		ret[tableFile] = data
	}
	return ret
}
