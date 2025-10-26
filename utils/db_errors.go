package utils

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// FormatDBError converts database errors to human-readable error messages
func FormatDBError(err error) error {
	if err == nil {
		return nil
	}

	// Handle no rows found
	if errors.Is(err, pgx.ErrNoRows) {
		return errors.New("record not found")
	}

	// Handle PostgreSQL errors using pgconn
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return errors.New(formatUniqueConstraint(pgErr.ConstraintName))
		case "23503": // foreign_key_violation
			return errors.New("referenced record does not exist")
		case "23502": // not_null_violation
			field := formatFieldName(pgErr.ColumnName)
			return errors.New(field + " is required")
		case "23514": // check_violation
			return errors.New("invalid data provided")
		case "42P01": // undefined_table
			return errors.New("database table not found")
		default:
			return errors.New("database operation failed")
		}
	}

	// Return generic message for other errors
	return errors.New("an error occurred")
}

// formatUniqueConstraint converts constraint names to readable field names
func formatUniqueConstraint(constraint string) string {
	if constraint == "" {
		return "Record already exists"
	}

	// Handle common patterns: table_field_key
	if strings.Contains(constraint, "_") {
		parts := strings.Split(constraint, "_")
		if len(parts) >= 2 {
			field := parts[len(parts)-2] // Get field name (second to last part)
			switch field {
			case "email":
				return "Email already exists"
			case "username":
				return "Username already exists"
			case "phone":
				return "Phone number already exists"
			default:
				return formatFieldName(field) + " already exists"
			}
		}
	}

	return "Record already exists"
}

// formatFieldName converts database field names to readable format
func formatFieldName(field string) string {
	if field == "" {
		return "Field"
	}

	// Convert snake_case to Title Case
	words := strings.Split(field, "_")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}

	return strings.Join(words, " ")
}
