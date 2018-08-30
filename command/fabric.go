package command

import (
	"context"
	"errors"
	"fmt"
	"gamelinkBot/prot"
	"log"
	"regexp"
	"strings"
)

var ageRegexp, idRegexp, sexRegexp, delRegexp, registrationRegexp *regexp.Regexp

type CountFabric struct{}

type CountCommand struct {
	params   []*prot.OneCriteriaStruct
	res      Responder
	userName string
}

func init() {
	var err error
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
	SharedParser().RegisterFabric(CountFabric{})
}

func (c CountFabric) RequireAdmin() bool {
	return false
}

func (c CountFabric) Require() []string {
	return []string{"count"}
}

func (c CountFabric) TryParse(req RequesterResponder) (Command, error) {
	var command CountCommand
	ind := strings.Index(req.Request(), " ")
	if req.Request()[:ind] != "/count" {
		return nil, nil
	}
	params := strings.Split(req.Request()[ind+1:], " ")
	for _, v := range params {
		var matches []string
		if v == "" {
			continue
		}
		matches = ageRegexp.FindStringSubmatch(v)
		if matches != nil {
			appendToMultiCriteria(&command.params, matches)
			continue
		}
		matches = idRegexp.FindStringSubmatch(v)
		if matches != nil {
			appendToMultiCriteria(&command.params, matches)
			continue
		}
		matches = sexRegexp.FindStringSubmatch(v)
		if matches != nil {
			appendToMultiCriteria(&command.params, matches)
			continue
		}
		matches = delRegexp.FindStringSubmatch(v)
		if matches != nil {
			appendToMultiCriteria(&command.params, matches)
			continue
		}
		matches = registrationRegexp.FindStringSubmatch(v)
		if matches != nil {
			appendToMultiCriteria(&command.params, matches)
			continue
		}
		return nil, errors.New(fmt.Sprintf("wrong param %s", v))
	}
	command.userName = req.UserName()
	command.res = req
	return command, nil
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
		criteria.Op = prot.OneCriteriaStruct_e
		criteria.Value = matches[5]
		*multiCriteria = append(*multiCriteria, &criteria)
	} else if matches[8] != "" && matches[11] != "" {
		criteria.Op = prot.OneCriteriaStruct_l
		criteria.Value = matches[11]

		*multiCriteria = append(*multiCriteria, &criteria)

		secondCriteria.Op = prot.OneCriteriaStruct_g
		secondCriteria.Value = matches[8]

		*multiCriteria = append(*multiCriteria, &secondCriteria)
	}
}

func (cc CountCommand) Execute(ctx context.Context) {
	r, err := SharedClient().Count(ctx, &prot.MultiCriteriaRequest{cc.params})
	if err != nil {
		cc.res.Respond(err.Error())
		return
	}
	cc.res.Respond(r.String())
}
