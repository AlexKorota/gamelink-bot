package service

import (
	"errors"
	"fmt"
	msg "gamelink-go/proto_msg"
	"gamelinkBot/command_list"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	ageRegexp, idRegexp, sexRegexp, delRegexp, registrationRegexp, permissionRegexp, pushRegexp, updRegexp, updatedAtRegexp, adsRegexp, paymentRegexp *regexp.Regexp
	UnknownCommandError                                                                                                                               error
)

func init() {
	var err error
	UnknownCommandError = errors.New("Unknown command")
	ageRegexp, err = regexp.Compile("(((^age)\\s*(=\\s*([0-9]{1,2}$)|\\[\\s*((([0-9]{1,2})))\\s*;\\s*((([0-9]{1,2})))\\s*\\]$)))")
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

	adsRegexp, err = regexp.Compile("(((watched_ads)\\s*(=\\s*(0|1)$)))")
	if err != nil {
		log.Fatal(err)
	}

	paymentRegexp, err = regexp.Compile("(((made_payment)\\s*(=\\s*(0|1)$)))")
	if err != nil {
		log.Fatal(err)
	}

	registrationRegexp, err = regexp.Compile("(((created_at)\\s*(=\\s*((0?[1-9]|1[0-9]|2[0-9]|3[01])\\.(0?[1-9]|1[012])\\.[0-9]{4}$)|\\[\\s*((0?[1-9]|1[0-9]|2[0-9]|3[01])\\.(0?[1-9]|1[012])\\.[0-9]{4})\\s*;\\s*((0?[1-9]|1[0-9]|2[0-9]|3[01])\\.(0?[1-9]|1[012])\\.[0-9]{4})\\]$)))") //(((created_at)\s*(=\s*((0[1-9]|1[0-9]|2[0-9]|3[01])\.(0[1-9]|1[012])\.[0-9]{4}$)|\[\s*((0[1-9]|1[0-9]|2[0-9]|3[01])\.(0[1-9]|1[012])\.[0-9]{4})\s*;\s*((0[1-9]|1[0-9]|2[0-9]|3[01])\.(0[1-9]|1[012])\.[0-9]{4})\]$)))
	if err != nil {
		log.Fatal(err)
	}
	permissionRegexp, err = regexp.Compile("(\\w+)\\s*(\\[((\\s*(count|find|delete|send_push|update|get_user)\\s*;)*\\s*(count|find|delete|send_push|update|get_user))\\s*])?")
	if err != nil {
		log.Fatal(err)
	}
	pushRegexp, err = regexp.Compile("(((message)))\\s*{{1}\\s*((.+))}{1}")
	if err != nil {
		log.Fatal(err)
	}
	updRegexp, err = regexp.Compile("(set|delete)\\[(vk_id|fb_id|sex|age|country|deleted)](=(.+))?")
	if err != nil {
		log.Fatal(err)
	}
	updatedAtRegexp, err = regexp.Compile("(((updated_at)(=((0?[1-9]|1[0-9]|2[0-9]|3[01])\\.(0?[1-9]|1[012])\\.[0-9]{4}$)|\\[((0?[1-9]|1[0-9]|2[0-9]|3[01])\\.(0?[1-9]|1[012])\\.[0-9]{4});((0?[1-9]|1[0-9]|2[0-9]|3[01])\\.(0?[1-9]|1[012])\\.[0-9]{4})\\]$)))")
	if err != nil {
		log.Fatal(err)
	}

}

func ParseRequest(params []string, cmd string) ([]*msg.OneCriteriaStruct, []*msg.UpdateCriteriaStruct, string, error) {
	var multiCriteria []*msg.OneCriteriaStruct
	var updateCriteria []*msg.UpdateCriteriaStruct
	var message string
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
		matches = adsRegexp.FindStringSubmatch(v)
		if matches != nil {
			appendToMultiCriteria(&multiCriteria, matches)
			continue
		}
		matches = paymentRegexp.FindStringSubmatch(v)
		if matches != nil {
			appendToMultiCriteria(&multiCriteria, matches)
			continue
		}
		matches = registrationRegexp.FindStringSubmatch(v)
		if matches != nil {
			matches, err := appendTimeParams(matches)
			if err != nil {
				return nil, nil, "", err
			}
			appendToMultiCriteria(&multiCriteria, matches)
			continue
		}
		matches = updatedAtRegexp.FindStringSubmatch(v)
		if matches != nil {
			matches, err := appendTimeParams(matches)
			if err != nil {
				return nil, nil, "", err
			}
			appendToMultiCriteria(&multiCriteria, matches)
			continue
		}
		//This check only for send push command
		if cmd == "/"+command_list.CommandSendPush {
			matches = pushRegexp.FindStringSubmatch(v)
			if matches != nil {
				message = matches[5]
				continue
			}
		}
		//This check only for update command
		if cmd == "/"+command_list.CommandUpdate {
			matches = updRegexp.FindStringSubmatch(v)
			if matches != nil {
				appendToUpdateCriteria(&updateCriteria, matches)
				continue
			}
		}
		return nil, nil, "", errors.New(fmt.Sprintf("wrong param %s", v))
	}
	return multiCriteria, updateCriteria, message, nil
}

