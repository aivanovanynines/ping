package app

import (
	"fmt"

	"github.com/aivanovanynines/ping/internal/models"
)

type fetcherRequest struct {
	TargetAPI string `json:"target_api"`
	Rules     []rule `json:"rules"`
}

type rule struct {
	Field    string `json:"field"`
	Type     string `json:"type"`
	ToString string `json:"to_string"`
}

func (fr *fetcherRequest) getRules() (models.ModifyRules, error) {
	modifyRules := make(models.ModifyRules, 0, len(fr.Rules))
	for _, requestRule := range fr.Rules {
		r, err := models.RuleFactory(
			models.ModifyRuleType(requestRule.Type),
			requestRule.Field,
			requestRule.ToString,
		)
		if err != nil {
			return nil, fmt.Errorf("create rule %s. %w", requestRule.Type, err)
		}

		err = r.Validate()
		if err != nil {
			return nil, fmt.Errorf("validation error %w", err)
		}

		modifyRules = append(modifyRules, r)
	}

	return modifyRules, nil
}
