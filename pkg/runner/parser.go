package runner

import (
	"fmt"
	"strconv"

	"github.com/codingconcepts/di/pkg/model"
)

func csvLineToArgs(header, line []string, types model.ColumnTypes) ([]any, error) {
	args := make([]any, len(line))

	for i, raw := range line {
		parsed, err := toGoDataType(types[header[i]], raw)
		if err != nil {
			return nil, fmt.Errorf("converting string to go data type: %w", err)
		}

		args[i] = parsed
	}

	return args, nil
}

func toGoDataType(c *model.Column, rawValue string) (any, error) {
	switch c.Type {
	case "uuid", "text", "crdb_internal_region":
		return fmt.Sprintf("%v", rawValue), nil
	case "int8", "int16", "int32", "int64":
		return strconv.ParseInt(rawValue, 10, 64)
	case "bool":
		return strconv.ParseBool(rawValue)
	default:
		return nil, fmt.Errorf("invalid data type: %q", c.Type)
	}
}
