package main

import (
	"time"
)

type Indicator struct {
	ID          uint  	  `gorm:"primaryKey"`
	Name        string	  `gorm:"size:50"`
	Importance  int
	UnitMeasure string    `gorm:"size:50"`
	Dynamics    []Dynamic `gorm:"foreignKey:IndicatorID"`
}

type Enterprise struct {
	ID             uint  	 `gorm:"primaryKey"`
	Name           string	 `gorm:"size:50"`
	BankRequisites string    `gorm:"size:30"`
	PhoneNumber    int64
	ContactPerson  string    `gorm:"size:50"`
	Dynamics       []Dynamic `gorm:"foreignKey:EnterpriseID"`
}

type Dynamic struct {
	ID           uint 		`gorm:"primaryKey"`
	IndicatorID  uint
	EnterpriseID uint
	Date         time.Time  `gorm:"type:datetime"`
	Value        int64
}