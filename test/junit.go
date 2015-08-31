package test

import (
	"github.com/coduno/api/model"
	"golang.org/x/net/context"
)

func init() {
	RegisterTester(Junit, junit)
}

func junit(ctx context.Context, params map[string]string, sub model.KeyedSubmission) error {
	return nil
}
