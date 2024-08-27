package main

import (
	//"strconv"
	//"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	//"fyne.io/fyne/v2/layout"
	//"fyne.io/fyne/v2/widget"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var NextPage = make(chan int)
var myApp = app.New()

func main() {
	db, err := gorm.Open(sqlite.Open("practice.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	defer func() {
		db, err := db.DB()
		if err != nil {
			panic("failed to get database from gorm DB: " + err.Error())
		}
		err = db.Close()
		if err != nil {
			panic("failed to close database: " + err.Error())
		}
	}()

	err = db.AutoMigrate(&Indicator{}, &Dynamic{}, &Enterprise{})

	if err != nil {
		panic("failed migration")
	}

	myWindow := myApp.NewWindow("Practice")
	myWindow.Resize(fyne.NewSize(600, 400))
	content := PageMenuCanvas()
	myWindow.SetContent(content)
	// Main Window Pages manager
	go func() {
		for {
			switch <-NextPage {
			case PageStartMenu:
				content := PageMenuCanvas()
				myWindow.SetContent(content)
			case PageTableIndicators:
				content := PageTableIndicatorsCanvas(db)
				myWindow.SetContent(content)
			case PageTableEnterprises:
				content := PageTableEnterprisesCanvas(db)
				myWindow.SetContent(content)
			case PageTableDynamics:
				content := PageTableDynamicsCanvas(db)
				myWindow.SetContent(content)
			}
		}
	}()
	myWindow.ShowAndRun()
}
