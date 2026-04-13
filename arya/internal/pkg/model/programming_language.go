package model

import (
	"strings"

	"github.com/CodefriendOrg/arya/internal/pb/github.com/CodefriendOrg/daenerys/pkg/daenerys"
)

const programmingLanguagePrefix = "LANGUAGE_"

// ToProto .
func (e ProgrammingLanguage) ToProto() daenerys.ProgrammingLanguage {
	return daenerys.ProgrammingLanguage(daenerys.ProgrammingLanguage_value[programmingLanguagePrefix+e.String()]) //nolint:unconvert
}

// ConvertProgrammingLanguage .
func ConvertProgrammingLanguage(language daenerys.ProgrammingLanguage) ProgrammingLanguage {
	return ProgrammingLanguage(strings.TrimPrefix(string(language.String()), programmingLanguagePrefix)) //nolint:unconvert
}
