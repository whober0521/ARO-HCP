package poc

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type VarHandler interface {
	ReplaceVariables(variables Variables) (Variables, error)
}

type varHandler struct {
	prefix           string // usually the '$'
	quoteSymbolLeft  string // usually the '(' or '{' or '{{'
	quoteSymbolRight string // usually the '(' or '{' or '{{'
	region           string
	user             string
}

func NewDefaultVarHandler(region, user string) VarHandler {
	return &varHandler{
		prefix:           "$",
		quoteSymbolLeft:  "(",
		quoteSymbolRight: ")",
		region:           region,
		user:             user,
	}
}

// replace variables values recursively
// variable in variable is not supported yet, sample: ${${AAA}}
func (vh *varHandler) ReplaceVariables(variables Variables) (Variables, error) {
	replacedVariables := Variables{}
	pendingVariables := Variables{}
	re := regexp.MustCompile(fmt.Sprintf(`\%s\%s[^%s]+\%s`, vh.prefix, vh.quoteSymbolLeft, vh.quoteSymbolRight, vh.quoteSymbolRight))
	for k, v := range variables {
		if !re.MatchString(v) {
			replacedVariables[k] = v
		} else {
			pendingVariables[k] = v
		}
	}

	if len(replacedVariables) == 0 {
		return nil, errors.New("no raw variables values found")
	}

	replacedVariables["region"] = vh.region
	replacedVariables["user"] = vh.user

	var err error
	variableReplace := func(origin string) string {
		trimedName := strings.TrimSuffix(strings.TrimPrefix(origin, fmt.Sprintf("%s%s", vh.prefix, vh.quoteSymbolLeft)), vh.quoteSymbolRight)
		replacedValue, ok := replacedVariables[trimedName]
		if ok {
			return replacedValue
		} else {
			return origin
		}
	}

	// currently we only support one level replacing
	needToContrinue := true
	maxLevel := 100
	currentLevel := 0
	for needToContrinue {
		needToContrinue = false
		currentLevel++
		if currentLevel >= maxLevel {
			err = errors.New("max recursive level reached")
			break
		}
		itemToRemove := []string{}
		for k, v := range pendingVariables {
			// if the value does not contains
			replacedVariable := re.ReplaceAllStringFunc(v, variableReplace)

			if !re.MatchString(replacedVariable) {
				replacedVariables[k] = replacedVariable
				itemToRemove = append(itemToRemove, k)
				needToContrinue = true
			} else {
				pendingVariables[k] = replacedVariable
			}
		}
		for _, itemK := range itemToRemove {
			delete(pendingVariables, itemK)
		}
	}

	return replacedVariables, err
}
