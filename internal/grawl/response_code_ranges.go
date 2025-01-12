package grawl

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type ResponseCodeRange struct {
	minCode int
	maxCode int
}

type ResponseCodeRanges struct {
	ranges []ResponseCodeRange
}

func (r *ResponseCodeRanges) IsError(responseCode int) bool {
	for _, errorCodeRange := range r.ranges {
		if responseCode >= errorCodeRange.minCode && responseCode <= errorCodeRange.maxCode {
			return true
		}
	}
	return false
}

func NewResponseCodeRanges(responseCodeFlags []string) (*ResponseCodeRanges, error) {
	var ranges = ResponseCodeRanges{}

	regexRange := regexp.MustCompile(`(\d+)\s*-\s*(\d+)`)

	for _, responseCodeFlag := range responseCodeFlags {
		responseCodeFlag = strings.TrimSpace(responseCodeFlag)
		matches := regexRange.FindStringSubmatch(responseCodeFlag)

		if len(matches) < 3 {
			val, err := strconv.Atoi(responseCodeFlag)
			if err != nil {
				return nil, fmt.Errorf("could not parse http-error-codes: %s", responseCodeFlag)
			}

			resSingle := ResponseCodeRange{
				minCode: val,
				maxCode: val,
			}

			ranges.ranges = append(ranges.ranges, resSingle)
			continue
		}

		minVal, err := strconv.Atoi(strings.TrimSpace(matches[1]))
		if err != nil {
			return nil, err
		}

		maxVal, err := strconv.Atoi(strings.TrimSpace(matches[2]))
		if err != nil {
			return nil, err
		}

		res := ResponseCodeRange{
			minCode: minVal,
			maxCode: maxVal,
		}

		ranges.ranges = append(ranges.ranges, res)
	}

	if len(ranges.ranges) <= 0 {
		ranges.ranges = append(ranges.ranges, ResponseCodeRange{minCode: 400, maxCode: 599})
	}

	return &ranges, nil
}
