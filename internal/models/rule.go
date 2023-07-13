package models

import (
	"errors"
	"fmt"
)

func RuleFactory(ruleType ModifyRuleType, field string, toString string) (Rule, error) {
	if !getPossibleModifyRuleTypes()[ruleType] {
		return nil, fmt.Errorf("unknown rule type %s", ruleType)
	}

	switch ruleType {
	case ModifyRuleTypeIncrement:
		return &IncrementRule{
			Field: field,
		}, nil
	case ModifyRuleTypeReplace:
		return &ReplaceRule{
			Field:    field,
			ToString: toString,
		}, nil
	}

	return nil, fmt.Errorf("unknown rule %s", ruleType)
}

type IncrementRule struct {
	Field string
}

func (ir IncrementRule) Validate() error {
	var err error

	if ir.Field == "" {
		err = errors.Join(err, fmt.Errorf("incrementRule.Field can not be empty"))
	}

	return err
}

func (ir IncrementRule) Apply(data map[string]any) error {
	value, ok := data[ir.Field]
	if !ok {
		return fmt.Errorf("field '%s' not exists", ir.Field)
	}

	switch value.(type) {
	case int:
		newValue := value.(int)
		newValue++
		data[ir.Field] = newValue
	case int8:
		newValue := value.(int8)
		newValue++
		data[ir.Field] = newValue
	case int16:
		newValue := value.(int16)
		newValue++
		data[ir.Field] = newValue
	case int32:
		newValue := value.(int32)
		newValue++
		data[ir.Field] = newValue
	case int64:
		newValue := value.(int64)
		newValue++
		data[ir.Field] = newValue
	case float32:
		newValue := value.(float32)
		newValue++
		data[ir.Field] = newValue
	case float64:
		newValue := value.(float64)
		newValue++
		data[ir.Field] = newValue
	}

	return nil
}

type ReplaceRule struct {
	Field    string
	ToString string
}

func (rr ReplaceRule) Validate() error {
	var err error

	if rr.Field == "" {
		err = errors.Join(err, fmt.Errorf("replaceRule.Field can not be empty"))
	}

	return err
}

func (rr ReplaceRule) Apply(data map[string]any) error {
	value, ok := data[rr.Field]
	if !ok {
		return fmt.Errorf("field '%s' not exists", rr.Field)
	}

	_, ok = value.(string)
	if !ok {
		return fmt.Errorf("field '%s' not string", rr.Field)
	}

	data[rr.Field] = rr.ToString

	return nil
}
