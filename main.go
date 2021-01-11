package main

import (
	"flag"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type item struct {
	name     string
	voteTime time.Time
	number   int
}

var lock sync.Mutex
var votes = map[string]item{}

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":80", "ip:port")
	flag.Parse()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(home))
	})
	mux.HandleFunc("/avg", func(writer http.ResponseWriter, request *http.Request) {
		rows := strings.Builder{}
		lock.Lock()
		defer lock.Unlock()
		var sum int
		var temp []item
		for _, i := range votes {
			temp = append(temp, i)
			sum += i.number
		}
		avg := float64(sum) / float64(len(votes))
		sort.Slice(temp, func(i, j int) bool {
			a := math.Abs(float64(temp[i].number) - avg/2)
			b := math.Abs(float64(temp[j].number) - avg/2)
			if a == b {
				return temp[i].voteTime.Before(temp[j].voteTime)
			}
			return a < b
		})
		for i, it := range temp {
			rows.WriteString(fmt.Sprintf("<tr>\n<td>%d</td>\n<td>%s</td>\n<td>%d</td>\n<td>%s</td>\n</tr>", i, it.name, it.number, it.voteTime.Format("2006-01-02 15:04:05.999999999")))
		}
		var output = fmt.Sprintf("%s%s%s<br>总计：%d, 参与人数：%d, 平均数为%f, 平均数一半为%f", tableHeader, rows.String(), tableTail, sum, len(votes), avg, avg/2)
		_, _ = writer.Write([]byte(output))
	})

	mux.HandleFunc("/add", func(writer http.ResponseWriter, request *http.Request) {
		name := request.PostFormValue("name")
		num, err := strconv.Atoi(request.PostFormValue("number"))
		lock.Lock()
		defer lock.Unlock()
		if name == "" {
			_, _ = writer.Write([]byte("姓名不能为空"))
			return
		}
		if _, ok := votes[name]; ok {
			_, _ = writer.Write([]byte(name + "，请勿重复提交"))
			return
		}
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		votes[name] = item{
			voteTime: time.Now(),
			number:   num,
			name:     name,
		}
		_, _ = writer.Write([]byte("参与成功,等待公布答案！"))
	})
	mux.HandleFunc("/clear", func(writer http.ResponseWriter, request *http.Request) {
		lock.Lock()
		defer lock.Unlock()
		votes = map[string]item{}
		_, _ = writer.Write([]byte("ok"))
	})
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		fmt.Println(err)
	}
}
