package test

import (
	"bytes"
	crand "crypto/rand"
	"fmt"
	"framework/security"
	"framework/utils"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime/debug"
	"runtime/trace"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

// func TestHeap(t *testing.T) {
// 	h := &IntHeap{2, 1, 5, 6, 4, 3, 7, 9, 8, 0} // 创建slice
// 	heap.Init(h)                                // 初始化heap
// 	t.Log(*h)
// 	t.Log(heap.Pop(h)) // 调用pop
// 	heap.Push(h, 6)    // 调用push
// 	t.Log(*h)
// 	for len(*h) > 0 {
// 		t.Logf("%d ", heap.Pop(h))
// 	}
// 	//t.Error("")
// }

func TestSwith(t *testing.T) {
	for index := 1; index < 100; index++ {
		switch {
		case index%15 == 0:
			t.Log("fizzbuzz")
		case index%3 == 0:
			t.Log("fizz")
		case index%5 == 0:
			t.Log("buzz")
		default:
			t.Log(index)
		}
	}
}
func TestStringBuffer(t *testing.T) {
	var buffer bytes.Buffer
	for index := 0; index < 99; index++ {
		buffer.WriteString(strconv.Itoa(index))
	}
	s := buffer.String()
	fmt.Printf(s)
}
func TestFunc(t *testing.T) {
	variadic(1, 2, 3)
}
func variadic(numbers ...int) {
	fmt.Printf("Type: %T\t Content: %d\n", numbers, numbers)
}

func TestHttp(t *testing.T) {
	res, err := http.Get("http://www.richengke.com")
	if err != nil {
		fmt.Printf("http error:%s", err)
		panic("error")
	} else {
		defer res.Body.Close()
		contents, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("http read error:%s", err)
			os.Exit(1)
		} else {
			fmt.Printf("contents %s\n", string(contents))
		}
	}
}

func FuncWithInterface(emptyinterface interface{}) {
	switch t := emptyinterface.(type) {
	case string:
		fmt.Print("type: string\t")
	case int:
		fmt.Print("type: int\t")
	case bool:
		fmt.Print("type: bool\t")
	case float64:
		fmt.Print("type: float64\t")
	default:
		fmt.Printf("type: %v\t", t)
	}

	fmt.Printf("data: %#v\n", emptyinterface)
}
func TestInterface(t *testing.T) {
	var interfaces = [3]interface{}{}
	interfaces[0], interfaces[1], interfaces[2] = 1, 1.1, "goo"

	fmt.Printf("interfaces %v\n", interfaces)

	for _, m := range interfaces {
		FuncWithInterface(m)
	}
}

func TestRand(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("math/rand:")
	for i := 0; i < 10; i++ {
		fmt.Println(i, rand.Intn(127))
	}

	fmt.Println("crypto/rand:")
	b := make([]byte, 3)
	for i := 0; i < 10; i++ {
		crand.Read(b)
		number := uint(b[0]) | uint(b[1])<<8 | uint(b[2])<<16
		fmt.Println(i, number)
	}
}

func TestSort(t *testing.T) {
	strs := []string{"c", "a", "b"}
	sort.Strings(strs)
	fmt.Println("strings:", strs)

	ints := []int{2, 1, 5}
	sort.Ints(ints)
	fmt.Println("ints:", ints)

	s := sort.IntsAreSorted(ints)
	fmt.Println("Sorted: ", s)
}

func TestSpilit(t *testing.T) {
	variable := "Lorem Ipsum Dolor Sit Amet"
	strs := strings.Split(variable, " ")
	for _, s := range strs {
		fmt.Println(s)
	}
	var b bytes.Buffer
	trace.Start(&b)
	log.Printf("log inside a function")
	//debug.PrintStack()
	log.Printf("%v\n", string(debug.Stack()))

	log.SetFlags(log.Ltime | log.Lshortfile)

	log.Println("first log output")
	log.Printf("second log output\n")
	//fn()
	log.Printf("trace: %v\n", b.String())
}
func TestChan(t *testing.T) {
	chanFinish := make(chan bool, 3)
	defer close(chanFinish)
	for i := 0; i < 3; i++ {
		go func(n int) {
			fmt.Println(n)
			chanFinish <- true
		}(i)
	}
	for i := 0; i < 3; i++ {
		b := <-chanFinish
		fmt.Println(b)
	}
}

func TestFile(t *testing.T) {
	filename := "./main_test.go"
	finfo, err := os.Stat(filename)
	if err != nil {
		fmt.Println(filename, "does not exist.")
	} else {
		if finfo.IsDir() {
			fmt.Println(filename, "is a dir.")
		} else {
			fileContent, err := ioutil.ReadFile(filename)
			if err != nil {
				panic(err)
			}
			fstr := string(fileContent)
			fmt.Print(fstr)
		}
	}
}
func TestFileBuffer(t *testing.T) {
	filename := "./main_test.go"
	fs, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fs.Close()

	filename2 := "./main_test.go.bak"
	fs2, _ := os.Create(filename2)
	defer fs2.Close()

	buf := make([]byte, 1024)
	for {
		n, err := fs.Read(buf)
		if err == nil {
			fs2.Write(buf[:n])
		} else if err != io.EOF {
			panic(err)
		} else if err == io.EOF {
			break
		}
	}
	fmt.Println(filename2, "created!")
}

func callAnother(f func(i int) int, x int) string {
	return strconv.Itoa(f(x * 2))
}
func Test_CallBack(t *testing.T) {
	i := 0
	callAnother(func(n int) int {
		i = n * 2
		return i
	}, 2)
	fmt.Println(i)
}
func initTimeSeq() func() int {
	t := time.Now().UnixNano()
	return func() int {
		return int(time.Now().UnixNano() - t)
	}
}

func Test_FunctionClosure(t *testing.T) {
	timeSince := initTimeSeq()

	fmt.Println(timeSince())

	fmt.Println(timeSince())

	time.Sleep(1 * time.Second)

	fmt.Println(timeSince())

	time.Sleep(120 * time.Millisecond)

	fmt.Println(timeSince())

	timeSince = initTimeSeq()

	time.Sleep(1300 * time.Millisecond)

	fmt.Println(timeSince())

	fmt.Println(timeSince())

	fmt.Println(timeSince())
}

func Test_InArray(t *testing.T) {
	ints := []int{2, 4, 6}
	tofind := 2
	b, i := utils.InArray(tofind, ints)
	fmt.Println(b, i)
}

func Test_Map(t *testing.T) {
	m := map[string]int{}
	m["a"], m["b"] = 1, 2
	fmt.Printf("%v\n", m)
}
func Test_Aes(t *testing.T) {
	key := "0123456789abcdef"
	text := "hello world中国"
	var err error
	if s1, err := security.AesEncryptBase64(text, key); err == nil {
		fmt.Printf("%v\n", s1)
		if s2, err := security.AesDecryptBase64(s1, key); err == nil {
			fmt.Printf("%v\n", s2)
		}
	}
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
