package db

import (
	"fmt"
	"strings"
)

var (
	// The fmt / string sub here only inserts the table name (constant strings defined here in the repo).
	// This is done so that we can easily change the table names in the future.
	//
	// All actual data is interpolated using the sql driver(s) .. so it's properly escaped.

	sqlFactionChildren = strings.ReplaceAll(`WITH RECURSIVE recurse AS (
	SELECT faction_table.id
	    FROM faction_table
	    WHERE faction_table.parent_faction_id = $1
	    AND faction_table.parent_faction_relation IN (%s)
		UNION
	    SELECT faction_table.id
	    FROM faction_table
	    JOIN recurse ON faction_table.parent_faction_id = recurse.id
	    WHERE faction_table.parent_faction_relation IN (%s)
	)
	SELECT * FROM faction_table JOIN recurse on faction_table.id = recurse.id;`, "faction_table", tableFactions)

	sqlFactionParents = strings.ReplaceAll(`WITH RECURSIVE recurse AS (
	    SELECT faction_table.parent_faction_id
	    FROM faction_table
	    WHERE faction_table.id = $1
	    AND faction_table.parent_faction_relation IN (%s)
		UNION
	    SELECT faction_table.parent_faction_id
	    FROM faction_table
	    JOIN recurse ON faction_table.id = recurse.parent_faction_id
	    AND faction_table.parent_faction_relation IN (%s)
	)
	SELECT * FROM faction_table JOIN recurse on faction_table.id = recurse.parent_faction_id;`, "faction_table", tableFactions)

	sqlInsertFamilies = fmt.Sprintf(`INSERT INTO %s (
	    id, race, culture, area_id, faction_id, 
	    ethos_altruism, ethos_ambition, ethos_tradition, ethos_pacifism, ethos_piety, ethos_caution,
	    is_child_bearing, max_child_bearing_tick,  pregnancy_end,
	    male_id, female_id,
	    ma_grandma_id, ma_grandpa_id, pa_grandma_id, pa_grandpa_id,
	    number_of_children, random
	) VALUES (
	    :id, :race, :culture, :area_id, :faction_id, 
	    :ethos_altruism, :ethos_ambition, :ethos_tradition, :ethos_pacifism, :ethos_piety, :ethos_caution,
	    :is_child_bearing, :max_child_bearing_tick, :pregnancy_end,
	    :male_id, :female_id,
	    :ma_grandma_id, :ma_grandpa_id, :pa_grandma_id, :pa_grandpa_id,
	    :number_of_children, :random
	) ON CONFLICT (id) DO UPDATE SET
	    race=EXCLUDED.race,
	    culture=EXCLUDED.culture,
	    ethos_altruism=EXCLUDED.ethos_altruism,
	    ethos_ambition=EXCLUDED.ethos_ambition,
	    ethos_tradition=EXCLUDED.ethos_tradition,
	    ethos_pacifism=EXCLUDED.ethos_pacifism,
	    ethos_piety=EXCLUDED.ethos_piety,
	    ethos_caution=EXCLUDED.ethos_caution,
	    area_id=EXCLUDED.area_id,
	    faction_id=EXCLUDED.faction_id,
	    is_child_bearing=EXCLUDED.is_child_bearing,
	    pregnancy_end=EXCLUDED.pregnancy_end,
	    number_of_children=EXCLUDED.number_of_children
	;`, tableFamilies)

	sqlInsertFactions = fmt.Sprintf(`INSERT INTO %s (
            id, name, home_area_id, hq_plot_id,
	    ethos_altruism, ethos_ambition, ethos_tradition, ethos_pacifism, ethos_piety, ethos_caution,
            action_frequency_ticks,
            leadership, structure, wealth, cohesion, corruption,
	    is_covert,
	    government_id, is_government,
	    religion_id, is_religion,
	    is_member_by_birth,
	    espionage_offense, espionage_defense,
	    military_offense, military_defense,
            parent_faction_id, parent_faction_relation,
	    members, plots, areas
	) VALUES (
	    :id, :name,
	    :home_area_id, :hq_plot_id,
	    :ethos_altruism, :ethos_ambition, :ethos_tradition, :ethos_pacifism, :ethos_piety, :ethos_caution,
	    :action_frequency_ticks,
	    :leadership, :structure, :wealth, :cohesion, :corruption,
	    :is_covert,
	    :government_id, :is_government,
	    :religion_id, :is_religion,
	    :is_member_by_birth,
	    :espionage_offense, :espionage_defense,
	    :military_offense, :military_defense,
	    :parent_faction_id, :parent_faction_relation,
	    :members, :plots, :areas
	) ON CONFLICT (id) DO UPDATE SET
	    name=EXCLUDED.name,
	    home_area_id=EXCLUDED.home_area_id,
	    hq_plot_id=EXCLUDED.hq_plot_id,
	    ethos_altruism=EXCLUDED.ethos_altruism,
	    ethos_ambition=EXCLUDED.ethos_ambition,
	    ethos_tradition=EXCLUDED.ethos_tradition,
	    ethos_pacifism=EXCLUDED.ethos_pacifism,
	    ethos_piety=EXCLUDED.ethos_piety,
	    ethos_caution=EXCLUDED.ethos_caution,
	    action_frequency_ticks=EXCLUDED.action_frequency_ticks,
	    leadership=EXCLUDED.leadership,
	    structure=EXCLUDED.structure,
	    wealth=EXCLUDED.wealth,
	    cohesion=EXCLUDED.cohesion,
	    corruption=EXCLUDED.corruption,
	    is_covert=EXCLUDED.is_covert,
	    government_id=EXCLUDED.government_id,
	    is_government=EXCLUDED.is_government,
	    religion_id=EXCLUDED.religion_id,
	    is_religion=EXCLUDED.is_religion,
	    is_member_by_birth=EXCLUDED.is_member_by_birth,
	    espionage_offense=EXCLUDED.espionage_offense,
	    espionage_defense=EXCLUDED.espionage_defense,
	    military_offense=EXCLUDED.military_offense,
	    military_defense=EXCLUDED.military_defense,
	    parent_faction_id=EXCLUDED.parent_faction_id,
	    parent_faction_relation=EXCLUDED.parent_faction_relation,
	    members=EXCLUDED.members,
	    plots=EXCLUDED.plots,
	    areas=EXCLUDED.areas
	;`, tableFactions)

	sqlInsertAreas = fmt.Sprintf(`INSERT INTO %s (id, government_id, biome, random)
			VALUES (:id, :government_id, :biome, :random)
		        ON CONFLICT (id) DO UPDATE SET
		            government_id=EXCLUDED.government_id,
			    biome=EXCLUDED.biome
		        ;`,
		tableAreas,
	)

	sqlInsertGovernments = fmt.Sprintf(`INSERT INTO %s (
	    id, tax_rate, tax_frequency
	) VALUES (
	    :id, :tax_rate, :tax_frequency
	) ON CONFLICT (id) DO UPDATE SET 
	    tax_rate=EXCLUDED.tax_rate,
	    tax_frequency=EXCLUDED.tax_frequency
	;`, tableGovernments)

	sqlInsertLaws = fmt.Sprintf(`INSERT INTO %s (
	    source_id, meta_key, meta_val, illegal
	) VALUES (
	    :source_id, :meta_key, :meta_val, :illegal
	);`, tableLaws)

	sqlInsertPlots = fmt.Sprintf(`INSERT INTO %s (
	    id, area_id, faction_id, hidden, value, size, commodity, yield
	) VALUES (
	    :id, :area_id, :faction_id, :hidden, :value, :size, :commodity, :yield
	) ON CONFLICT (id) DO UPDATE SET 
	    faction_id=EXCLUDED.faction_id,
	    hidden=EXCLUDED.hidden,
	    value=EXCLUDED.value,
	    size=EXCLUDED.size,
	    commodity=EXCLUDED.commodity,
	    yield=EXCLUDED.yield
	;`, tablePlots)

	sqlInsertPeople = fmt.Sprintf(`INSERT INTO %s (
	    id,
	    ethos_altruism, ethos_ambition, ethos_tradition, ethos_pacifism, ethos_piety, ethos_caution,
	    first_name, last_name, birth_family_id, race, culture,
	    area_id, job_id,
	    birth_tick, death_tick, is_male, adulthood_tick,
	    preferred_profession, preferred_faction_id,
	    death_meta_reason, death_meta_key, death_meta_val, natural_death_tick,
	    random
	) VALUES (
	    :id,
	    :ethos_altruism, :ethos_ambition, :ethos_tradition, :ethos_pacifism, :ethos_piety, :ethos_caution,
	    :first_name, :last_name, :birth_family_id, :race, :culture,
	    :area_id, :job_id,
	    :birth_tick, :death_tick,
	    :is_male, :adulthood_tick,
	    :preferred_profession, :preferred_faction_id,
	    :death_meta_reason, :death_meta_key, :death_meta_val, :natural_death_tick,
	    :random
	) ON CONFLICT (id) DO UPDATE SET
	    ethos_altruism=EXCLUDED.ethos_altruism,
	    ethos_ambition=EXCLUDED.ethos_ambition,
	    ethos_tradition=EXCLUDED.ethos_tradition,
	    ethos_pacifism=EXCLUDED.ethos_pacifism,
	    ethos_piety=EXCLUDED.ethos_piety,
	    ethos_caution=EXCLUDED.ethos_caution,
	    race=EXCLUDED.race,
	    culture=EXCLUDED.culture,
	    area_id=EXCLUDED.area_id,
	    job_id=EXCLUDED.job_id,
	    preferred_profession=EXCLUDED.preferred_profession,
	    preferred_faction_id=EXCLUDED.preferred_faction_id,
	    death_tick=EXCLUDED.death_tick,
	    death_meta_reason=EXCLUDED.death_meta_reason,
	    death_meta_key=EXCLUDED.death_meta_key,
	    death_meta_val=EXCLUDED.death_meta_val
	;`, tablePeople)

	sqlInsertEvents = fmt.Sprintf(`INSERT INTO %s (
	    id, type, tick, message,
	    subject_meta_key, subject_meta_val,
	    cause_meta_key, cause_meta_val
	) VALUES (	
	    :id, :type, :tick, :message,
	    :subject_meta_key, :subject_meta_val,
	    :cause_meta_key, :cause_meta_val
	);`, tableEvents)

	sqlInsertJobs = fmt.Sprintf(`INSERT INTO %s (
	    id, parent_job_id,
	    source_faction_id, source_area_id,
	    action,
	    target_faction_id, target_area_id, target_meta_key, target_meta_val,
	    people_min, people_max,
	    tick_created, tick_starts, tick_ends,
	    secrecy,
	    is_illegal,
	    state
	) VALUES (
	    :id, :parent_job_id,
	    :source_faction_id, :source_area_id,
	    :action,
	    :target_faction_id, :target_area_id, :target_meta_key, :target_meta_val,
	    :people_min, :people_max,
	    :tick_created, :tick_starts, :tick_ends,
	    :secrecy,
	    :is_illegal,
	    :state
	) ON CONFLICT (id) DO UPDATE SET 
	    state=EXCLUDED.state,
	    secrecy=EXCLUDED.secrecy,
	    is_illegal=EXCLUDED.is_illegal,
	    tick_starts=EXCLUDED.tick_starts,
	    tick_ends=EXCLUDED.tick_ends,
	    people_now=EXCLUDED.people_now,
	    target_area_id=EXCLUDED.target_area_id,
	    target_meta_key=EXCLUDED.target_meta_key,
	    target_meta_val=EXCLUDED.target_meta_val
	;`, tableJobs)
)
