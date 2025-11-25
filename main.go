package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ledongthuc/pdf"
	"github.com/joho/godotenv"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "notree",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})
	err = godotenv.Load()
	if err != nil {
		log.Println("注意: .envファイルが見つかりません (環境変数が直接設定されている場合はOK)")
	}

	if err != nil {
		println("Error:", err.Error())
	}

	// --- A. ローカルファイルでテストしたい場合 ---
		// S3を使わず、手元の "test.pdf" を読むならここを有効化して、下のBをコメントアウト
		/*
		content, err := extractText("test.pdf")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(cleanText(content))
		return
		*/

		// --- B. 本番用 (S3から取得) ---
		bucketName := os.Getenv("BUCKET_NAME")
		objectKey := os.Getenv("FILE_KEY")

		if bucketName == "" || objectKey == "" {
			log.Fatal("エラー: 環境変数 BUCKET_NAME と FILE_KEY が必要です")
		}

		// 1. AWS設定読み込み
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-1"))
		if err != nil {
			log.Fatalf("AWS設定エラー: %v", err)
		}

		// 2. S3クライアント
		s3Client := s3.NewFromConfig(cfg)

		// 3. ダウンロード
		fmt.Printf("S3からダウンロード中... s3://%s/%s\n", bucketName, objectKey)
		file, err := downloadFromS3(s3Client, bucketName, objectKey)
		if err != nil {
			log.Fatalf("ダウンロード失敗: %v", err)
		}
		defer os.Remove(file.Name()) // 終わったら一時ファイル削除

		// 4. 解析
		content, err := extractText(file.Name())
		if err != nil {
			log.Fatalf("解析失敗: %v", err)
		}

		// 5. 結果表示
		fmt.Println("----- 解析結果 -----")
		fmt.Println(cleanText(content)) // クリーニングして表示
		fmt.Println("----- 完了 -----")
}

// S3ダウンロード関数
func downloadFromS3(client *s3.Client, bucket, key string) (*os.File, error) {
	resp, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
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

// PDF解析関数 (ledongthuc/pdf 使用版)
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
		if err != nil {
			continue
		}
		content += text + "\n"
	}
	return content, nil
}

// テキスト整形関数 (ゴミ取り)
func cleanText(text string) string {
	text = strings.ReplaceAll(text, "\n\n", "PLACEHOLDER")
	text = strings.ReplaceAll(text, "\n", "")
	text = strings.ReplaceAll(text, "PLACEHOLDER", "\n\n")
	return text
}