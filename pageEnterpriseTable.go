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

func PageTableEnterprisesCanvas(db *gorm.DB) fyne.CanvasObject {
	data := []Enterprise{}
	db.Find(&data)
	table := widget.NewTable(
		func() (int, int) { return len(data) + 1, 5 },
		func() fyne.CanvasObject { return widget.NewLabel("wide content") },
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Row == 0 {
				switch i.Col {
				case 0:
					o.(*widget.Label).SetText("ID")
				case 1:
					o.(*widget.Label).SetText("Название")
				case 2:
					o.(*widget.Label).SetText("Реквезиты")
				case 3:
					o.(*widget.Label).SetText("Телефон")
				case 4:
					o.(*widget.Label).SetText("Контактное лицо")
				}
			} else {
				switch i.Col {
				case 0:
					o.(*widget.Label).SetText(strconv.Itoa(int(data[i.Row - 1].ID)))
				case 1:
					o.(*widget.Label).SetText(data[i.Row - 1].Name)
				case 2:
					o.(*widget.Label).SetText(data[i.Row - 1].BankRequisites)
				case 3:
					o.(*widget.Label).SetText(strconv.FormatInt(data[i.Row - 1].PhoneNumber, 10))
				case 4:
					o.(*widget.Label).SetText(data[i.Row - 1].ContactPerson)
				}
			}
		},
	)
	content := container.NewVSplit(
		container.NewVBox(
			container.NewCenter(
				container.NewHBox(
					widget.NewButton("Стартовое меню", func() { NextPage <- PageStartMenu }),	
					widget.NewLabel("Предприятия"),
			)),
			container.NewHBox(
				widget.NewButton("Добавить предприятие", func() {
					w := myApp.NewWindow("Добавить предприятие")
					nameEntry := widget.NewEntry()
					requisitesEntry := widget.NewEntry()
					phoneEntry := widget.NewEntry()
					contactEntry := widget.NewEntry()
					form := widget.NewForm(
						widget.NewFormItem("Имя", nameEntry),
						widget.NewFormItem("Реквезиты", requisitesEntry),
						widget.NewFormItem("Номер телефона", phoneEntry),
						widget.NewFormItem("Контактное лицо", contactEntry),
					)
					form.OnCancel = func() {
						w.Close()
					}
					form.OnSubmit = func() {
						name := nameEntry.Text
						phone, err := strconv.ParseInt(phoneEntry.Text, 10, 64)
						if err != nil {
							w.Close()
							ErrorWindow(err)
							return
						}
						req := requisitesEntry.Text
						contact := contactEntry.Text
						enterprise := Enterprise{0, name, req, phone, contact, nil}
						db.Create(&enterprise)
						NextPage <- PageTableEnterprises
					}
					w.SetContent(form)
					w.Resize(fyne.NewSize(300, 240))
					w.Show()
				}),
				
				widget.NewButton(" Изменить предприятие ", func() {
					w := myApp.NewWindow("Изменить предприятие")
					idEntry := widget.NewEntry()
					nameEntry := widget.NewEntry()
					requisitesEntry := widget.NewEntry()
					phoneEntry := widget.NewEntry()
					contactEntry := widget.NewEntry()
					form := widget.NewForm(
						widget.NewFormItem("ID", idEntry),
						widget.NewFormItem("Имя", nameEntry),
						widget.NewFormItem("Реквезиты", requisitesEntry),
						widget.NewFormItem("Номер телефона", phoneEntry),
						widget.NewFormItem("Контактное лицо", contactEntry),
					)
					form.OnSubmit = func() {
						name := nameEntry.Text
						phone, err := strconv.ParseInt(phoneEntry.Text, 10, 64)
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
						req := requisitesEntry.Text
						contact := contactEntry.Text
						enterprise := Enterprise{}
						result := db.First(&enterprise, id)
						if result.Error != nil {
							w.Close()
							ErrorWindow(result.Error)
							return
						}
						enterprise.Name = name
						enterprise.BankRequisites = req
						enterprise.PhoneNumber = phone
						enterprise.ContactPerson = contact
						db.Save(&enterprise)
						NextPage <- PageTableEnterprises
					}
					
					form.OnCancel = func() {
						w.Close()
					}

					content := container.NewVBox(
						widget.NewLabel("Обьект с этим ID будет изменен"),
						form,
					)

					w.SetContent(content)
					w.Resize(fyne.NewSize(300, 200))
					w.Show()
				}),

				widget.NewButton(" Удалить предприятие", func() {
					w := myApp.NewWindow("Удалить предприятие")
					idEntry := widget.NewEntry()
					content := container.NewVBox(
						widget.NewLabel("Укажите ID обьекта для удаления"),
						idEntry,
						widget.NewButton("Удалить", func() {
							id, err := strconv.Atoi(idEntry.Text)
							if err != nil {
								ErrorWindow(err)
							} else {
								result := db.Delete(&Enterprise{}, id)
								if result.Error != nil {
									ErrorWindow(result.Error)
									return
								}
								NextPage <- PageTableEnterprises
							}
						}),
					)
					w.SetContent(content)
					w.Show()
				}),
				widget.NewButton("Применить фильтр", func() {
					w := myApp.NewWindow("Применить фильтр")
					nameEntry := widget.NewEntry()
					requisitesEntry := widget.NewEntry()
					phoneEntry := widget.NewEntry()
					contactEntry := widget.NewEntry()
					form := widget.NewForm(
						widget.NewFormItem("Имя", nameEntry),
						widget.NewFormItem("Реквезиты", requisitesEntry),
						widget.NewFormItem("Номер телефона", phoneEntry),
						widget.NewFormItem("Контактное лицо", contactEntry),
					)
					form.OnCancel = func() {
						w.Close()
					}
					form.OnSubmit = func() {
						db.Find(&data)
						name := nameEntry.Text
						phone, _ := strconv.ParseInt(phoneEntry.Text, 10, 64)
						req := requisitesEntry.Text
						contact := contactEntry.Text
						dataFiltrated := make([]Enterprise, 0)
						for _, obj := range data {
							if (obj.Name == name || name == "") &&
								(obj.BankRequisites == req || req == "") &&
								(obj.PhoneNumber == phone || phone == 0) &&
								(obj.ContactPerson == contact || contact == "") { 
									dataFiltrated = append(dataFiltrated, obj) 
							}
						}
						data = dataFiltrated
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