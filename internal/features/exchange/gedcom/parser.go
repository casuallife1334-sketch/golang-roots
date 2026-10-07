package gedcom

import (
	"bufio"
	"fmt"
	"genealogy-tree/internal/core/domain"
	"io"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type recordLine struct {
	level int
	xref  string
	tag   string
	value string
	line  int
}

type individual struct {
	xref       string
	name       string
	givenName  string
	nameLine   int
	sex        string
	birthDate  *time.Time
	deathDate  *time.Time
	notes      []string
	activePart string
}

type family struct {
	xref          string
	husband       string
	wife          string
	children      []string
	marriageDate  *time.Time
	marriagePlace string
	notes         []string
	activePart    string
}

// Parse reads the supported GEDCOM 5.5.1 records without writing anything to
// the database. Unsupported tags become warnings so a valid file can still be
// previewed and imported.
func Parse(input io.Reader) (Snapshot, error) {
	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, 1024), 4*1024*1024)

	var individuals []individual
	var families []family
	var currentIndividual *individual
	var currentFamily *family
	var headerSeen bool
	var activeTopLevel string
	var warnings []Issue
	var errors []Issue
	lineNumber := 0

	flush := func() {
		if currentIndividual != nil {
			individuals = append(individuals, *currentIndividual)
			currentIndividual = nil
		}
		if currentFamily != nil {
			families = append(families, *currentFamily)
			currentFamily = nil
		}
	}

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSuffix(scanner.Text(), "\r")
		if lineNumber == 1 {
			line = strings.TrimPrefix(line, "\ufeff")
		}
		if !utf8.ValidString(line) {
			errors = append(errors, Issue{Line: lineNumber, Message: "GEDCOM is not valid UTF-8"})
			continue
		}
		parsed, err := parseLine(line, lineNumber)
		if err != nil {
			errors = append(errors, Issue{Line: lineNumber, Message: err.Error()})
			continue
		}
		parsed.line = lineNumber

		if parsed.level == 0 {
			flush()
			activeTopLevel = parsed.tag
			switch parsed.tag {
			case "HEAD":
				headerSeen = true
			case "SUBM":
				// Submitter records are metadata and do not affect the import.
			case "INDI":
				if parsed.xref == "" {
					errors = append(errors, Issue{Line: lineNumber, Message: "INDI record has no xref"})
					continue
				}
				currentIndividual = &individual{xref: parsed.xref}
			case "FAM":
				if parsed.xref == "" {
					errors = append(errors, Issue{Line: lineNumber, Message: "FAM record has no xref"})
					continue
				}
				currentFamily = &family{xref: parsed.xref}
			case "TRLR":
			default:
				warnings = append(warnings, Issue{Line: lineNumber, Message: "unsupported top-level record " + parsed.tag})
			}
			continue
		}

		if currentIndividual != nil {
			parseIndividualLine(currentIndividual, parsed, &warnings)
			continue
		}
		if currentFamily != nil {
			parseFamilyLine(currentFamily, parsed, &warnings)
			continue
		}
		if activeTopLevel == "HEAD" || activeTopLevel == "SUBM" {
			continue
		}
		warnings = append(warnings, Issue{Line: lineNumber, Message: "record field has no active record"})
	}
	if err := scanner.Err(); err != nil {
		return Snapshot{}, fmt.Errorf("read GEDCOM: %w", err)
	}
	flush()

	if !headerSeen {
		warnings = append(warnings, Issue{Message: "GEDCOM HEAD record is missing"})
	}
	return buildSnapshot(individuals, families, warnings, errors), nil
}

func parseLine(line string, lineNumber int) (recordLine, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return recordLine{}, fmt.Errorf("line %d is not a valid GEDCOM line", lineNumber)
	}
	level, err := strconv.Atoi(fields[0])
	if err != nil || level < 0 {
		return recordLine{}, fmt.Errorf("invalid level")
	}
	if len(fields) >= 3 && strings.HasPrefix(fields[1], "@") && strings.HasSuffix(fields[1], "@") {
		return recordLine{level: level, xref: strings.Trim(fields[1], "@"), tag: fields[2], value: strings.TrimSpace(strings.TrimPrefix(line, fields[0]+" "+fields[1]+" "+fields[2]))}, nil
	}
	return recordLine{level: level, tag: fields[1], value: strings.TrimSpace(strings.TrimPrefix(line, fields[0]+" "+fields[1]))}, nil
}

