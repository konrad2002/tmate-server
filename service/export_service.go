package service

import (
	"bytes"
	"fmt"
	"github.com/konrad2002/tmate-server/dto"
	"github.com/konrad2002/tmate-server/model"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ExportService struct {
	memberService MemberService
}

func NewExportService(ms MemberService) ExportService {
	return ExportService{
		memberService: ms,
	}
}

func (es *ExportService) ExportFromQueryId(queryId primitive.ObjectID, sortField string, sortDirection int) (*bytes.Buffer, error) {
	members, fields, query, err := es.memberService.GetAllByQueryId(queryId, sortField, sortDirection)
	if err != nil {
		return nil, err
	}

	return es.membersAndFieldsToExcel(dto.QueryResultDto{
		Members: *members,
		Fields:  *fields,
		Query:   *query,
	})
}

func (es *ExportService) membersAndFieldsToExcel(result dto.QueryResultDto) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	index, err := f.NewSheet("Abfrage")
	_, err = f.NewSheet("Mitglieder")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	f.DeleteSheet("Sheet1")

	// query details
	f.SetCellValue("Abfrage", "A1", "tMate Abfrage")
	f.SetCellValue("Abfrage", "A2", "Abfrage:")
	f.SetCellValue("Abfrage", "A3", "Ergebnisse")
	f.SetCellValue("Abfrage", "A4", "Zeitpunkt")
	f.SetCellValue("Abfrage", "B2", result.Query.Name)
	f.SetCellValue("Abfrage", "B3", len(result.Members))
	f.SetCellValue("Abfrage", "B4", time.Now().Format("02.01.2006 15:04:05"))

	for i, field := range result.Fields {
		f.SetCellValue("Mitglieder", fmt.Sprintf("%s%d", numberToColumn(i), 1), field.DisplayName)
	}

	for i, member := range result.Members {
		for j, field := range result.Fields {
			if field.Type == model.Date {
				f.SetCellValue("Mitglieder", fmt.Sprintf("%s%d", numberToColumn(j), i+2), (member.Data[field.Name]).(primitive.DateTime).Time().Format("02.01.2006"))
				continue
			}
			if field.Type == model.Select || field.Type == model.MultiSelect {
				if member.Data[field.Name] != nil {
					f.SetCellValue("Mitglieder", fmt.Sprintf("%s%d", numberToColumn(j), i+2), field.Data.Options[member.Data[field.Name].(string)])
				}
				continue
			}
			f.SetCellValue("Mitglieder", fmt.Sprintf("%s%d", numberToColumn(j), i+2), member.Data[field.Name])
		}
	}

	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	if err := f.SaveAs(fmt.Sprintf("Mitglieder_%s.xlsx", time.Now().Format("2006-01-02_15-04-05"))); err != nil {
		fmt.Println(err)
	}

	// Save to memory buffer
	buf := new(bytes.Buffer)
	err = f.Write(buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func numberToColumn(n int) string {
	column := ""
	for n > 0 {
		n-- // Adjust for 1-based indexing
		column = string(rune(n%26+'A')) + column
		n /= 26
	}
	return column
}
