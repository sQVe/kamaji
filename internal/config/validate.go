package config

import (
	"fmt"
	"strings"

	"github.com/sqve/kamaji/internal/domain"
)

// ValidationError represents a single validation error with field path and message.
type ValidationError struct {
	Field   string
	Message string
}

// ValidateSprint validates a sprint configuration and returns all validation errors.
func ValidateSprint(s *domain.Sprint) []ValidationError {
	var errors []ValidationError

	errors = validateRequired("name", s.Name, errors)

	for i, ticket := range s.Tickets {
		ticketPrefix := fmt.Sprintf("tickets[%d]", i)
		errors = validateRequired(ticketPrefix+".name", ticket.Name, errors)
		errors = validateNotEmpty(ticketPrefix+".description", ticket.Description, errors)

		for j, task := range ticket.Tasks {
			taskField := fmt.Sprintf("%s.tasks[%d].description", ticketPrefix, j)
			errors = validateRequired(taskField, task.Description, errors)
			if task.Description != "" {
				errors = validateNotEmpty(taskField, task.Description, errors)
			}
		}
	}

	return errors
}

func validateRequired(field, value string, errors []ValidationError) []ValidationError {
	if value == "" {
		return append(errors, ValidationError{
			Field:   field,
			Message: "required",
		})
	}
	return errors
}

func validateNotEmpty(field, value string, errors []ValidationError) []ValidationError {
	if strings.TrimSpace(value) == "" && value != "" {
		return append(errors, ValidationError{
			Field:   field,
			Message: "cannot be empty",
		})
	}
	return errors
}
