package service

import (
	"errors"
	"fmt"
	"gamelinkBot/prot"
	"regexp"
	"strings"
)

func ParseRequest(params []string) ([]*prot.OneCriteriaStruct, error) {
	var criteriaNameEnum int32
	var criteria prot.OneCriteriaStruct
	var multiCriteria []*prot.OneCriteriaStruct
	for _, v := range params {
		strings.Trim(v, " ")
		arr := strings.Split(v, "=")
		if val, ok := prot.OneCriteriaStruct_Criteria_value[arr[0]]; !ok {
			return nil, errors.New(fmt.Sprintf("Invalid criteria %s", arr[0]))
		} else {
			criteriaNameEnum = val
		}
		rangeSearch, err := regexp.MatchString("^[[0-9]*;[0-9]*]$", arr[1])
		if err != nil {
			return nil, err
		}
		oneValueSearch, err := regexp.MatchString("^[A-z0-9А-яЁё]*$", arr[1])
		if err != nil {
			return nil, err
		}
		if rangeSearch {
			arr := strings.Split(arr[1], ";")
			criteria.Cr = prot.OneCriteriaStruct_Criteria(criteriaNameEnum)
			criteria.Op = prot.OneCriteriaStruct_Option(2)
			criteria.Value = strings.Trim(arr[0], "[]")
			multiCriteria = append(multiCriteria, &criteria)
			criteria.Op = prot.OneCriteriaStruct_Option(0)
			criteria.Value = strings.Trim(arr[1], "[]")
			multiCriteria = append(multiCriteria, &criteria)
		} else if oneValueSearch {
			criteria.Cr = prot.OneCriteriaStruct_Criteria(criteriaNameEnum)
			criteria.Op = prot.OneCriteriaStruct_Option(1)
			criteria.Value = arr[1]
			multiCriteria = append(multiCriteria, &criteria)
		} else {
			return nil, errors.New(fmt.Sprintf("Invalid criteria value %s", arr[1]))
		}
	}

	return multiCriteria, nil
}
