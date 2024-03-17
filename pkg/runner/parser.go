package runner

import (
	"fmt"
	"strconv"
	"time"

	"github.com/codingconcepts/di/pkg/model"
)

func (r *Runner) csvLineToArgs(header, line []string) ([]any, error) {
	args := make([]any, len(line))

	for i, raw := range line {
		parsed, err := r.toGoDataType(r.types[header[i]], raw)
		if err != nil {
			return nil, fmt.Errorf("converting string to go data type: %w", err)
		}

		args[i] = parsed
	}

	return args, nil
}

func (r *Runner) toGoDataType(c *model.Column, rawValue string) (any, error) {
	switch c.Type {
	case "uuid", "text", "crdb_internal_region":
		return fmt.Sprintf("%v", rawValue), nil
	case "int2", "int4", "int8":
		return strconv.ParseInt(rawValue, 10, 64)
	case "float4", "float8", "numeric":
		return strconv.ParseFloat(rawValue, 64)
	case "bool":
		return strconv.ParseBool(rawValue)
	case "date", "time", "timestamptz":
		return time.Parse(r.formatHelpers[c.Type], rawValue)
	default:
		return nil, fmt.Errorf("invalid data type: %q", c.Type)
	}
}
