package db

import (
	"fmt"
	"strings"

	"github.com/voidshard/faction/internal/dbutils"
)

func sqlFromPlotFilters(token string, in []*PlotFilter) (string, []interface{}, error) {
	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return "", nil, err
	}

	var (
		ors   []string
		args  []interface{}
		where string
	)

	for _, f := range in {
		ands := []string{}

		if dbutils.IsValidID(f.ID) {
			ands = append(ands, fmt.Sprintf("id = $%d", len(args)+1))
			args = append(args, f.ID)
		}
		if dbutils.IsValidID(f.OwnerFactionID) {
			ands = append(ands, fmt.Sprintf("owner_faction_id = $%d", len(args)+1))
			args = append(args, f.OwnerFactionID)
		}
		if dbutils.IsValidID(f.AreaID) {
			ands = append(ands, fmt.Sprintf("area_id = $%d", len(args)+1))
			args = append(args, f.AreaID)
		}

		if len(ands) > 0 {
			ors = append(ors, fmt.Sprintf("(%s)", strings.Join(ands, " AND ")))
		}
	}

	if len(ors) > 0 {
		where = fmt.Sprintf("WHERE %s", strings.Join(ors, " OR "))
	}

	return fmt.Sprintf(`SELECT
		id, is_head_quarters, area_id, owner_faction_id, size
	    FROM %s
	    %s
	    ORDER BY id ASC LIMIT $%d OFFSET $%d;`,
		tablePlots,
		where,
		len(args)+1,
		len(args)+2,
	), append(args, tk.Limit, tk.Offset), nil
}

func sqlFromPersonFilters(token string, in []*PersonFilter) (string, []interface{}, error) {
	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return "", nil, err
	}

	var (
		ors   []string
		args  []interface{}
		where string
	)

	for _, f := range in {
		ands := []string{}

		if dbutils.IsValidID(f.ID) {
			ands = append(ands, fmt.Sprintf("id = $%d", len(args)+1))
			args = append(args, f.ID)
		}
		if dbutils.IsValidID(f.JobID) {
			ands = append(ands, fmt.Sprintf("job_id = $%d", len(args)+1))
			args = append(args, f.JobID)
		}
		if dbutils.IsValidID(f.AreaID) {
			ands = append(ands, fmt.Sprintf("area_id = $%d", len(args)+1))
			args = append(args, f.AreaID)
		}
		if !f.IncludeDead { // add this by default
			ands = append(ands, fmt.Sprintf("death_tick is null"))
		}
		if f.MinEthos != nil {
			ands = append(ands, fmt.Sprintf("ethos_altruism >= $%d", len(args)+1))
			args = append(args, f.MinEthos.Altruism)
			ands = append(ands, fmt.Sprintf("ethos_ambition >= $%d", len(args)+1))
			args = append(args, f.MinEthos.Ambition)
			ands = append(ands, fmt.Sprintf("ethos_tradition >= $%d", len(args)+1))
			args = append(args, f.MinEthos.Tradition)
			ands = append(ands, fmt.Sprintf("ethos_pacifism >= $%d", len(args)+1))
			args = append(args, f.MinEthos.Pacifism)
			ands = append(ands, fmt.Sprintf("ethos_piety >= $%d", len(args)+1))
			args = append(args, f.MinEthos.Piety)
			ands = append(ands, fmt.Sprintf("ethos_caution >= $%d", len(args)+1))
			args = append(args, f.MinEthos.Caution)
		}
		if f.MaxEthos != nil {
			ands = append(ands, fmt.Sprintf("ethos_altruism <= $%d", len(args)+1))
			args = append(args, f.MaxEthos.Altruism)
			ands = append(ands, fmt.Sprintf("ethos_ambition <= $%d", len(args)+1))
			args = append(args, f.MaxEthos.Ambition)
			ands = append(ands, fmt.Sprintf("ethos_tradition <= $%d", len(args)+1))
			args = append(args, f.MaxEthos.Tradition)
			ands = append(ands, fmt.Sprintf("ethos_pacifism <= $%d", len(args)+1))
			args = append(args, f.MaxEthos.Pacifism)
			ands = append(ands, fmt.Sprintf("ethos_piety <= $%d", len(args)+1))
			args = append(args, f.MaxEthos.Piety)
			ands = append(ands, fmt.Sprintf("ethos_caution <= $%d", len(args)+1))
			args = append(args, f.MaxEthos.Caution)
		}

		if len(ands) > 0 {
			ors = append(ors, fmt.Sprintf("(%s)", strings.Join(ands, " AND ")))
		}
	}

	if len(ors) > 0 {
		where = fmt.Sprintf("WHERE %s", strings.Join(ors, " OR "))
	}

	return fmt.Sprintf(`SELECT
		    id,
		    ethos_altruism,
		    ethos_ambition,
		    ethos_tradition,
		    ethos_pacifism,
		    ethos_piety,
		    ethos_caution,
		    area_id,
		    job_id,
		    birth_tick,
		    death_tick,
		    is_male
		FROM %s
		%s
		ORDER BY id ASC LIMIT $%d OFFSET $%d;`,
		tablePeople,
		where,
		len(args)+1,
		len(args)+2,
	), append(args, tk.Limit, tk.Offset), nil
}

func sqlFromAreaFilters(token string, in []*AreaFilter) (string, []interface{}, error) {
	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return "", nil, err
	}

	var (
		ors   []string
		args  []interface{}
		where string
	)

	for _, f := range in {
		ands := []string{}

		if dbutils.IsValidID(f.ID) {
			ands = append(ands, fmt.Sprintf("id = $%d", len(args)+1))
			args = append(args, f.ID)
		}
		if dbutils.IsValidID(f.GoverningFactionID) {
			ands = append(ands, fmt.Sprintf("governing_faction_id = $%d", len(args)+1))
			args = append(args, f.GoverningFactionID)
		}

		if len(ands) > 0 {
			ors = append(ors, fmt.Sprintf("(%s)", strings.Join(ands, " AND ")))
		}
	}

	if len(ors) > 0 {
		where = fmt.Sprintf("WHERE %s", strings.Join(ors, " OR "))
	}

	return fmt.Sprintf(`SELECT id, governing_faction_id
		FROM %s
		%s 
		ORDER BY id ASC LIMIT $%d OFFSET $%d;`,
		tableAreas,
		where,
		len(args)+1,
		len(args)+2,
	), append(args, tk.Limit, tk.Offset), nil
}
