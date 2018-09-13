package service

import (
	"errors"
	"fmt"
	"gamelinkBot/iface"
	"log"
	"regexp"
	"strings"
)

var ageRegexp, idRegexp, sexRegexp, delRegexp, registrationRegexp, permissionRegexp *regexp.Regexp
var err error

func init() {
	ageRegexp, err = regexp.Compile("(((age)\\s*(=\\s*([0-9]{1,2}$)|\\[\\s*((([0-9]{1,2})))\\s*;\\s*((([0-9]{1,2})))\\s*\\]$)))")
	if err != nil {
		log.Fatal(err)
	}
	idRegexp, err = regexp.Compile("(((id|vk_id|fb_id)\\s*(=\\s*([0-9]{1,20}$)|\\[\\s*((([0-9]{1,20})))\\s*;\\s*((([0-9]{1,20})))\\s*\\]$)))")
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
	permissionRegexp, err = regexp.Compile("(\\w+)\\s*(\\[((\\s*(count|find|delete|send_push|update|get_user)\\s*;)*\\s*(count|find|delete|send_push|update|get_user))\\s*])?")
	if err != nil {
		log.Fatal(err)
	}
}

func ParseRequest(params []string) ([]*iface.OneCriteriaStruct, error) {
	var multiCriteria []*iface.OneCriteriaStruct
	for _, v := range params {
		var matches []string
		if v == "" {
			continue
		}
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
			appendToMultiCriteria(&multiCriteria, matches)
			continue
		}
		return nil, errors.New(fmt.Sprintf("wrong param %s", v))
	}
	return multiCriteria, nil
}

func appendToMultiCriteria(multiCriteria *[]*iface.OneCriteriaStruct, matches []string) {
	var criteria, secondCriteria iface.OneCriteriaStruct
	if matches[3] != "" {
		if val, ok := iface.OneCriteriaStruct_Criteria_value[matches[3]]; ok {
			criteria.Cr = iface.OneCriteriaStruct_Criteria(val)
			secondCriteria.Cr = iface.OneCriteriaStruct_Criteria(val)
		} else {
			// Стоит ли тут добавить обработку ошибки на случай, если критерий не нашелся в енуме?
		}
	}
	if matches[5] != "" {
		criteria.Op = iface.OneCriteriaStruct_e
		criteria.Value = matches[5]
		*multiCriteria = append(*multiCriteria, &criteria)
	} else if matches[8] != "" && matches[11] != "" {
		criteria.Op = iface.OneCriteriaStruct_l
		criteria.Value = matches[11]

		*multiCriteria = append(*multiCriteria, &criteria)

		secondCriteria.Op = iface.OneCriteriaStruct_g
		secondCriteria.Value = matches[8]

		*multiCriteria = append(*multiCriteria, &secondCriteria)
	}
}

func CompareParseCommand(str, cmd string) ([]*iface.OneCriteriaStruct, error) {
	ind := strings.Index(str, " ")
	if ind < 0 || str[:ind] != cmd {
		return nil, errors.New("wrong command")
	}
	params := strings.Split(str[ind+1:], " ")
	return ParseRequest(params)
}

func CompareParsePermissionCommand(str, cmd string) (string, []string, error) {
	ind := strings.Index(str, " ")
	if ind < 0 || str[:ind] != cmd {
		return "", nil, errors.New("wrong permission command")
	}
	return ParsePermissionRequest(str[ind+1:])
}

func ParsePermissionRequest(params string) (string, []string, error) {
	var matches []string
	matches = permissionRegexp.FindStringSubmatch(params)
	if matches == nil {
		return "", nil, errors.New("bad admin request")
	}
	userName := matches[1]
	permissions := strings.Split(matches[3], ";")
	for k, v := range permissions {
		permissions[k] = strings.Trim(v, " ")
	}
	if matches == nil {
		return "", nil, errors.New("there is no available params")
	}
	return userName, permissions, nil
}
