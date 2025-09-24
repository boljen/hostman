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

// ExtractProjectSections scans the given data and returns all section contents for the
// specified project. Sections are delimited by lines in the form:
//   "### hostman-project-start <project> ###"
// and either
//   "### hostman-project-end <project> ###" or "### hostman-project-stop <project> ###"
// The returned slice contains the raw content between the marker lines (without the
// marker lines themselves). If a start marker is found without a matching end/stop
// marker, an error is returned.
func ExtractProjectSections(data string, project string) ([]string, error) {
	start := "### hostman-project-start " + project + " ###"
	end := "### hostman-project-end " + project + " ###"
	stop := "### hostman-project-stop " + project + " ###"

	var results []string
	searchFrom := 0
	for {
		iStart := strings.Index(data[searchFrom:], start)
		if iStart == -1 {
			break
		}
		iStart += searchFrom

		// Content begins just after the start marker. If the next character is a newline,
		// skip exactly one to start content on the following line (matching how sections
		// are written by InsertOrReplaceSection).
		contentStart := iStart + len(start)
		if contentStart < len(data) && data[contentStart] == '\n' {
			contentStart++
		}

		// Find nearest matching end or stop marker after contentStart
		rel := data[contentStart:]
		iEndRel := strings.Index(rel, end)
		iStopRel := strings.Index(rel, stop)

		// Choose the first occurring positive index among end/stop
		var contentEnd int
		switch {
		case iEndRel == -1 && iStopRel == -1:
			return nil, errors.New("start marker without matching end/stop marker for project: " + project)
		case iEndRel == -1:
			contentEnd = iStopRel
		case iStopRel == -1:
			contentEnd = iEndRel
		default:
			if iEndRel < iStopRel {
				contentEnd = iEndRel
			} else {
				contentEnd = iStopRel
			}
		}

		content := rel[:contentEnd]
		results = append(results, content)

		// Advance searchFrom to just after the marker we matched
		searchFrom = contentStart + contentEnd
		if strings.HasPrefix(data[searchFrom:], end) {
			searchFrom += len(end)
		} else if strings.HasPrefix(data[searchFrom:], stop) {
			searchFrom += len(stop)
		}
		// If there is a trailing newline after the end/stop marker, skip exactly one
		if searchFrom < len(data) && data[searchFrom] == '\n' {
			searchFrom++
		}
	}

	return results, nil
}
