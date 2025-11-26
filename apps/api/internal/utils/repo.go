package utils

import (
	"github.com/huandu/go-sqlbuilder"
)

func GetNotNilAssignments(
	params any,
	ub *sqlbuilder.UpdateBuilder,
) []string {
	fields := ExtractNonNilFieldsByTag(params, "db")
	assignments := make([]string, 0, len(fields))
	for dbTag, value := range fields {
		assignments = append(assignments, ub.Assign(dbTag, value))
	}
	return assignments
}
