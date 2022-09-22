package main

import (
	"fmt"
	"strings"
)

func main() {
	var data string = "--name_class=QUENTIN+FIBREUX,ETIENNE+ALATIENNE,DIMITRI+PAYEEEEEET"

	if strings.Contains(data, "name_class") {
		dt := strings.Split(data, "=")
		dt = dt[1:]
		datal := strings.Join(dt, "=")
		datasplit := strings.Split(datal, ",")
		for _, value := range datasplit {
			datatype := strings.Split(value, "+")
			dataname := strings.Split(value, "+")
			datatype = datatype[:1]
			dataname = dataname[1:]
			fmt.Println("datatype", datatype)
			fmt.Println("dataname", dataname)
		}
	}
}
