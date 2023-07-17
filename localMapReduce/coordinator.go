package mapReduce

import (
	"fmt"
)

func Coordinator() {
	fmt.Println("In coordinator")

	Mapper()
}
