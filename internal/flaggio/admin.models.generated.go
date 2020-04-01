// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package flaggio

import (
	"fmt"
	"io"
	"strconv"
)

type Ruler interface {
	IsRuler()
}

type FlagResults struct {
	Flags []*Flag `json:"flags"`
	Total int     `json:"total"`
}

type NewConstraint struct {
	Property  string        `json:"property"`
	Operation Operation     `json:"operation"`
	Values    []interface{} `json:"values"`
}

type NewDistribution struct {
	VariantID  string `json:"variantId"`
	Percentage int    `json:"percentage"`
}

type NewFlag struct {
	Key         string  `json:"key"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type NewFlagRule struct {
	Constraints   []*NewConstraint   `json:"constraints"`
	Distributions []*NewDistribution `json:"distributions"`
}

type NewSegment struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type NewSegmentRule struct {
	Constraints []*NewConstraint `json:"constraints"`
}

type NewVariant struct {
	Description *string     `json:"description"`
	Value       interface{} `json:"value"`
}

type UpdateFlag struct {
	Key                   *string `json:"key"`
	Name                  *string `json:"name"`
	Description           *string `json:"description"`
	Enabled               *bool   `json:"enabled"`
	DefaultVariantWhenOn  *string `json:"defaultVariantWhenOn"`
	DefaultVariantWhenOff *string `json:"defaultVariantWhenOff"`
}

type UpdateFlagRule struct {
	Constraints   []*NewConstraint   `json:"constraints"`
	Distributions []*NewDistribution `json:"distributions"`
}

type UpdateSegment struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type UpdateSegmentRule struct {
	Constraints []*NewConstraint `json:"constraints"`
}

type UpdateVariant struct {
	Description *string     `json:"description"`
	Value       interface{} `json:"value"`
}

type Operation string

const (
	OperationOneOf            Operation = "ONE_OF"
	OperationNotOneOf         Operation = "NOT_ONE_OF"
	OperationGreater          Operation = "GREATER"
	OperationGreaterOrEqual   Operation = "GREATER_OR_EQUAL"
	OperationLower            Operation = "LOWER"
	OperationLowerOrEqual     Operation = "LOWER_OR_EQUAL"
	OperationExists           Operation = "EXISTS"
	OperationDoesntExist      Operation = "DOESNT_EXIST"
	OperationContains         Operation = "CONTAINS"
	OperationDoesntContain    Operation = "DOESNT_CONTAIN"
	OperationStartsWith       Operation = "STARTS_WITH"
	OperationDoesntStartWith  Operation = "DOESNT_START_WITH"
	OperationEndsWith         Operation = "ENDS_WITH"
	OperationDoesntEndWith    Operation = "DOESNT_END_WITH"
	OperationMatchesRegex     Operation = "MATCHES_REGEX"
	OperationDoesntMatchRegex Operation = "DOESNT_MATCH_REGEX"
	OperationIsInSegment      Operation = "IS_IN_SEGMENT"
	OperationIsntInSegment    Operation = "ISNT_IN_SEGMENT"
	OperationIsInNetwork      Operation = "IS_IN_NETWORK"
)

var AllOperation = []Operation{
	OperationOneOf,
	OperationNotOneOf,
	OperationGreater,
	OperationGreaterOrEqual,
	OperationLower,
	OperationLowerOrEqual,
	OperationExists,
	OperationDoesntExist,
	OperationContains,
	OperationDoesntContain,
	OperationStartsWith,
	OperationDoesntStartWith,
	OperationEndsWith,
	OperationDoesntEndWith,
	OperationMatchesRegex,
	OperationDoesntMatchRegex,
	OperationIsInSegment,
	OperationIsntInSegment,
	OperationIsInNetwork,
}

func (e Operation) IsValid() bool {
	switch e {
	case OperationOneOf, OperationNotOneOf, OperationGreater, OperationGreaterOrEqual, OperationLower, OperationLowerOrEqual, OperationExists, OperationDoesntExist, OperationContains, OperationDoesntContain, OperationStartsWith, OperationDoesntStartWith, OperationEndsWith, OperationDoesntEndWith, OperationMatchesRegex, OperationDoesntMatchRegex, OperationIsInSegment, OperationIsntInSegment, OperationIsInNetwork:
		return true
	}
	return false
}

func (e Operation) String() string {
	return string(e)
}

func (e *Operation) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Operation(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Operation", str)
	}
	return nil
}

func (e Operation) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
