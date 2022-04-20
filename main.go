package main

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
*/
import "C"

import (
	"fmt"
	"regexp"
	"strings"
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
	action := C.GoString(input)
	clearArgs := cleanInput(argv, int(argc))

	if !dbConnected {
		fmt.Printf("database not connected\n")
		printInArma(output, outputsize, "db not con")
		return
	}
	switch action {
	case "clear_old_veh":
		printInArma(output, outputsize, deleteOldVehicles(parseVehicleArray(0, clearArgs)))
		return

	case "count_vehicles":
		printInArma(output, outputsize, countVeh(clearArgs[0], parseVehicleArray(1, clearArgs)))
		return
	case "perf":
		fmt.Println(clearArgs)
		printInArma(output, outputsize, insertCpu(clearArgs[0], clearArgs[1], clearArgs[2], clearArgs[3], clearArgs[4]))
		return
	default:
		temp := fmt.Sprintf("Undefined '%s' command", action)
		printInArma(output, outputsize, temp)
		return
	}
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
	var vars []string
	var inputSplit = strings.Split(input, "~")
	var function = ""
	for idx := range inputSplit {
		if idx == 0 {
			function = inputSplit[idx]
		} else {
			vars = append(vars, inputSplit[idx])
		}
	}

	switch function {
	case "connectDB":
		if dbConnected {
			return "connected"
		}
		err := ReadConfig()
		if err != nil {
			return "config error : " + err.Error()
		}
		db, err = ConnectDatabase()
		if err != nil {
			return "non " + err.Error()
		}
		/*
			cdb, err = ConnectClickHouse()
			if err != nil {
				return "non " + err.Error()
			}
		*/
		return "connected"
	case "perf":
		return insertCpu(vars[0], vars[1], vars[2], vars[3], vars[4])
	default:
		return fmt.Sprintf("command %s is undefined", function)
	}

	return ""
}

func cleanInput(argv **C.char, argc int) []string {
	fmt.Printf("cleanInput\n")
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

func main() {}
