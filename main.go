package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go [height] [iterations]")
		os.Exit(1)
	}
	height, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Invalid height:", os.Args[1])
		os.Exit(1)
	}
	iterations, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Invalid iterations:", os.Args[2])
		os.Exit(1)
	}

	var totalDiff int64

	for i := 0; i < iterations; i++ {
		blockTime, err := getBlockTime(height)
		if err != nil {
			fmt.Println("Failed to get block time:", err)
			os.Exit(1)
		}

		prevBlockTime, err := getBlockTime(height - 1)
		if err != nil {
			fmt.Println("Failed to get previous block time:", err)
			os.Exit(1)
		}

		diff := blockTime - prevBlockTime

		totalDiff += diff

		height--
	}

	// 밀리초 단위로 변환하여 평균 값을 계산합니다.
	average := float64(totalDiff) / float64(iterations) / 1000.0

	fmt.Printf("Average time difference: %.4f seconds\n", average)
}

func getBlockTime(height int) (int64, error) {
	url := fmt.Sprintf("http://localhost:26657/block?height=%d", height)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}
	blockHeader := result["result"].(map[string]interface{})["block"].(map[string]interface{})["header"].(map[string]interface{})
	blockTimeStr := blockHeader["time"].(string)

	blockTime, err := time.Parse(time.RFC3339Nano, blockTimeStr)
	if err != nil {
		return 0, err
	}

	return blockTime.UnixNano() / int64(time.Millisecond), nil
}