func appendToMultiCriteria(multiCriteria *[]*msg.OneCriteriaStruct, matches []string) {
	var criteria, secondCriteria msg.OneCriteriaStruct
	if matches[3] != "" {
		if val, ok := msg.OneCriteriaStruct_Criteria_value[matches[3]]; ok {
			criteria.Cr = msg.OneCriteriaStruct_Criteria(val)
			secondCriteria.Cr = msg.OneCriteriaStruct_Criteria(val)
		}
	}
	if matches[5] != "" {
		criteria.Op = msg.OneCriteriaStruct_e
		criteria.Value = matches[5]
		*multiCriteria = append(*multiCriteria, &criteria)
	} else if matches[8] != "" && matches[11] != "" {
		criteria.Op = msg.OneCriteriaStruct_l
		criteria.Value = matches[11]

		*multiCriteria = append(*multiCriteria, &criteria)

		secondCriteria.Op = msg.OneCriteriaStruct_g
		secondCriteria.Value = matches[8]
		*multiCriteria = append(*multiCriteria, &secondCriteria)
	}
}

func appendToUpdateCriteria(updateCriteria *[]*msg.UpdateCriteriaStruct, matches []string) {
	var criteria msg.UpdateCriteriaStruct
	if matches[2] != "" {
		if val, ok := msg.UpdateCriteriaStruct_UpdCriteria_value[matches[2]]; ok {
			criteria.Ucr = msg.UpdateCriteriaStruct_UpdCriteria(val)
		}
	}
	if matches[1] == msg.UpdateCriteriaStruct_set.String() {
		criteria.Uop = msg.UpdateCriteriaStruct_set
	} else if matches[1] == msg.UpdateCriteriaStruct_delete.String() {
		criteria.Uop = msg.UpdateCriteriaStruct_delete
	}
	if matches[4] != "" {
		criteria.Value = matches[4]
	}
	*updateCriteria = append(*updateCriteria, &criteria)
}

func CompareParseCommand(str, cmd string) ([]*msg.OneCriteriaStruct, []*msg.UpdateCriteriaStruct, string, error) {
	ind := strings.Index(str, " ")
	if ind < 0 || str[:ind] != cmd {
		return nil, nil, "", UnknownCommandError
	}
	messageInd := strings.Index(str, "message")
	var params []string
	if messageInd < 0 {
		params = strings.Split(str[ind+1:], " ")
	} else {
		params = strings.Split(str[ind+1:messageInd], " ")
		params = append(params, str[messageInd:])
	}

	return ParseRequest(params, cmd)
}

func CompareParsePermissionCommand(str, cmd string) (string, []string, error) {
	ind := strings.Index(str, " ")
	if ind < 0 || str[:ind] != cmd {
		return "", nil, UnknownCommandError
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

func appendTimeParams(matches []string) ([]string, error) {
	if matches[5] != "" {
		t, err := stringToTime(matches[5])
		if err != nil {
			return nil, err
		}
		s := time.Date(t.Year(), t.Month(), t.Day(), 00, 00, 00, 00, time.Local).Unix()
		matches[8] = strconv.Itoa(int(s))

		e := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 59, time.Local).Unix()
		matches[11] = strconv.Itoa(int(e))
		matches[5] = ""
	} else if matches[8] != "" && matches[11] != "" {
		t1, err := stringToTime(matches[8])
		if err != nil {
			return nil, err
		}
		s := time.Date(t1.Year(), t1.Month(), t1.Day(), 00, 00, 00, 00, time.Local).Unix()
		matches[8] = strconv.Itoa(int(s))

		t2, err := stringToTime(matches[11])
		if err != nil {
			return nil, err
		}
		e := time.Date(t2.Year(), t2.Month(), t2.Day(), 23, 59, 59, 59, time.Local).Unix()
		matches[11] = strconv.Itoa(int(e))
	}
	return matches, nil
}

func stringToTime(date string) (time.Time, error) {
	lay := "2.1.2006"
	t, err := time.ParseInLocation(lay, date, time.Local)
	if err != nil {
		return t.AddDate(0, 0, 0), err
	}
	return t, nil
}