func parseIndividualLine(person *individual, line recordLine, warnings *[]Issue) {
	if line.level > 2 && (person.activePart == "OBJE" || person.activePart == "NOTE" || person.activePart == "RESI" || person.activePart == "ADOP") {
		if person.activePart == "NOTE" && (line.tag == "CONC" || line.tag == "CONT") && line.value != "" {
			person.notes = append(person.notes, line.value)
		}
		return
	}
	if line.level != 1 && line.level != 2 {
		*warnings = append(*warnings, Issue{Line: line.line, Message: "unsupported INDI nesting level"})
		return
	}
	if line.level == 1 {
		person.activePart = line.tag
		switch line.tag {
		case "NAME":
			person.name = line.value
			person.nameLine = line.line
		case "SEX":
			person.sex = strings.ToUpper(line.value)
		case "BIRT", "DEAT", "NOTE", "FAMS", "FAMC", "RESI", "OBJE", "OCCU", "ADOP", "_UID":
			if line.tag == "NOTE" && line.value != "" {
				person.notes = append(person.notes, line.value)
			}
		default:
			*warnings = append(*warnings, Issue{Line: line.line, Message: "unsupported INDI tag " + line.tag})
		}
		return
	}
	switch {
	case line.tag == "DATE" && person.activePart == "BIRT":
		person.birthDate = parseDate(line.value, line.line, warnings)
	case line.tag == "DATE" && person.activePart == "DEAT":
		person.deathDate = parseDate(line.value, line.line, warnings)
	case line.tag == "CONC" || line.tag == "CONT":
		if person.activePart == "NOTE" && line.value != "" {
			person.notes = append(person.notes, line.value)
		}
	case person.activePart == "OBJE":
		// GEDCOM attachments refer to files outside the import payload.
	case person.activePart == "RESI" || person.activePart == "BIRT" || person.activePart == "DEAT":
		if line.tag == "PLAC" {
			return
		}
	case line.tag == "GIVN":
		person.givenName = line.value
	case line.tag == "SURN" || line.tag == "PEDI" || (line.tag == "FAMC" && person.activePart == "ADOP"):
		// These values duplicate NAME or describe a family link already imported.
	default:
		*warnings = append(*warnings, Issue{Line: line.line, Message: "unsupported INDI field " + line.tag})
	}
}

func parseFamilyLine(record *family, line recordLine, warnings *[]Issue) {
	if line.level > 2 && record.activePart == "NOTE" && (line.tag == "CONC" || line.tag == "CONT") {
		if line.value != "" {
			record.notes = append(record.notes, line.value)
		}
		return
	}
	if line.level != 1 && line.level != 2 {
		*warnings = append(*warnings, Issue{Message: "unsupported FAM nesting level"})
		return
	}
	if line.level == 1 {
		record.activePart = line.tag
		switch line.tag {
		case "HUSB":
			record.husband = strings.Trim(line.value, "@")
		case "WIFE":
			record.wife = strings.Trim(line.value, "@")
		case "CHIL":
			record.children = append(record.children, strings.Trim(line.value, "@"))
		case "MARR", "NOTE":
			if line.tag == "NOTE" && line.value != "" {
				record.notes = append(record.notes, line.value)
			}
		default:
			*warnings = append(*warnings, Issue{Line: line.line, Message: "unsupported FAM tag " + line.tag})
		}
		return
	}
	if line.tag == "NOTE" && record.activePart == "MARR" {
		record.notes = append(record.notes, line.value)
		record.activePart = "NOTE"
		return
	}
	switch {
	case line.tag == "DATE" && record.activePart == "MARR":
		record.marriageDate = parseDate(line.value, line.line, warnings)
	case line.tag == "PLAC" && record.activePart == "MARR":
		record.marriagePlace = line.value
	case line.tag == "CONC" || line.tag == "CONT":
		if record.activePart == "NOTE" && line.value != "" {
			record.notes = append(record.notes, line.value)
		}
	default:
		*warnings = append(*warnings, Issue{Line: line.line, Message: "unsupported FAM field " + line.tag})
	}
}

func parseDate(value string, line int, warnings *[]Issue) *time.Time {
	date, err := parseGEDCOMDate(value)
	if err != nil {
		*warnings = append(*warnings, Issue{Line: line, Message: fmt.Sprintf("unsupported date %q: %s", value, err)})
		return nil
	}
	return date
}

