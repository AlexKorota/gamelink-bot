package service

import (
	"errors"
	"fmt"
	"gamelinkBot/prot"
	"log"
	"regexp"
)

var ageRegexp, idRegexp, sexRegexp, delRegexp, registrationRegexp *regexp.Regexp
var err error

func init() {
	ageRegexp, err = regexp.Compile("(((age)\\s*(=\\s*([0-9]{1,2}$)|\\[\\s*([0-9]{1,2})\\s*;\\s*([0-9]{1,2})\\s*\\]$)))")
	if err != nil {
		log.Fatal(err)
	}
	idRegexp, err = regexp.Compile("(((id|vk_id|fb_id)\\s*(=\\s*([0-9]{1,20}$)|\\[\\s*([0-9]{1,20})\\s*;\\s*([0-9]{1,20})\\s*\\]$)))")
	if err != nil {
		log.Fatal(err)
	}
	sexRegexp, err = regexp.Compile("(((sex)\\s*(=\\s*(f|m)$)))")
	if err != nil {
		log.Fatal(err)
	}
	delRegexp, err = regexp.Compile("(((deleted)\\s*(=\\s*(0|1)$)))")
	if err != nil {
		log.Fatal(err)
	}
	registrationRegexp, err = regexp.Compile("(((created_at)\\s*(=\\s*((0[1-9]|1[0-9]|2[0-9]|3[01])\\.(0[1-9]|1[012])\\.[0-9]{4}$)|\\[\\s*((0[1-9]|1[0-9]|2[0-9]|3[01])\\.(0[1-9]|1[012])\\.[0-9]{4})\\s*;\\s*((0[1-9]|1[0-9]|2[0-9]|3[01])\\.(0[1-9]|1[012])\\.[0-9]{4})\\]$)))") //(((created_at)\s*(=\s*((0[1-9]|1[0-9]|2[0-9]|3[01])\.(0[1-9]|1[012])\.[0-9]{4}$)|\[\s*((0[1-9]|1[0-9]|2[0-9]|3[01])\.(0[1-9]|1[012])\.[0-9]{4})\s*;\s*((0[1-9]|1[0-9]|2[0-9]|3[01])\.(0[1-9]|1[012])\.[0-9]{4})\]$)))
	if err != nil {
		log.Fatal(err)
	}
}

func ParseRequest(params []string) ([]*prot.OneCriteriaStruct, error) {
	var multiCriteria []*prot.OneCriteriaStruct
	for _, v := range params {
		var matches []string
		matches = ageRegexp.FindStringSubmatch(v)
		if matches != nil {
			appendToMultiCriteria(&multiCriteria, matches)
			continue
		}
		matches = idRegexp.FindStringSubmatch(v)
		if matches != nil {
			appendToMultiCriteria(&multiCriteria, matches)
			continue
		}
		matches = sexRegexp.FindStringSubmatch(v)
		if matches != nil {
			appendToMultiCriteria(&multiCriteria, matches)
			continue
		}
		matches = delRegexp.FindStringSubmatch(v)
		if matches != nil {
			appendToMultiCriteria(&multiCriteria, matches)
			continue
		}
		matches = registrationRegexp.FindStringSubmatch(v)
		if matches != nil {
			if matches[8] != "" && matches[11] != "" {
				matches[6] = matches[8]
				matches[7] = matches[11]
			}
			appendToMultiCriteria(&multiCriteria, matches)
			continue
		}
		return nil, errors.New(fmt.Sprintf("wrong param %s", v))
	}
	return multiCriteria, nil
}

func appendToMultiCriteria(multiCriteria *[]*prot.OneCriteriaStruct, matches []string) {
	var criteria, secondCriteria prot.OneCriteriaStruct
	if matches[3] != "" {
		if val, ok := prot.OneCriteriaStruct_Criteria_value[matches[3]]; ok {
			criteria.Cr = prot.OneCriteriaStruct_Criteria(val)
			secondCriteria.Cr = prot.OneCriteriaStruct_Criteria(val)
		} else {
			// Стоит ли тут добавить обработку ошибки на случай, если критерий не нашелся в енуме?
		}
	}
	if matches[5] != "" {
		criteria.Op = prot.OneCriteriaStruct_Option(2)
		criteria.Value = matches[5]
		*multiCriteria = append(*multiCriteria, &criteria)
	} else if matches[6] != "" && matches[7] != "" {
		criteria.Op = prot.OneCriteriaStruct_Option(1)
		criteria.Value = matches[7]

		*multiCriteria = append(*multiCriteria, &criteria)

		secondCriteria.Op = prot.OneCriteriaStruct_Option(3)
		secondCriteria.Value = matches[6]

		*multiCriteria = append(*multiCriteria, &secondCriteria)
	}
}
