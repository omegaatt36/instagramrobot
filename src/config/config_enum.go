// Code generated by go-enum DO NOT EDIT.
// Version: 0.5.3
// Revision: 8e2c93debfc66888870b2dfd86e70c79a70c920f
// Build Date: 2022-11-09T16:39:46Z
// Built By: goreleaser

package config

import (
	"errors"
	"fmt"
)

const (
	// EnvDevelopment is a env of type development.
	EnvDevelopment env = "development"
	// EnvProduction is a env of type production.
	EnvProduction env = "production"
)

var ErrInvalidenv = errors.New("not a valid env")

// String implements the Stringer interface.
func (x env) String() string {
	return string(x)
}

// String implements the Stringer interface.
func (x env) IsValid() bool {
	_, err := Parseenv(string(x))
	return err == nil
}

var _envValue = map[string]env{
	"development": EnvDevelopment,
	"production":  EnvProduction,
}

// Parseenv attempts to convert a string to a env.
func Parseenv(name string) (env, error) {
	if x, ok := _envValue[name]; ok {
		return x, nil
	}
	return env(""), fmt.Errorf("%s is %w", name, ErrInvalidenv)
}