func parseGEDCOMDate(value string) (*time.Time, error) {
	fields := strings.Fields(strings.TrimSpace(value))
	if len(fields) == 0 {
		return nil, fmt.Errorf("date is empty")
	}
	if len(fields) > 1 {
		switch strings.ToUpper(fields[0]) {
		case "ABT", "BEF", "AFT", "CAL", "EST":
			return nil, fmt.Errorf("date qualifier %q is not representable by a DATE column", fields[0])
		}
	}
	month := 1
	day := 1
	yearField := fields[len(fields)-1]
	if len(fields) == 3 {
		parsedDay, err := strconv.Atoi(fields[0])
		if err != nil || parsedDay < 1 || parsedDay > 31 {
			return nil, fmt.Errorf("invalid day")
		}
		day = parsedDay
		month = gedcomMonth(fields[1])
		if month == 0 {
			return nil, fmt.Errorf("invalid month")
		}
	} else if len(fields) == 2 {
		month = gedcomMonth(fields[0])
		if month == 0 {
			return nil, fmt.Errorf("invalid month")
		}
	} else if len(fields) != 1 {
		return nil, fmt.Errorf("unsupported date format")
	}
	year, err := strconv.Atoi(yearField)
	if err != nil || year < 1 || year > 9999 {
		return nil, fmt.Errorf("invalid year")
	}
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if date.Day() != day || int(date.Month()) != month {
		return nil, fmt.Errorf("invalid calendar date")
	}
	return &date, nil
}

func gedcomMonth(value string) int {
	months := map[string]int{"JAN": 1, "FEB": 2, "MAR": 3, "APR": 4, "MAY": 5, "JUN": 6, "JUL": 7, "AUG": 8, "SEP": 9, "OCT": 10, "NOV": 11, "DEC": 12}
	return months[strings.ToUpper(value)]
}

func buildSnapshot(individuals []individual, families []family, warnings, errors []Issue) Snapshot {
	snapshot := Snapshot{Warnings: warnings, Errors: errors}
	seenPersons := make(map[string]struct{}, len(individuals))
	for _, person := range individuals {
		if _, exists := seenPersons[person.xref]; exists {
			snapshot.Errors = append(snapshot.Errors, Issue{Message: "duplicate INDI xref " + person.xref})
			continue
		}
		seenPersons[person.xref] = struct{}{}
		firstName, patronymic, lastName, warning := splitName(person.name, person.givenName)
		if warning != "" {
			snapshot.Warnings = append(snapshot.Warnings, Issue{
				Line:    person.nameLine,
				Message: fmt.Sprintf("INDI %s: %s", person.xref, warning),
			})
		}
		input := domain.CreatePersonInput{
			FirstName:  firstName,
			Patronymic: patronymic,
			LastName:   lastName,
			BirthDate:  person.birthDate,
			DeathDate:  person.deathDate,
			Gender:     parseGender(person.sex, &snapshot.Warnings),
			Metadata:   map[string]any{},
		}
		if len(person.notes) > 0 {
			input.Metadata["comment"] = strings.Join(person.notes, "\n")
		}
		snapshot.Persons = append(snapshot.Persons, PersonRecord{ExternalID: person.xref, Input: input})
	}

	personIDs := make(map[string]struct{}, len(snapshot.Persons))
	for _, person := range snapshot.Persons {
		personIDs[person.ExternalID] = struct{}{}
	}
	seenRelationships := make(map[string]struct{})
	for _, record := range families {
		metadata := familyMetadata(record)
		if record.husband != "" && record.wife != "" {
			addRelationship(&snapshot, seenRelationships, personIDs, record.husband, record.wife, domain.RelationshipSpouse, nil, metadata, record.xref)
		}
		parents := []string{record.husband, record.wife}
		for _, parent := range parents {
			if parent == "" {
				continue
			}
			for _, child := range record.children {
				addRelationship(&snapshot, seenRelationships, personIDs, parent, child, domain.RelationshipParentChild, direction(domain.DirectionParent), metadata, record.xref)
			}
		}
	}
	if len(snapshot.Persons) == 0 {
		snapshot.Errors = append(snapshot.Errors, Issue{Message: "GEDCOM contains no importable INDI records"})
	}
	return snapshot
}

