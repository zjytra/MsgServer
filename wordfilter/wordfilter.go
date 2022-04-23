
package wordfilter

import (
    "io/ioutil"
    "os"
    "path"
    "strings"
)

var (
    ConfExample *ConfigFilter
)

//屏蔽词节点
type FilterNode struct {
    NodeStr rune                  //内容
    Subli   map[rune]*FilterNode //屏蔽子集合
    IsEnd   bool                  //是否为结束
}

type ConfigFilter struct {
    FilterList map[rune]*FilterNode //屏蔽字树
}
//加载词库
func InitConfigFilter(configpath string) {
    ConfExample = new(ConfigFilter)
    ConfExample.FilterList = make(map[rune]*FilterNode)
    //我这里用的是一个文本文件，一行表示一个屏蔽词
    file, err := os.Open(path.Join(configpath, "list.txt"))
    if err != nil {
        panic(err)
    }
    barr, _ := ioutil.ReadAll(file)
    bstr := string(barr)
    bstr = strings.ReplaceAll(bstr, "\r", "")
    rows := strings.Split(bstr, "\n")
    for _, row := range rows {
        rowr := []rune(row)
        wordLen := len(rowr)
        node := ConfExample.FilterList //先从根map来
        for i := 0; i < wordLen; i++ {
            subNode, subok := node[rowr[i]] //第一个取根
            if !subok {//没有找到就添加
                subNode = createNode(node,rowr[i])
            }
            subNode.IsEnd = wordLen -1 == i
            node = subNode.Subli //又再下一阶寻找
        }
    }
}

func createNode(fnode map[rune]*FilterNode,char rune) *FilterNode {
    fmd,ok  := fnode[char]
    if !ok { //没有找到的时候才创建
        fmd = new(FilterNode)
        fmd.NodeStr = char
        fmd.Subli = make(map[rune]*FilterNode)
        fnode[char] = fmd
    }
    return fmd
}
 


 
//屏蔽字操作，这个方法就是外部调用的入口方法
func FilterChack(data string) string{
    filterli := ConfExample.FilterList
    arr := []rune(data)
    wordLen := len(arr)
    //逐个遍历输入字符串
    for i := 0; i < wordLen; i++ {
        if rune(' ') == arr[i]  { //遇见空格又从根来
            filterli = ConfExample.FilterList
            continue
        }
        fmd, ok := filterli[arr[i]]
        if !ok { //没有找到没有在屏蔽词中
            filterli = ConfExample.FilterList
            //再从根再找一次 针对中间断掉的屏蔽字
            fmd, ok = filterli[arr[i]]
            if !ok {
                continue
            }
        }
        //找到了就替换
        arr[i] = rune('*')
        if fmd.IsEnd {
            filterli = ConfExample.FilterList
        }else {
            filterli = fmd.Subli //下一个在子树里面查找
        }

    }
    return string(arr)
}
 
