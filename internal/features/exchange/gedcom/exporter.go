package gedcom

import (
	"bufio"
	"bytes"
	"fmt"
	"genealogy-tree/internal/core/domain"
	"io"
	"sort"
	"strings"
	"time"
)

type familyGroup struct {
	key      string
	partners []string
	parents  []string
	children map[string]struct{}
	metadata map[string]any
}

type familyReferences struct {
	spouse []string
	child  []string
}

func Export(persons []domain.Person, relationships []domain.Relationship, output io.Writer) error {
	personByID := make(map[string]domain.Person, len(persons))
	for _, person := range persons {
		personByID[person.ID] = person
	}

	groups := make(map[string]*familyGroup)
	parentSets := make(map[string][]string)
	parentMetadata := make(map[string]map[string]any)
	for _, relationship := range relationships {
		switch relationship.Type {
		case domain.RelationshipSpouse:
			partners := []string{relationship.Person1ID, relationship.Person2ID}
			sort.Strings(partners)
			key := "spouse:" + strings.Join(partners, ":")
			groups[key] = &familyGroup{key: key, partners: partners, children: map[string]struct{}{}, metadata: relationship.Metadata}
		case domain.RelationshipParentChild:
			parent := relationship.Person1ID
			child := relationship.Person2ID
			if relationship.Direction != nil && *relationship.Direction == domain.DirectionChild {
				parent, child = child, parent
			}
			parentSets[child] = append(parentSets[child], parent)
			parentMetadata[child+"\x00"+parent] = relationship.Metadata
		}
	}
	for child, parents := range parentSets {
		parents = uniqueSorted(parents)
		if len(parents) > 2 {
			for _, parent := range parents {
				key := "parent:" + parent + ":" + child
				group := groups[key]
				if group == nil {
					group = &familyGroup{key: key, parents: []string{parent}, children: map[string]struct{}{}, metadata: parentMetadata[child+"\x00"+parent]}
					groups[key] = group
				}
				group.children[child] = struct{}{}
			}
			continue
		}
		group := matchingSpouseGroup(groups, parents)
		if group == nil {
			key := "parents:" + strings.Join(parents, ":")
			group = groups[key]
			if group == nil {
				group = &familyGroup{key: key, parents: parents, children: map[string]struct{}{}, metadata: parentMetadata[child+"\x00"+parents[0]]}
				groups[key] = group
			}
		}
		group.children[child] = struct{}{}
	}

	orderedPersons := append([]domain.Person(nil), persons...)
	sort.Slice(orderedPersons, func(i, j int) bool { return orderedPersons[i].ID < orderedPersons[j].ID })
	personXrefs := make(map[string]string, len(orderedPersons))
	for index, person := range orderedPersons {
		personXrefs[person.ID] = fmt.Sprintf("@I%d@", index+1)
	}

	orderedGroups := make([]*familyGroup, 0, len(groups))
	for _, group := range groups {
		orderedGroups = append(orderedGroups, group)
	}
	sort.Slice(orderedGroups, func(i, j int) bool { return orderedGroups[i].key < orderedGroups[j].key })
	familyRefs := make(map[string]*familyReferences)
	for index, group := range orderedGroups {
		xref := fmt.Sprintf("@F%d@", index+1)
		members := append(append([]string{}, group.partners...), group.parents...)
		for _, partner := range members {
			if familyRefs[partner] == nil {
				familyRefs[partner] = &familyReferences{}
			}
			familyRefs[partner].spouse = append(familyRefs[partner].spouse, xref)
		}
		for child := range group.children {
			if familyRefs[child] == nil {
				familyRefs[child] = &familyReferences{}
			}
			familyRefs[child].child = append(familyRefs[child].child, xref)
		}
	}

	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	writeLine(writer, "0 HEAD")
	writeLine(writer, "1 SOUR Roots")
	writeLine(writer, "1 CHAR UTF-8")
	writeLine(writer, "1 GEDC")
	writeLine(writer, "2 VERS 5.5.1")
	for _, person := range orderedPersons {
		writePerson(writer, person, personXrefs[person.ID], familyRefs[person.ID])
	}
	for index, group := range orderedGroups {
		writeFamily(writer, group, fmt.Sprintf("@F%d@", index+1), personXrefs, personByID)
	}
	writeLine(writer, "0 TRLR")
	if err := writer.Flush(); err != nil {
		return err
	}
	_, err := output.Write(buffer.Bytes())
	return err
}

func matchingSpouseGroup(groups map[string]*familyGroup, parents []string) *familyGroup {
	if len(parents) == 0 {
		return nil
	}
	var match *familyGroup
	for _, group := range groups {
		if len(group.partners) != 2 {
			continue
		}
		matches := 0
		for _, parent := range parents {
			if contains(group.partners, parent) {
				matches++
			}
		}
		if matches != len(parents) {
			continue
		}
		if match == nil || group.key < match.key {
			match = group
		}
	}
	if len(parents) == 1 {
		matches := 0
		for _, group := range groups {
			if len(group.partners) == 2 && contains(group.partners, parents[0]) {
				matches++
			}
		}
		if matches != 1 {
			return nil
		}
	}
	return match
}

