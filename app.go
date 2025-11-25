package main

import (
	"context"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Content string `json:"content"`
	IsDone  bool   `json:"is_done"`
}

// App struct
type App struct {
	ctx context.Context
	db  *gorm.DB
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	db, err := gorm.Open(sqlite.Open("notree.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Task{})

	a.db = db

	var count int64
	a.db.Model(&Task{}).Count(&count)
	if count == 0 {
		a.db.Create(&Task{Content: "notreeの開発環境構築", IsDone: true})
		a.db.Create(&Task{Content: "GoでDB接続テスト", IsDone: false})
		a.db.Create(&Task{Content: "試験の勉強", IsDone: false})
	}
}

func (a *App) GetTasks() []Task {
	var tasks []Task

	result := a.db.Find(&tasks)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return tasks
}

// Greet returns a greeting for the given name
func (a *App) AddTask(content string) []Task {
	// INSERT INTO tasks ...
	a.db.Create(&Task{Content: content, IsDone: false})
	return a.GetTasks() // 最新リストを返す
}