package main

import (
	"fmt"
	"http-speed-test/bytes"
	"io"
	"log"
	"net/http"
	"time"
)

type countingWriter struct {
	writeBytes uint64
}

func (r *countingWriter) Write(p []byte) (int, error) {
	r.writeBytes += uint64(len(p))
	return len(p), nil
}

func RunClient(serverAddress string) {
	// size := 10 * 1024 * 1024 * 1024
	start := time.Now()

	resp, err := http.Get(serverAddress)
	if err != nil {
		log.Fatal(err)
	}

	writer := &countingWriter{}
	_, err = io.Copy(writer, resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	duration := time.Now().Sub(start)
	log.Printf("downloaded %s for %s, speed %s", bytes.IBytes(writer.writeBytes), duration, FormatSpeed(writer.writeBytes, duration))
}

func FormatSpeed(bytes uint64, duration time.Duration) string {
	speed := float64(bytes) / duration.Seconds()
	if speed < 1024 {
		return fmt.Sprintf("%f b/s", speed)
	}
	if speed < 1024*1024 {
		return fmt.Sprintf("%f Kib/s", speed/1024)
	}
	if speed < 1024*1024*1024 {
		return fmt.Sprintf("%f Mib/s", speed/1024/1024)
	}
	return fmt.Sprintf("%f Gib/s", speed/1024/1024/1024)
}

func main() {
	// downloadUrl := "https://corbit-dev-868303926763-us-east-1.s3.amazonaws.com/cms/prod/world/-LssKL9rSbi5lM4P08qh/topic/-M0T9huUkb5E7ypz0xBN/9762068c4c00cfbd82f28cc323649ace.json"
	downloadUrl := "https://d2niex7nhy7zda.cloudfront.net/cms/prod/publish/117/world/-LssKL9rSbi5lM4P08qh/world_sum.zip"
	count := 100
	for i := 0; i < count; i++ {
		RunClient(downloadUrl)
	}
}
