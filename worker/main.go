package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"github.com/ledongthuc/pdf"
)

func main() {
	_ = godotenv.Load()

	bucketName := os.Getenv("BUCKET_NAME")
	objectKey := os.Getenv("FILE_KEY")

	var textContent string
	var err error

	// --- 1. ファイル取得 ---
	if bucketName == "" || objectKey == "" {
		fmt.Println("⚠️ S3設定なし: ローカルの test.pdf を使用")
		textContent, err = extractText("test.pdf")
	} else {
		fmt.Printf("🚀 S3モード: s3://%s/%s\n", bucketName, objectKey)
		cfg, _ := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-1"))
		s3Client := s3.NewFromConfig(cfg)

		file, err := downloadFromS3(s3Client, bucketName, objectKey)
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(file.Name()) // 終わったら消す

		textContent, err = extractText(file.Name())
	}

	if err != nil {
		log.Fatalf("PDF解析失敗: %v", err)
	}

	// --- 2. 結果出力 (整形なしでそのまま出す) ---
	fmt.Println("----- 解析結果 (Raw) -----")
	// 簡易的な改行削除だけやっておく（AIに渡すときは改行がない方が扱いやすい場合が多いので）
	fmt.Println(cleanText(textContent))
	fmt.Println("--------------------------")
}

// 以下、変更なし（コピペ用）
func downloadFromS3(client *s3.Client, bucket, key string) (*os.File, error) {
	resp, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucket, Key: &key,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	tmpFile, err := os.CreateTemp("", "worker-*.pdf")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return nil, err
	}
	return tmpFile, nil
}

func extractText(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return "", err
	}
	r, err := pdf.NewReader(f, st.Size())
	if err != nil {
		return "", err
	}
	var content string
	for i := 1; i <= r.NumPage(); i++ {
		p := r.Page(i)
		text, err := p.GetPlainText(nil)
		if err == nil {
			content += text + "\n"
		}
	}
	return content, nil
}

func cleanText(text string) string {
	// 連続する改行だけ残して、単発の改行は繋げてしまう（文章をつなげるため）
	text = strings.ReplaceAll(text, "\n\n", "PLACEHOLDER")
	text = strings.ReplaceAll(text, "\n", "")
	text = strings.ReplaceAll(text, "PLACEHOLDER", "\n\n")
	return text
}
