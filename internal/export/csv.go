package export

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"noticeword/internal/model"
)

var Header = []string{"id", "community", "title", "body", "description", "status", "character_count", "revision"}

func Encode(records []model.Record) (string, error) {
	var builder strings.Builder
	writer := csv.NewWriter(&builder)
	if err := writer.Write(Header); err != nil {
		return "", err
	}
	for _, record := range records {
		row := []string{record.ID, record.Community, record.Title, record.Body, record.Description, string(record.Status), fmt.Sprint(record.CharacterCount), fmt.Sprint(record.Revision)}
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}
	return builder.String(), nil
}

func Decode(input string) ([]model.ImportRow, error) {
	reader := csv.NewReader(strings.NewReader(input))
	reader.FieldsPerRecord = -1
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}
	if !sameHeader(header) {
		return nil, fmt.Errorf("unexpected csv header")
	}
	rows := make([]model.ImportRow, 0)
	for {
		values, readErr := reader.Read()
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return nil, readErr
		}
		if len(values) != len(Header) {
			return nil, fmt.Errorf("row has %d fields", len(values))
		}
		rows = append(rows, model.ImportRow{ID: values[0], Community: values[1], Title: values[2], Body: values[3], Description: values[4], Tags: splitTags(values[7])})
	}
	return rows, nil
}

func sameHeader(values []string) bool {
	if len(values) != len(Header) {
		return false
	}
	for i := range Header {
		if values[i] != Header[i] {
			return false
		}
	}
	return true
}

func splitTags(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, "|")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if strings.TrimSpace(part) != "" {
			result = append(result, strings.TrimSpace(part))
		}
	}
	return result
}
