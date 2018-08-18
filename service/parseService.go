package service

import (
	"errors"
	"fmt"
	"gamelinkBot/prot"
	"regexp"
)

func ParseRequest(params []string) ([]*prot.OneCriteriaStruct, error) {
	//var criteria prot.OneCriteriaStruct
	var multiCriteria []*prot.OneCriteriaStruct
	ageRegexp, err := regexp.Compile("(((age)\\s*(=\\s*([0-9]{1,2}$)|\\[\\s*([0-9]{1,2})\\s*;\\s*([0-9]{1,2})\\s*\\]$)))")
	if err != nil {
		return nil, errors.New("age regexp error")
	}
	idRegexp, err := regexp.Compile("(((id|vk_id|fb_id)\\s*(=\\s*([0-9]{1,20}$)|\\[\\s*([0-9]{1,20})\\s*;\\s*([0-9]{1,20})\\s*\\]$)))")
	if err != nil {
		return nil, errors.New("id regexp error")
	}
	sexRegexp, err := regexp.Compile("(((sex)\\s*(=\\s*(f|m)$)))")
	if err != nil {
		return nil, errors.New("sex regexp error")
	}
	delRegexp, err := regexp.Compile("(((deleted)\\s*(=\\s*(0|1)$)))")
	if err != nil {
		return nil, errors.New("delete regexp error")
	}
	registrationRegexp, err := regexp.Compile("(((created_at)\\s*(=\\s*((0[1-9]|1[0-9]|2[0-9]|3[01])\\.(0[1-9]|1[012])\\.[0-9]{4}$)|\\[\\s*((0[1-9]|1[0-9]|2[0-9]|3[01])\\.(0[1-9]|1[012])\\.[0-9]{4})\\s*;\\s*((0[1-9]|1[0-9]|2[0-9]|3[01])\\.(0[1-9]|1[012])\\.[0-9]{4})\\]$)))") //(((created_at)\s*(=\s*((0[1-9]|1[0-9]|2[0-9]|3[01])\.(0[1-9]|1[012])\.[0-9]{4}$)|\[\s*((0[1-9]|1[0-9]|2[0-9]|3[01])\.(0[1-9]|1[012])\.[0-9]{4})\s*;\s*((0[1-9]|1[0-9]|2[0-9]|3[01])\.(0[1-9]|1[012])\.[0-9]{4})\]$)))
	if err != nil {
		return nil, errors.New("delete regexp error")
	}
	for _, v := range params {
		var result []string
		result = ageRegexp.FindStringSubmatch(v)
		if result != nil {
			appendToMultiCriteria(&multiCriteria, result)
			continue
		}
		result = idRegexp.FindStringSubmatch(v)
		if result != nil {
			appendToMultiCriteria(&multiCriteria, result)
			continue
		}
		result = sexRegexp.FindStringSubmatch(v)
		if result != nil {
			appendToMultiCriteria(&multiCriteria, result)
			continue
		}
		result = delRegexp.FindStringSubmatch(v)
		if result != nil {
			appendToMultiCriteria(&multiCriteria, result)
			continue
		}
		result = registrationRegexp.FindStringSubmatch(v)
		if result != nil {
			if result[8] != "" && result[11] != "" {
				result[6] = result[8]
				result[7] = result[11]
			}
			appendToMultiCriteria(&multiCriteria, result)
			continue
		}
		return nil, errors.New(fmt.Sprintf("wrong param %s", v))
	}
	fmt.Println(multiCriteria)
	return multiCriteria, nil
}

func appendToMultiCriteria(multiCriteria *[]*prot.OneCriteriaStruct, result []string) {
	var criteria, secondCriteria prot.OneCriteriaStruct
	fmt.Println(result)
	if result[3] != "" {
		if val, ok := prot.OneCriteriaStruct_Criteria_value[result[3]]; ok {
			criteria.Cr = prot.OneCriteriaStruct_Criteria(val)
			secondCriteria.Cr = prot.OneCriteriaStruct_Criteria(val)
		} else {
			// Стоит ли тут добавить обработку ошибки на случай, если критерий не нашелся в енуме?
		}
	}
	if result[5] != "" {
		criteria.Op = prot.OneCriteriaStruct_Option(2)
		criteria.Value = result[5]
		*multiCriteria = append(*multiCriteria, &criteria)
	} else if result[6] != "" && result[7] != "" {
		criteria.Op = prot.OneCriteriaStruct_Option(1)
		criteria.Value = result[7]

		*multiCriteria = append(*multiCriteria, &criteria)

		secondCriteria.Op = prot.OneCriteriaStruct_Option(3)
		secondCriteria.Value = result[6]

		*multiCriteria = append(*multiCriteria, &secondCriteria)
	}
}