func writePerson(writer *bufio.Writer, person domain.Person, xref string, families *familyReferences) {
	writeLine(writer, "0 "+xref+" INDI")
	firstName := strings.TrimSpace(person.FirstName)
	if person.Patronymic != nil && strings.TrimSpace(*person.Patronymic) != "" {
		firstName += " " + strings.TrimSpace(*person.Patronymic)
	}
	writeLine(writer, "1 NAME "+sanitize(firstName)+" /"+sanitize(person.LastName)+"/")
	if person.Gender != nil {
		gender := map[domain.Gender]string{domain.GenderMale: "M", domain.GenderFemale: "F", domain.GenderOther: "U"}[*person.Gender]
		if gender != "" {
			writeLine(writer, "1 SEX "+gender)
		}
	}
	if person.BirthDate != nil {
		writeLine(writer, "1 BIRT")
		writeLine(writer, "2 DATE "+formatDate(person.BirthDate))
	}
	if person.DeathDate != nil {
		writeLine(writer, "1 DEAT")
		writeLine(writer, "2 DATE "+formatDate(person.DeathDate))
	}
	if families != nil {
		for _, family := range families.spouse {
			writeLine(writer, "1 FAMS "+family)
		}
		for _, family := range families.child {
			writeLine(writer, "1 FAMC "+family)
		}
	}
	if note, ok := metadataString(person.Metadata, "comment"); ok {
		writeNote(writer, note)
	}
	for _, note := range metadataNotes(person.Metadata) {
		writeNote(writer, note)
	}
}

func writeFamily(writer *bufio.Writer, group *familyGroup, xref string, personXrefs map[string]string, personByID map[string]domain.Person) {
	writeLine(writer, "0 "+xref+" FAM")
	partners := append([]string(nil), group.partners...)
	if len(partners) == 2 {
		writePartners(writer, partners, personXrefs, personByID)
	} else {
		writeParents(writer, group.parents, personXrefs, personByID)
	}
	children := make([]string, 0, len(group.children))
	for child := range group.children {
		children = append(children, child)
	}
	sort.Strings(children)
	for _, child := range children {
		if xref := personXrefs[child]; xref != "" {
			writeLine(writer, "1 CHIL "+xref)
		}
	}
	if gedcom, ok := group.metadata["gedcom"].(map[string]any); ok {
		if date, ok := gedcom["marriage_date"].(string); ok && date != "" {
			writeLine(writer, "1 MARR")
			writeLine(writer, "2 DATE "+formatStoredDate(date))
			if place, ok := gedcom["marriage_place"].(string); ok && place != "" {
				writeLine(writer, "2 PLAC "+sanitize(place))
			}
		}
		for _, note := range metadataNotes(gedcom) {
			writeNote(writer, note)
		}
	}
}

func writePartners(writer *bufio.Writer, partners []string, xrefs map[string]string, people map[string]domain.Person) {
	tags := []string{"HUSB", "WIFE"}
	if people[partners[0]].Gender != nil && *people[partners[0]].Gender == domain.GenderFemale {
		tags = []string{"WIFE", "HUSB"}
	}
	for index, partner := range partners {
		writeLine(writer, "1 "+tags[index]+" "+xrefs[partner])
	}
}

func writeParents(writer *bufio.Writer, parents []string, xrefs map[string]string, people map[string]domain.Person) {
	tags := []string{"HUSB", "WIFE"}
	if len(parents) > 0 {
		if person, ok := people[parents[0]]; ok && person.Gender != nil && *person.Gender == domain.GenderFemale {
			tags = []string{"WIFE", "HUSB"}
		}
	}
	for index, parent := range parents {
		if index < len(tags) {
			writeLine(writer, "1 "+tags[index]+" "+xrefs[parent])
		}
	}
}

func writeNote(writer *bufio.Writer, note string) {
	parts := strings.Split(strings.ReplaceAll(note, "\r\n", "\n"), "\n")
	if len(parts) == 0 {
		return
	}
	writeLine(writer, "1 NOTE "+sanitize(parts[0]))
	for _, part := range parts[1:] {
		writeLine(writer, "2 CONT "+sanitize(part))
	}
}

func writeLine(writer *bufio.Writer, line string) {
	_, _ = writer.WriteString(line + "\n")
}

func formatDate(value *time.Time) string {
	return strings.ToUpper(value.Format("2 Jan 2006"))
}

func formatStoredDate(value string) string {
	if len(value) == len("2006-01-02") {
		return value[8:10] + " " + map[string]string{"01": "JAN", "02": "FEB", "03": "MAR", "04": "APR", "05": "MAY", "06": "JUN", "07": "JUL", "08": "AUG", "09": "SEP", "10": "OCT", "11": "NOV", "12": "DEC"}[value[5:7]] + " " + value[:4]
	}
	return sanitize(value)
}

func metadataString(metadata map[string]any, key string) (string, bool) {
	value, ok := metadata[key].(string)
	return value, ok
}

func metadataNotes(metadata map[string]any) []string {
	value, ok := metadata["notes"]
	if !ok {
		if nested, nestedOK := metadata["gedcom"].(map[string]any); nestedOK {
			value, ok = nested["notes"]
		}
	}
	if !ok {
		return nil
	}
	switch notes := value.(type) {
	case []string:
		return notes
	case []any:
		result := make([]string, 0, len(notes))
		for _, note := range notes {
			if value, ok := note.(string); ok {
				result = append(result, value)
			}
		}
		return result
	default:
		return nil
	}
}

func sanitize(value string) string {
	return strings.NewReplacer("\r", " ", "\n", " ").Replace(strings.TrimSpace(value))
}

func uniqueSorted(values []string) []string {
	set := map[string]struct{}{}
	for _, value := range values {
		set[value] = struct{}{}
	}
	result := make([]string, 0, len(set))
	for value := range set {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
