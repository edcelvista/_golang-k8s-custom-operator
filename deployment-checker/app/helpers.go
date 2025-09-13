package main

import (
	"encoding/base64"
	"fmt"
	"os"
)

func int64ToInt32Ptr(i int64) int32 {
	i32 := int32(i)
	return i32
}

func int64ToInt32PtrP(i int64) *int32 {
	i32 := int32(i) // Convert int64 to int32
	return &i32     // Return pointer to int32
}

func readFileContent(file string, isBase64Encoded bool) (any, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}

	if isBase64Encoded {
		dataBase64 := base64.StdEncoding.EncodeToString(data)
		return string(dataBase64), nil
	}

	return string(data), nil
}
