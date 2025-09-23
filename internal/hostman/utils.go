package hostman

import (
	"errors"
	"fmt"
	"strings"
)

func InsertOrReplaceSection(
	data string,
	start string,
	end string,
	content string,
) (string, error) {
	iStart := strings.Index(data, start)
	iEnd := strings.Index(data, end)

	if (iStart == -1 || iEnd == -1) && (iStart != iEnd) {
		return "", errors.New("There must either be no start or end section, or both must be present")
	}

	nStart := fmt.Sprintln(start)
	nContent := fmt.Sprintln(content)
	nEnd := fmt.Sprintln(end)

	if iStart == -1 && iEnd == -1 {
		return data + nStart + nContent + nEnd, nil
	}

	before := data[:iStart]
	after := data[iEnd+len(end):]

	if strings.HasPrefix(after, "\n") {
		after = after[1:]
	}

	replacement := nStart + nContent + nEnd
	return before + replacement + after, nil
}
