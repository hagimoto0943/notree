package main

import (
	"context"
	"fmt"
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("❌ エラー: %v", err)
	}
}

func run() error {
	ctx := context.Background()

	// 1. 設定
	cfg, err := NewConfigFromEnv()
	if err != nil {
		// ローカルモードの分岐などはここでやる
		return err
	}
	fmt.Printf("🚀 処理開始: %s\n", cfg.ObjectKey)

	// 2. S3準備
	s3Client, err := NewS3Client(ctx, cfg.Region)
	if err != nil {
		return err
	}

	// 3. ダウンロード
	file, err := s3Client.Download(ctx, cfg.BucketName, cfg.ObjectKey)
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	// 4. 解析
	text, err := ExtractText(file)
	if err != nil {
		return err
	}

	// 5. 結果
	fmt.Println("----- 結果 -----")
	fmt.Println(text)

	return nil
}
