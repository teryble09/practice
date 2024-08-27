package main

import (
	//"time"

	"fyne.io/fyne/v2"
	//"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	//"gorm.io/gorm"
)

const (
	PageStartMenu = iota
	PageTableIndicators
	PageTableEnterprises
	PageTableDynamics
)

func ErrorWindow(err error) {
	w := myApp.NewWindow("Ошибка")
	w.SetContent(widget.NewLabel(err.Error()))
	w.Show()
}

func PageMenuCanvas() fyne.CanvasObject {
	return container.NewVBox(
		container.NewCenter(widget.NewLabel("Выберите нужную опцию")),
		widget.NewLabel(""),
		widget.NewButton("1) Показатели", func() { NextPage <- PageTableIndicators }),
		widget.NewButton("2) Предприятия", func() { NextPage <- PageTableEnterprises }),
		widget.NewButton("3) Динамика", func() { NextPage <- PageTableDynamics }),
	)
}
