package inquiry_service

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

// Helper functions for data formatting

// formatBalance converts pgtype.Numeric to formatted string
func formatBalance(balance pgtype.Numeric) string {
	if !balance.Valid {
		return ""
	}
	balanceFloat, _ := balance.Float64Value()
	return fmt.Sprintf("%.2f", balanceFloat.Float64)
}

// formatDate converts pgtype.Date to ISO date string (YYYY-MM-DD)
func formatDate(date pgtype.Date) string {
	if !date.Valid {
		return ""
	}
	return date.Time.Format("2006-01-02")
}

// formatTimestamp converts pgtype.Timestamp to ISO datetime string
func formatTimestamp(timestamp pgtype.Timestamp) string {
	if !timestamp.Valid {
		return ""
	}
	return timestamp.Time.Format("2006-01-02T15:04:05Z07:00")
}

// formatText converts pgtype.Text to string
func formatText(text pgtype.Text) string {
	return text.String
}
