package main

import (
	"command/src/utils"
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type ConfigType struct {
	ExcludeDirs []string `yaml:"excludeDirs"` // 排除目录
	FileExts    []string `yaml:"fileExts"`    // 直接算做注释的文件
	Language    map[string]struct {
		Multi  []string `yaml:"multi"`
		Single []string `yaml:"single"`
	} `yaml:"language"`
}

var config ConfigType
var template = `excludeDirs: [".git"]
fileExts: ["md"]
language:
  go:
    single: ["//[^\\n]+\\n?"]
  ts:
    multi: ["(?s)/\\*.*?\\*/"]
    single: ["//[^\\n]+\\n?"]
`

func init() {
	homeDir, _ := os.UserHomeDir()
	configFile := homeDir + "/command/codestat.yml"
	// configFile := "codestat.yml"
	data, err := os.ReadFile(configFile)
	if err != nil {
		os.Create(configFile)
		os.WriteFile(configFile, []byte(template), 0777)
		data, _ = os.ReadFile(configFile)
	}

	if err := yaml.Unmarshal([]byte(data), &config); err != nil {
		panic(err)
	}
}

func main() {
	type Result struct {
		TotalSize  int
		MultiSize  int
		SingleSize int
		TotalLine  int
		MultiLine  int
		SingleLine int
		Ext        string
	}
	fileCollect := map[string]Result{}

	filename := "./"      // 默认当前目录
	if len(os.Args) > 1 { // 指定文件夹
		filename = os.Args[1]
	}
	isDetail := false
	if len(os.Args) > 2 && os.Args[2] == "-d" { // 是否显示每个文件的注释行数
		isDetail = true
	}

	// 计算注释行数
	calculation := func(file *utils.FileItem) {
		if file.Ext == "" {
			return
		}
		data, err := os.ReadFile(file.Path)
		if err != nil {
			panic(err)
		}

		//直接算作注释的文件
		if utils.SliceIncluded(config.FileExts, file.Ext[1:]) {
			content := string(data)
			fileCollect[file.Path] = Result{
				TotalSize: len(content),
				MultiSize: len(content),
				TotalLine: len(strings.Split(content, "\n")),
				MultiLine: len(strings.Split(content, "\n")),
				Ext:       file.Ext,
			}
			return
		}

		// 代码行数统计
		content := string(data)
		totalLine := 0
		strs := strings.Split(content, "\n")
		for _, str := range strs {
			if strings.TrimSpace(str) == "" { // 空行不统计
				continue
			}
			totalLine++
		}

		val := config.Language[file.Ext[1:]]
		if len(val.Multi) == 0 && len(val.Single) == 0 { // 没配置任何信息
			return
		}

		// 多行注释统计
		multiCount := 0
		multiLine := 0
		for _, multi := range val.Multi {
			reg := regexp.MustCompile(multi)
			matches := reg.FindAllStringSubmatch(content, -1)
			for _, match := range matches {
				multiCount += len(match[0])
				multiLine += len(strings.Split(match[0], "\n"))
				content = strings.Replace(content, match[0], "", 1)
			}
		}

		// 单行注释统计
		singleCount := 0
		singleLine := 0
		for _, single := range val.Single {
			reg := regexp.MustCompile(single)
			matches := reg.FindAllStringSubmatch(content, -1)
			for _, match := range matches {
				singleCount += len(match[0])
				singleLine++
			}
		}

		// 文件收集
		fileCollect[file.Path] = Result{
			TotalSize:  len(string(data)),
			MultiSize:  multiCount,
			SingleSize: singleCount,
			TotalLine:  totalLine,
			MultiLine:  multiLine,
			SingleLine: singleLine,
			Ext:        file.Ext,
		}
	}

	if f, _ := os.Stat(filename); f.IsDir() {
		utils.FileCatalog(filename, func(file *utils.FileItem) bool {
			if utils.SliceIncluded(config.ExcludeDirs, file.Name) {
				return false
			}
			calculation(file)
			return true
		})
	} else {
		file, _ := utils.FileInfo(filename)
		calculation(&file)
	}

	type ExtType struct {
		Code   int
		Multi  int
		Single int
		Size   int
		Files  int
	}
	extCollect := map[string]ExtType{}
	stat := Result{}
	if isDetail {
		fmt.Printf("%s\t%s\t%s\t%s\n", "Total", "Multi", "Single", "Path")
	}
	for path, item := range fileCollect {
		stat.TotalSize += item.TotalSize
		stat.MultiSize += item.MultiSize
		stat.SingleSize += item.SingleSize
		stat.TotalLine += item.TotalLine
		stat.MultiLine += item.MultiLine
		stat.SingleLine += item.SingleLine

		// fmt.Printf("%d\t%d\t%d\t%s\n", item.TotalSize, item.MultiSize, item.SingleSize, path)
		if isDetail {
			fmt.Printf("%d\t%d\t%d\t%s\n", item.TotalLine, item.MultiLine, item.SingleLine, path)
		}

		extItem := extCollect[item.Ext]
		extCollect[item.Ext] = ExtType{
			Code:   extItem.Code + item.TotalLine,
			Multi:  extItem.Multi + item.MultiLine,
			Single: extItem.Single + item.SingleLine,
			Size:   extItem.Size + item.TotalSize,
			Files:  extItem.Files + 1,
		}
	}
	fmt.Println("\nExt\tFiles\tCode\tMulti\tSingle\tByte")
	for ext, item := range extCollect {
		fmt.Printf("%s\t%d\t%d\t%d\t%d\t%d\n", ext, item.Files, item.Code, item.Multi, item.Single, item.Size)
	}

	percent := fmt.Sprintf("%.2f", float64(stat.MultiLine+stat.SingleLine)/float64(stat.TotalLine)*100) + "%"
	fmt.Println("\n统计：", stat.TotalLine, "行，", stat.TotalSize, "byte")
	fmt.Println("注释：", stat.MultiLine+stat.SingleLine, "行，", stat.MultiSize+stat.SingleSize, "byte", "\t占比："+percent)
}