func splitName(value, givenName string) (string, *string, string, string) {
	value = strings.TrimSpace(value)
	start := strings.Index(value, "/")
	if start >= 0 {
		end := strings.Index(value[start+1:], "/")
		if end >= 0 {
			end += start + 1
			firstName := strings.TrimSpace(value[:start])
			lastName := strings.TrimSpace(value[start+1 : end])
			if strings.HasPrefix(lastName, "(") && strings.HasSuffix(lastName, ")") {
				lastName = strings.TrimSpace(lastName[1 : len(lastName)-1])
			}
			warning := ""
			if firstName == "" {
				firstName = "Unknown"
				warning = "NAME has no given name"
			}
			if lastName == "" {
				lastName = "Unknown"
				if warning != "" {
					warning = "NAME has no given name or surname"
				} else {
					warning = "NAME has no surname"
				}
			}
			firstName, patronymic := splitGivenName(firstName, givenName)
			return firstName, patronymic, lastName, warning
		}
	}

	fields := strings.Fields(value)
	if len(fields) >= 2 {
		firstName, patronymic := splitGivenName(strings.Join(fields[:len(fields)-1], " "), givenName)
		return firstName, patronymic, fields[len(fields)-1], "NAME does not use the standard surname slash format"
	}
	if len(fields) == 1 {
		firstName, patronymic := splitGivenName(fields[0], givenName)
		return firstName, patronymic, "Unknown", "NAME has no surname"
	}
	return "Unknown", nil, "Unknown", "NAME is empty"
}

func splitGivenName(name, givenName string) (string, *string) {
	if strings.TrimSpace(givenName) != "" {
		name = strings.TrimSpace(givenName)
	}
	parts := strings.Fields(name)
	if len(parts) < 2 || !looksLikePatronymic(parts[len(parts)-1]) {
		return name, nil
	}
	patronymic := parts[len(parts)-1]
	firstName := strings.Join(parts[:len(parts)-1], " ")
	return firstName, &patronymic
}

func looksLikePatronymic(value string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	for _, suffix := range []string{"ович", "евич", "ич", "овна", "евна", "ична", "инична", "ovich", "evich", "ovna", "evna", "ichna"} {
		if strings.HasSuffix(value, suffix) {
			return true
		}
	}
	return false
}

func parseGender(value string, warnings *[]Issue) *domain.Gender {
	var gender domain.Gender
	switch strings.ToUpper(value) {
	case "M":
		gender = domain.GenderMale
	case "F":
		gender = domain.GenderFemale
	case "", "U":
		return nil
	default:
		*warnings = append(*warnings, Issue{Message: "unsupported SEX value " + value})
		gender = domain.GenderOther
	}
	return &gender
}

func familyMetadata(record family) map[string]any {
	metadata := map[string]any{"family_xref": record.xref}
	gedcom := map[string]any{}
	if record.marriageDate != nil {
		gedcom["marriage_date"] = dateValue(record.marriageDate)
	}
	if record.marriagePlace != "" {
		gedcom["marriage_place"] = record.marriagePlace
	}
	if len(record.notes) > 0 {
		gedcom["notes"] = record.notes
	}
	if len(gedcom) > 0 {
		metadata["gedcom"] = gedcom
	}
	return metadata
}

func addRelationship(snapshot *Snapshot, seen map[string]struct{}, personIDs map[string]struct{}, person1, person2 string, relationshipType domain.RelationshipType, relationshipDirection *domain.RelationshipDirection, metadata map[string]any, familyXref string) {
	if _, ok := personIDs[person1]; !ok {
		snapshot.Warnings = append(snapshot.Warnings, Issue{Message: fmt.Sprintf("FAM %s references unknown person %s; relationship skipped", familyXref, person1)})
		return
	}
	if _, ok := personIDs[person2]; !ok {
		snapshot.Warnings = append(snapshot.Warnings, Issue{Message: fmt.Sprintf("FAM %s references unknown person %s; relationship skipped", familyXref, person2)})
		return
	}
	if person1 == person2 {
		snapshot.Warnings = append(snapshot.Warnings, Issue{Message: fmt.Sprintf("FAM %s contains a self relationship; relationship skipped", familyXref)})
		return
	}
	if relationshipType == domain.RelationshipSpouse && person1 > person2 {
		person1, person2 = person2, person1
	}
	key := string(relationshipType) + ":" + person1 + ":" + person2
	if relationshipDirection != nil {
		key += ":" + string(*relationshipDirection)
	}
	if _, ok := seen[key]; ok {
		return
	}
	seen[key] = struct{}{}
	copyMetadata := map[string]any{}
	for key, value := range metadata {
		copyMetadata[key] = value
	}
	snapshot.Relationships = append(snapshot.Relationships, RelationshipRecord{
		Person1ExternalID: person1,
		Person2ExternalID: person2,
		Input: domain.CreateRelationshipInput{
			Person1ID: person1,
			Person2ID: person2,
			Type:      relationshipType,
			Direction: relationshipDirection,
			Metadata:  copyMetadata,
		},
	})
}

func direction(value domain.RelationshipDirection) *domain.RelationshipDirection {
	return &value
}
