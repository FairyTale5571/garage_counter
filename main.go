package main

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"regexp"
	"unsafe"
)
var dbConnected = false

//export goRVExtensionVersion
func goRVExtensionVersion(output *C.char, outputsize C.size_t) {
	result := C.CString("GRC 1.0")
	defer C.free(unsafe.Pointer(result))
	var size = C.strlen(result) + 1
	if size > outputsize {
		size = outputsize
	}
	C.memmove(unsafe.Pointer(output), unsafe.Pointer(result), size)
}

//export goRVExtensionArgs
func goRVExtensionArgs(output *C.char, outputsize C.size_t, input *C.char, argv **C.char, argc C.int) {
	//offset := unsafe.Sizeof(uintptr(0))
	action := C.GoString(input)
	clearArgs := cleanInput(argv, int(argc))

	switch action {

	case "count_vehicles":

		if !dbConnected {
			fmt.Printf("database not connected\n")
			printInArma(output, outputsize, "db not con")
			return
		}
		var vehicles []string
		for i:=1;i < len(clearArgs); i++ {
			vehicles = append(vehicles,clearArgs[i])
		}
		vehs := countVeh(clearArgs[0],vehicles)
		printInArma(output,outputsize,vehs)
		return

	default:
		temp := fmt.Sprintf("Undefined '%s' command", action)
		printInArma(output, outputsize, temp)
		return
	}
}


func struct2JSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func cleanInput(argv **C.char, argc int) []string {
	newArgs := make([]string, argc)
	offset := unsafe.Sizeof(uintptr(0))
	i := 0
	for i < argc {
		_arg := (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv)) + offset*uintptr(i)))
		arg := C.GoString(*_arg)
		arg = arg[1 : len(arg)-1]

		reArg := regexp.MustCompile(`""`)
		arg = reArg.ReplaceAllString(arg, `"`)

		newArgs[i] = arg
		i++
	}

	return newArgs
}

func printInArma(output *C.char, outputsize C.size_t, input string) {
	result := C.CString(input)
	defer C.free(unsafe.Pointer(result))
	var size = C.strlen(result) + 1
	if size > outputsize {
		size = outputsize
	}
	C.memmove(unsafe.Pointer(output), unsafe.Pointer(result), size)
}

//export goRVExtension
func goRVExtension(output *C.char, outputsize C.size_t, input *C.char) {
	id := returnMyData(C.GoString(input), nil)
	printInArma(output, outputsize, id)
}

func returnMyData(input string, errors error) string {
	switch input {
	case "connectDB":
		err := ReadConfig()
		if err != nil {
			return "config error : "+err.Error()
		}
		db, err = ConnectDatabase()
		if err != nil {
			return "non "+err.Error()
		}
		return "connected"
	default:
		return "undefined command"
	}
	return ""
}

func main() {}


