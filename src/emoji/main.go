package main

import (
    "fmt"
    "html"
    "strconv"
)

func main() {
    // 相關範圍可以參考這裡：http://apps.timwhitlock.info/emoji/tables/unicode
    emoji := [][]int{
        // 表情圖示的範圍。
        {128513, 128591},
        // 裝飾符號的範圍。
        {9986, 10160},
        // 交通工具還有地圖的範圍。
        {128640, 128704},
    }

    for _, value := range emoji {
        for x := value[0]; x < value[1]; x++ {
            // 將範圍轉成 Unicode 字串後反脫逸字元，這樣就會出現表情了。
            str := html.UnescapeString("&#" + strconv.Itoa(x) + ";")
            // 輸出表情。
            fmt.Println(str)
        }
    }
}