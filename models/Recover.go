package models

import (
	"fmt"

	"github.com/tealeg/xlsx"
)

func RegCase() {
	excelFileName := "./data.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		panic("invalid file name")
	}
	for _, sheet := range xlFile.Sheets {
		for i, row := range sheet.Rows {
			if i > 200 {
				break
			}
			if len(row.Cells) == 11 {
				c := &Contract{
					Seq:            row.Cells[0].String(),
					Contract_id:    row.Cells[1].String(),
					Client_name:    row.Cells[2].String(),
					Country:        row.Cells[3].String(),
					Project_type:   row.Cells[4].String(),
					Consulter:      row.Cells[5].String(),
					Consulter_name: row.Cells[6].String(),
					Secretary:      row.Cells[7].String(),
					Secretary_name: row.Cells[8].String(),
					Zhuan_an_date:  row.Cells[9].String(),
					Current_state:  row.Cells[10].String(),
				}
				AddContract(c)

			} else {
				for _, cell := range row.Cells {
					fmt.Printf("%s\n", cell.String())
				}
			}
		}
	}
}
