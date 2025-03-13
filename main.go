package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"example.com/protobuf-test/personpb"
	"google.golang.org/protobuf/proto"
)

// Định nghĩa struct cho JSON & Go Binary (gob)
type Person struct {
	Name  string
	Age   int
	Email string
}

// Tạo danh sách dữ liệu mẫu
func generatePeople(n int) []Person {
	people := make([]Person, n)
	for i := 0; i < n; i++ {
		people[i] = Person{
			Name:  fmt.Sprintf("User%d", i),
			Age:   rand.Intn(100),
			Email: fmt.Sprintf("user%d@example.com", i),
		}
	}
	return people
}

// Hàm đo thời gian thực thi
func measureTime(f func()) time.Duration {
	start := time.Now()
	f()
	return time.Since(start)
}

func main() {
	// 🔥 Tạo danh sách 1 triệu đối tượng
	numRecords := 1000000
	people := generatePeople(numRecords)

	// 🟢 JSON
	jsonTime := measureTime(func() {
		_, _ = json.Marshal(people)
	})
	jsonData, _ := json.Marshal(people)

	// 🟢 Go Binary (gob)
	var gobData bytes.Buffer
	gobEncoder := gob.NewEncoder(&gobData)
	gobTime := measureTime(func() {
		gobData.Reset()
		err := gobEncoder.Encode(people)
		if err != nil {
			log.Fatal("Gob encoding failed:", err)
		}
	})
	gobSize := gobData.Len()

	// 🟢 Protocol Buffers
	protoPeople := make([]*personpb.Person, numRecords)
	for i, p := range people {
		protoPeople[i] = &personpb.Person{
			Name:  p.Name,
			Age:   int32(p.Age),
			Email: p.Email,
		}
	}
	protoData, _ := proto.Marshal(&personpb.People{Persons: protoPeople})
	protoTime := measureTime(func() {
		_, _ = proto.Marshal(&personpb.People{Persons: protoPeople})
	})

	// Kết quả so sánh
	fmt.Println("JSON Size: ", len(jsonData), "bytes, Time: ", jsonTime)
	fmt.Println("Go Binary (gob) Size: ", gobSize, "bytes, Time: ", gobTime)
	fmt.Println("Protocol Buffers Size: ", len(protoData), "bytes, Time: ", protoTime)
}
