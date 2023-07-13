package models

import (
	"context"
	"net/url"
)

type FetcherUsecase interface {
	Fetch(ctx context.Context, targetAPI *url.URL, modifyRules ModifyRules) (map[string]any, error)
}

type Rule interface {
	Apply(data map[string]any) error
	Validate() error
}

type ModifyRules []Rule

func (mRs ModifyRules) ApplyTargetAPIData(data map[string]any) {
	for _, rule := range mRs {
		rule.Apply(data)
	}
}

type ModifyRuleType string

const (
	ModifyRuleTypeIncrement ModifyRuleType = "ModifyRuleTypeIncrement"
	ModifyRuleTypeReplace   ModifyRuleType = "ModifyRuleTypeReplace"
)

func getPossibleModifyRuleTypes() map[ModifyRuleType]bool {
	return map[ModifyRuleType]bool{
		ModifyRuleTypeIncrement: true,
		ModifyRuleTypeReplace:   true,
	}
}
