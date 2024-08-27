package main

import (
	"fmt"
	"strconv"
	"time"

	//"time"

	"fyne.io/fyne/v2"
	//"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"gorm.io/gorm"
)

func PageTableDynamicsCanvas(db *gorm.DB) fyne.CanvasObject {

	type DynamicWithRelation struct {
		Dynamic
		IndicatorName string
		EnterpriseName string
	}

	dataDynamics := []DynamicWithRelation{}
	err := db.Table("dynamics").
		Select("dynamics.*, indicators.name as indicator_name, enterprises.name as enterprise_name").
		Joins("JOIN indicators ON indicators.id = dynamics.indicator_id").
		Joins("JOIN enterprises ON enterprises.id = dynamics.enterprise_id").
		Find(&dataDynamics).Error

	if err != nil {
		ErrorWindow(fmt.Errorf("Не получилось загрузить данные из таблицы" + err.Error()))
		NextPage <- PageStartMenu
		return nil
	}

	dataIndicators := []Indicator{}
	indicatorMap := make(map[string]uint, 100)
	indicatorNames := make([]string, 0)
	db.Find(&dataIndicators)
	
	for _, obj := range dataIndicators {
		indicatorMap[obj.Name] = obj.ID
		indicatorNames = append(indicatorNames, obj.Name)
	}
	
	dataEnterprises := []Enterprise{}
	enterpriseMap := make(map[string]uint, 100)
	enterpriseNames := make([]string, 0)
	db.Find(&dataEnterprises)

	for _, obj := range dataEnterprises {
		enterpriseMap[obj.Name] = obj.ID
		enterpriseNames = append(enterpriseNames, obj.Name)
	}

	table := widget.NewTable(
		func() (int, int) { return len(dataDynamics) + 1, 5 },
		func() fyne.CanvasObject { return widget.NewLabel("wide content") },
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Row == 0 {
				switch i.Col {
				case 0:
					o.(*widget.Label).SetText("ID")
				case 1:
					o.(*widget.Label).SetText("Предприятие")
				case 2:
					o.(*widget.Label).SetText("Индикатор")
				case 3:
					o.(*widget.Label).SetText("Значение")
				case 4:
					o.(*widget.Label).SetText("Дата")
				}
			} else {
				switch i.Col {
				case 0:
					o.(*widget.Label).SetText(strconv.Itoa(int(dataDynamics[i.Row - 1].ID)))
				case 1:
					o.(*widget.Label).SetText(dataDynamics[i.Row - 1].EnterpriseName)
				case 2:
					o.(*widget.Label).SetText(dataDynamics[i.Row - 1].IndicatorName)
				case 3:
					o.(*widget.Label).SetText(strconv.FormatInt(dataDynamics[i.Row - 1].Value, 10))
				case 4:
					o.(*widget.Label).SetText(dataDynamics[i.Row - 1].Date.Format("02 01 2006"))
				}
			}
		},
	)
	content := container.NewVSplit(
		container.NewVBox(
			container.NewCenter(
				container.NewHBox(
					widget.NewButton("Стартовое меню", func() { NextPage <- PageStartMenu }),	
					widget.NewLabel("Динамика"),
			)),
			container.NewHBox(
				widget.NewButton("Добавить динамику", func() {
					w := myApp.NewWindow("Добавить динамику")
					enterpriseSelect := widget.NewSelect(enterpriseNames, func(s string) {})
					indicatorSelect := widget.NewSelect(indicatorNames, func(s string) {})
					valueEntry := widget.NewEntry()
					dateEntry := widget.NewEntry()
					form := widget.NewForm(
						widget.NewFormItem("Предприятие", enterpriseSelect),
						widget.NewFormItem("Индикатор", indicatorSelect),
						widget.NewFormItem("Значение", valueEntry),
						widget.NewFormItem("Дата (Формат:23 04 2006)", dateEntry),
					)
					form.OnCancel = func() {
						w.Close()
					}
					form.OnSubmit = func() {
						indicatorID := indicatorMap[indicatorSelect.Selected]
						enterpriseID := enterpriseMap[enterpriseSelect.Selected]
						value, err := strconv.ParseInt(valueEntry.Text, 10, 64)
						if err != nil {
							w.Close()
							ErrorWindow(err)
							return
						}
						date, err := time.Parse("02 01 2006", dateEntry.Text)
						if err != nil {
							w.Close()
							ErrorWindow(err)
							return
						}
						dynamic := Dynamic{0, indicatorID, enterpriseID, date, value}
						db.Create(&dynamic)
						NextPage <- PageTableDynamics
					}
					w.SetContent(form)
					w.Resize(fyne.NewSize(300, 240))
					w.Show()
				}),
				
				widget.NewButton(" Изменить динамику ", func() {
					w := myApp.NewWindow("Изменить динамику")
					enterpriseSelect := widget.NewSelect(enterpriseNames, func(s string) {})
					indicatorSelect := widget.NewSelect(indicatorNames, func(s string) {})
					valueEntry := widget.NewEntry()
					dateEntry := widget.NewEntry()
					idEntry := widget.NewEntry()
					form := widget.NewForm(
						widget.NewFormItem("ID", idEntry),
						widget.NewFormItem("Предприятие", enterpriseSelect),
						widget.NewFormItem("Индикатор", indicatorSelect),
						widget.NewFormItem("Значение", valueEntry),
						widget.NewFormItem("Дата (Формат:23 04 2006)", dateEntry),
					)
					form.OnCancel = func() {
						w.Close()
					}
					form.OnSubmit = func() {
						indicatorID := indicatorMap[indicatorSelect.Selected]
						enterpriseID := enterpriseMap[enterpriseSelect.Selected]
						value, err := strconv.ParseInt(valueEntry.Text, 10, 64)
						if err != nil {
							w.Close()
							ErrorWindow(err)
							return
						}
						id_int, err := strconv.Atoi(idEntry.Text)
						id := uint(id_int)
						if err != nil {
							w.Close()
							ErrorWindow(err)
							return
						}
						date, err := time.Parse("02 01 2006", dateEntry.Text)
						if err != nil {
							w.Close()
							ErrorWindow(err)
							return
						}
						dynamic := Dynamic{}
						result := db.First(&dynamic, id)
						if result.Error != nil {
							w.Close()
							ErrorWindow(err)
							return
						}
						dynamic.Date = date
						dynamic.EnterpriseID = enterpriseID
						dynamic.IndicatorID = indicatorID
						dynamic.Value = value
						db.Save(&dynamic)
						NextPage <- PageTableDynamics
					}

					content := container.NewVBox(
						widget.NewLabel("Обьект с этим ID будет изменен"),
						form,
					)

					w.SetContent(content)
					w.Resize(fyne.NewSize(300, 200))
					w.Show()
				}),

				widget.NewButton(" Удалить динамику", func() {
					w := myApp.NewWindow("Удалить динамику")
					idEntry := widget.NewEntry()
					content := container.NewVBox(
						widget.NewLabel("Укажите ID обьекта для удаления"),
						idEntry,
						widget.NewButton("Удалить", func() {
							id, err := strconv.Atoi(idEntry.Text)
							if err != nil {
								ErrorWindow(err)
							} else {
								result := db.Delete(&Dynamic{}, id)
								if result.Error != nil {
									ErrorWindow(result.Error)
									return
								}
								NextPage <- PageTableDynamics
							}
						}),
					)
					w.SetContent(content)
					w.Show()
				}),
				widget.NewButton("Применить фильтр", func() {
					w := myApp.NewWindow("Применить фильтр")
					enterpriseSelect := widget.NewSelect(enterpriseNames, func(s string) {})
					indicatorSelect := widget.NewSelect(indicatorNames, func(s string) {})
					valueEntry := widget.NewEntry()
					dateEntry := widget.NewEntry()
					form := widget.NewForm(
						widget.NewFormItem("Предприятие", enterpriseSelect),
						widget.NewFormItem("Индикатор", indicatorSelect),
						widget.NewFormItem("Значение", valueEntry),
						widget.NewFormItem("Дата (Формат:23 04 2006)", dateEntry),
					)
					form.OnCancel = func() {
						w.Close()
					}
					form.OnSubmit = func() {
						db.Table("dynamics").
							Select("dynamics.*, indicators.name as indicator_name, enterprises.name as enterprise_name").
							Joins("JOIN indicators ON indicators.id = dynamics.indicator_id").
							Joins("JOIN enterprises ON enterprises.id = dynamics.enterprise_id").
							Find(&dataDynamics)
						enterprise := enterpriseSelect.Selected
						indicator := indicatorSelect.Selected
						value, _ := strconv.ParseInt(valueEntry.Text, 10, 64)
						date, _ := time.Parse("02 01 2006", dateEntry.Text)
						dataFiltrated := make([]DynamicWithRelation, 0)
						for _, obj := range dataDynamics {
							if (obj.IndicatorName == indicator || indicator == "") &&
								(obj.EnterpriseName == enterprise || enterprise == "") &&
								(obj.Value == value || value == 0) &&
								(obj.Date == date || date == time.Time{}) { 
									dataFiltrated = append(dataFiltrated, obj) 
							}
						}
						dataDynamics = dataFiltrated
						table.Refresh()
					}
					w.SetContent(form)
					w.Resize(fyne.NewSize(300, 240))
					w.Show()
				}),
			),
		),
		table,
	)
	content.SetOffset(0)
	return content
}