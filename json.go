package valid

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Validator interface {
	Valid(ctx context.Context) (problems map[string]string)
}

var ErrValidation error
var ErrDecode error

func DecodeJSON[T Validator](r *http.Request) (T, map[string]string, error) {
	var v T

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		ErrDecode = fmt.Errorf("error decoding JSON: %w", err)
		return v, nil, ErrDecode
	}

	if problems := v.Valid(r.Context()); len(problems) != 0 {
		ErrValidation = fmt.Errorf("invalid %T: %d problem(s)", v, len(problems))
		return v, problems, ErrValidation
	}

	return v, nil, nil
}
