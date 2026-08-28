package store

import "fmt"

func AuditKey(recordID string, sequence int) string {
	return fmt.Sprintf("%s:%06d", recordID, sequence)
}

func WorkflowKey(recordID string, sequence int) string {
	return fmt.Sprintf("%s:%06d", recordID, sequence)
}

func AttachmentKey(recordID, name string) string {
	return recordID + ":" + name
}
