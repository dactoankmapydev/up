package main

import (
	"backup/helper"
	"backup/model"
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

const maxS = 15*1024*1024

func initUpload(bucket, key string) (string, string, string) {
	r, w := io.Pipe()
	go func() {
		// close the writer
		defer w.Close()

		// write json data to the PipeReader through the PipeWriter
		if err := json.NewEncoder(w).Encode(&model.ReqUploadS3{
			Bucket: bucket,
			Key: key,
		}); err != nil {
			log.Fatal(err)
		}
	}()

	resp, err := http.Post("http://localhost:3000/upload/multipart/init", "application/json", r)
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(resp.Body)
	request := model.ReqUploadS3{}
	errRequest := decoder.Decode(&request)
	if errRequest != nil {
		log.Println(errRequest)
	}

	//log.Println(request.Key, request.Bucket, request.UploadID)
	return request.Bucket, request.Key, request.UploadID
}

func openFile(filePath string) *os.File {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Error opening file!!!")
	}
	return file
}

func uploadFilePartRoutine(bucket, key, filePath string) {
	// init upload
	bucket, key, uploadID := initUpload(bucket, key)

	// use errGroup
	sem := semaphore.NewWeighted(int64(1000))
	group, ctx := errgroup.WithContext(context.Background())

	file := openFile(filePath)
	stats, _ := file.Stat()
	fileSize := stats.Size()
	defer file.Close()

	nBytes, nChunks := int64(0), int64(0)
	r := bufio.NewReader(file)
	buf := make([]byte, maxS)

	remaining := int(fileSize)
	var partNum = 0
	var currentSize int
	var completedParts []*s3.CompletedPart

	// use waitGroup
	//var wg sync.WaitGroup

	for {
		if remaining < maxS {
			currentSize = remaining
		} else {
			currentSize = maxS
		}
		bodyBuf := &bytes.Buffer{}
		bodyWriter := multipart.NewWriter(bodyBuf)
		fileWriter, err := bodyWriter.CreateFormFile("file", file.Name())
		if err != nil {
			log.Println("error writing to buffer", err)
		}

		// read content to buffer
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		nChunks++
		nBytes += int64(len(buf))
		remaining -= currentSize
		partNum++
		//log.Printf("Part %v complete, read %d bytes\n", partNum, len(buf))

		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		// io copy
		_, err = io.Copy(fileWriter, bytes.NewReader(buf))
		if err != nil {
			log.Println(err)
		}
		contentType := bodyWriter.FormDataContentType()
		_ = bodyWriter.Close()
		uri := fmt.Sprintf("http://localhost:3000/upload/multipart/upload_part?bucket=%s&key=%s&upload_id=%s&part_number=%d",bucket, key, uploadID, partNum)

		// use errGroup
		errAcquire := sem.Acquire(ctx, 1)
		if errAcquire != nil {
			log.Printf("Acquire err = %+v\n", err)
			continue
		}

		// use waitGroup
		//wg.Add(1)

		buffTemp := buf

		// use waitGroup go func() error {}()

		// use errGroup
		group.Go(func() error {
			// use waitGroup
			//defer wg.Done()

			// use errGroup
			defer sem.Release(1)

			resp, err := helper.HttpClient.PostRequestWithRetries(uri, buffTemp, contentType)
			if err != nil {
				log.Println(err)
			}
			//log.Printf("Part %v complete, read %d bytes\n", partNum, len(buffTemp))

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			data := string(body)
			var mapData model.DataUploadResp
			if err := json.Unmarshal([]byte(data), &mapData); err != nil {
				log.Println(err)
			}
			completedParts = append(completedParts, mapData.Parts)
			return nil
		})
	}
	// use waitGroup
	//wg.Wait()

	// use errGroup
	if err := group.Wait(); err != nil {
		log.Printf("g.Wait() err = %+v\n", err)
	}
	log.Println("Bytes:", nBytes, "Chunks:", nChunks)

	// complete upload
	completeUpload(bucket, key, uploadID, completedParts)
}


func completeUpload(bucket, key, uploadID string, completedParts []*s3.CompletedPart) {
	r, w := io.Pipe()
	go func() {
		defer w.Close()
		if err := json.NewEncoder(w).Encode(&model.ResponseDataS3{
			Bucket: bucket,
			Key: key,
			UploadID: uploadID,
			Parts: completedParts,
		}); err != nil {
			log.Fatal(err)
		}
	}()
	resp, err := http.Post("http://localhost:3000/upload/multipart/complete", "application/json", r)
	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(resp.Body)
	request := model.ResponseDataS3{}
	errRequest := decoder.Decode(&request)
	if errRequest != nil {
		log.Println(errRequest)
	}
	log.Println(request.UploadID, request.Key, request.Bucket, request.Etag)
}

func main() {
	filePath := "/home/dactoan/Downloads/Win10_20H2_v2_English_x64.iso"
	uploadFilePartRoutine("toannd-test-2", "win10_iso", filePath)

	/*filePath := "/home/dactoan/Downloads/go1.16.linux-amd64.tar.gz"
	uploadFilePartRoutine("toannd-test-2", "go1.16", filePath)*/
}
