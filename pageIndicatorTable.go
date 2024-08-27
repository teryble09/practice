package main

import (
	"strconv"
	//"time"

	"fyne.io/fyne/v2"
	//"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"gorm.io/gorm"
)

func PageTableIndicatorsCanvas(db *gorm.DB) fyne.CanvasObject {
	data := []Indicator{}
	db.Find(&data)
	table := widget.NewTable(
		func() (int, int) { return len(data) + 1, 4 },
		func() fyne.CanvasObject { return widget.NewLabel("wide content") },
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Row == 0 {
				switch i.Col {
				case 0:
					o.(*widget.Label).SetText("ID")
				case 1:
					o.(*widget.Label).SetText("Название")
				case 2:
					o.(*widget.Label).SetText("Важность")
				case 3:
					o.(*widget.Label).SetText("Единица измерения")
				}
			} else {
				switch i.Col {
				case 0:
					o.(*widget.Label).SetText(strconv.Itoa(int(data[i.Row - 1].ID)))
				case 1:
					o.(*widget.Label).SetText(data[i.Row - 1].Name)
				case 2:
					o.(*widget.Label).SetText(strconv.Itoa(data[i.Row - 1].Importance))
				case 3:
					o.(*widget.Label).SetText(data[i.Row - 1].UnitMeasure)
				}
			}
		},
	)
	content := container.NewVSplit(
		container.NewVBox(
			container.NewCenter(
				container.NewHBox(
					widget.NewButton("Стартовое меню", func() { NextPage <- PageStartMenu }),	
					widget.NewLabel("Показатели"),
			)),
			container.NewHBox(
				widget.NewButton("Добавить показатель", func() {
					w := myApp.NewWindow("Добавить показатель")
					nameEntry := widget.NewEntry()
					importanceEntry := widget.NewEntry()
					unitEntry := widget.NewEntry()
					form := widget.NewForm(
						widget.NewFormItem("Имя", nameEntry),
						widget.NewFormItem("Важность", importanceEntry),
						widget.NewFormItem("Единица измерения", unitEntry),
					)
					form.OnCancel = func() {
						w.Close()
					}
					form.OnSubmit = func() {
						name := nameEntry.Text
						importance, err := strconv.Atoi(importanceEntry.Text)
						if err != nil {
							w.Close()
							ErrorWindow(err)
							return
						}
						unit := unitEntry.Text
						ind := Indicator{0, name, importance, unit, nil}
						db.Create(&ind)
						NextPage <- PageTableIndicators
					}
					w.SetContent(form)
					w.Resize(fyne.NewSize(400, 170))
					w.Show()
				}),

				widget.NewButton(" Изменить показатель ", func() {
					w := myApp.NewWindow("Изменить показатель")
					idEntry := widget.NewEntry()
					nameEntry := widget.NewEntry()
					importanceEntry := widget.NewEntry()
					unitEntry := widget.NewEntry()
					form := widget.NewForm(
						widget.NewFormItem("ID", idEntry),
						widget.NewFormItem("Имя", nameEntry),
						widget.NewFormItem("Важность", importanceEntry),
						widget.NewFormItem("Единица измерения", unitEntry),
					)
					form.OnCancel = func() {
						w.Close()
					}
					form.OnSubmit = func() {
						name := nameEntry.Text
						id_int, err := strconv.Atoi(idEntry.Text)
						id := uint(id_int)
						if err != nil {
							w.Close()
							ErrorWindow(err)
							return
						}
						importance, err := strconv.Atoi(importanceEntry.Text)
						if err != nil {
							w.Close()
							ErrorWindow(err)
							return
						}
						unit := unitEntry.Text
						ind := Indicator{}
						result := db.First(&ind, id)
						if result.Error != nil {
							w.Close()
							ErrorWindow(result.Error)
							return
						}	
						ind.Name = name
						ind.Importance = importance
						ind.UnitMeasure = unit
						db.Save(&ind)
						NextPage <- PageTableIndicators
					}
					content := container.NewVBox(
						widget.NewLabel("Обьект с этим ID будет изменен"),
						form,
					)
					w.SetContent(content)
					w.Resize(fyne.NewSize(300, 200))
					w.Show()
				}),

				widget.NewButton(" Удалить показатель ", func() {
					w := myApp.NewWindow("Удалить показатель")
					idEntry := widget.NewEntry()
					content := container.NewVBox(
						widget.NewLabel("Укажите ID обьекта для удаления"),
						idEntry,
						widget.NewButton("Удалить", func() {
							id, err := strconv.Atoi(idEntry.Text)
							if err != nil {
								ErrorWindow(err)
							} else {
								db.Delete(&Indicator{}, id)
								NextPage <- PageTableIndicators
							}
						}),
					)
					w.SetContent(content)
					w.Show()
				}),

				widget.NewButton("Применить фильтр", func() {
					w := myApp.NewWindow("Применить фильтр")
					nameEntry := widget.NewEntry()
					importanceEntry := widget.NewEntry()
					unitEntry := widget.NewEntry()
					form := widget.NewForm(
						widget.NewFormItem("Имя", nameEntry),
						widget.NewFormItem("Важность", importanceEntry),
						widget.NewFormItem("Единица измерения", unitEntry),
					)
					form.OnCancel = func() {
						w.Close()
					}
					form.OnSubmit = func() {
						db.Find(&data)
						name := nameEntry.Text
						importance, err := strconv.Atoi(importanceEntry.Text)
						if importanceEntry.Text != "" && err != nil {
							w.Close()
							ErrorWindow(err)
						}
						unit := unitEntry.Text
						dataFiltrated := make([]Indicator, 0)
						for _, obj := range data {
							if (obj.Name != name && name != "") ||
								(obj.Importance != importance && importance != 0) ||
								(obj.UnitMeasure != unit && unit != "") {
								continue
							}
							dataFiltrated = append(dataFiltrated, obj)
						}
						data = dataFiltrated
						table.Refresh()
					}
					w.SetContent(form)
					w.Resize(fyne.NewSize(300, 170))
					w.Show()
				}),
			),
		),
		table,
	)
	content.SetOffset(0)
	return content
}