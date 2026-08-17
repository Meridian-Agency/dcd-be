package service

import (
	"regexp"
	"strings"
)

type CoverageService interface {
	CheckCoverage(postcode string) (covered bool, cleanPostcode string, message string)
}

type coverageService struct {
	nottinghamOutcodes map[string]bool
	postcodeRegex      *regexp.Regexp
}

func NewCoverageService() CoverageService {
	return &coverageService{
		nottinghamOutcodes: map[string]bool{
			"NG1": true, "NG2": true, "NG3": true, "NG4": true,
			"NG5": true, "NG6": true, "NG7": true, "NG8": true,
			"NG9": true, "NG15": true, "NG16": true, "DE7": true,
		},
		postcodeRegex: regexp.MustCompile(`^([A-Z]{1,2}\d{1,2}[A-Z]?)`),
	}
}

func (s *coverageService) CheckCoverage(postcode string) (bool, string, string) {
	clean := strings.ToUpper(strings.ReplaceAll(postcode, " ", ""))
	match := s.postcodeRegex.FindStringSubmatch(clean)

	if len(match) < 2 {
		return false, postcode, "Invalid UK postcode format."
	}

	outcode := match[1]
	if s.nottinghamOutcodes[outcode] {
		return true, postcode, "Your area is covered for mobile valeting!"
	}

	return false, postcode, "Outside mobile coverage area. Studio drop-off in Bulwell available."
}
