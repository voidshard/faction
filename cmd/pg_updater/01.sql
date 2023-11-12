\c faction
SET search_path TO faction;

--- Create tables
CREATE TABLE IF NOT EXISTS meta (
	    id VARCHAR(255) PRIMARY KEY,
	    str VARCHAR(255) DEFAULT '',
	    int INTEGER NOT NULL DEFAULT 0
	);
INSERT INTO meta (id, int) VALUES ('tick', 1) ON CONFLICT (id) DO NOTHING;
CREATE TABLE IF NOT EXISTS areas (
	    id VARCHAR(36) PRIMARY KEY,
	    government_id VARCHAR(36),
	    biome VARCHAR(255) DEFAULT '',
	    random INTEGER NOT NULL DEFAULT 0
	);
CREATE TABLE IF NOT EXISTS plots (
	    id VARCHAR(36) PRIMARY KEY,
	    area_id VARCHAR(36) NOT NULL,
	    faction_id VARCHAR(36) DEFAULT '',
	    hidden INTEGER NOT NULL DEFAULT 0,
	    value REAL NOT NULL DEFAULT 0,
	    size INTEGER NOT NULL DEFAULT 1,
	    commodity VARCHAR(255) DEFAULT '',
	    yield INTEGER NOT NULL DEFAULT 0
	);
CREATE TABLE IF NOT EXISTS governments (
	    id VARCHAR(36) PRIMARY KEY,
	    tax_rate REAL NOT NULL DEFAULT 0.10,
	    tax_frequency INTEGER NOT NULL DEFAULT 1
	);
CREATE TABLE IF NOT EXISTS laws (
	    source_id VARCHAR(36) NOT NULL,
	    meta_key VARCHAR(255) NOT NULL,
	    meta_val VARCHAR(255) NOT NULL,
	    illegal BOOLEAN NOT NULL DEFAULT FALSE,
	    UNIQUE(source_id, meta_key, meta_val)
	);
CREATE TABLE IF NOT EXISTS jobs (
	    id VARCHAR(36) PRIMARY KEY,
	    parent_job_id VARCHAR(36) DEFAULT '',
	    source_faction_id VARCHAR(36) NOT NULL,
	    source_area_id VARCHAR(36) NOT NULL,
	    action VARCHAR(255) NOT NULL,
            priority INTEGER NOT NULL DEFAULT 0,
	    conscription BOOLEAN NOT NULL DEFAULT FALSE,
	    target_area_id VARCHAR(36) NOT NULL,
	    target_faction_id VARCHAR(36) DEFAULT '',
	    target_meta_key VARCHAR(20) DEFAULT '',
	    target_meta_val VARCHAR(255) DEFAULT '',
	    people_min INTEGER NOT NULL DEFAULT 1,
	    people_max INTEGER NOT NULL DEFAULT 0,
	    people_now INTEGER NOT NULL DEFAULT 0,
	    tick_created INTEGER NOT NULL DEFAULT 0,
	    tick_starts INTEGER NOT NULL DEFAULT 0,
	    tick_ends INTEGER NOT NULL DEFAULT 1,
	    secrecy INTEGER NOT NULL DEFAULT 0,
	    is_illegal BOOLEAN NOT NULL DEFAULT FALSE,
	    state VARCHAR(20) NOT NULL DEFAULT 'pending'
	);
CREATE TABLE IF NOT EXISTS families (
	    id VARCHAR(36) PRIMARY KEY,
	    race VARCHAR(255) DEFAULT '',
	    culture VARCHAR(255) DEFAULT '',
            ethos_altruism INTEGER NOT NULL DEFAULT 0,
            ethos_ambition INTEGER NOT NULL DEFAULT 0,
            ethos_tradition INTEGER NOT NULL DEFAULT 0,
            ethos_pacifism INTEGER NOT NULL DEFAULT 0,
            ethos_piety INTEGER NOT NULL DEFAULT 0,
            ethos_caution INTEGER NOT NULL DEFAULT 0,
	    area_id VARCHAR(36) NOT NULL,
	    social_class VARCHAR(255) DEFAULT '',
	    faction_id VARCHAR(36) DEFAULT '',
	    is_child_bearing BOOLEAN NOT NULL DEFAULT FALSE,
	    male_id VARCHAR(36) NOT NULL,
	    female_id VARCHAR(36) NOT NULL,
	    max_child_bearing_tick INTEGER NOT NULL DEFAULT 0,
	    pregnancy_end INTEGER NOT NULL DEFAULT 0,
	    ma_grandma_id VARCHAR(36) DEFAULT '',
	    ma_grandpa_id VARCHAR(36) DEFAULT '',
	    pa_grandma_id VARCHAR(36) DEFAULT '',
	    pa_grandpa_id VARCHAR(36) DEFAULT '',
	    number_of_children INTEGER NOT NULL DEFAULT 0,
	    random INTEGER NOT NULL DEFAULT 0,
	    UNIQUE (male_id, female_id)
	);
CREATE TABLE IF NOT EXISTS people (
            id VARCHAR(36) PRIMARY KEY,
	    first_name VARCHAR(255) DEFAULT '',
	    last_name VARCHAR(255) DEFAULT '',
	    birth_family_id VARCHAR(36) DEFAULT '',
	    race VARCHAR(255) DEFAULT '',
	    culture VARCHAR(255) DEFAULT '',
            ethos_altruism INTEGER NOT NULL DEFAULT 0,
            ethos_ambition INTEGER NOT NULL DEFAULT 0,
            ethos_tradition INTEGER NOT NULL DEFAULT 0,
            ethos_pacifism INTEGER NOT NULL DEFAULT 0,
            ethos_piety INTEGER NOT NULL DEFAULT 0,
            ethos_caution INTEGER NOT NULL DEFAULT 0,
	    area_id VARCHAR(36) NOT NULL,
	    job_id VARCHAR(36) DEFAULT '',
	    birth_tick INTEGER NOT NULL DEFAULT 1,
	    death_tick INTEGER NOT NULL DEFAULT 0,
	    is_male BOOLEAN NOT NULL DEFAULT FALSE,
	    adulthood_tick INTEGER NOT NULL DEFAULT 0,
	    preferred_profession VARCHAR(255) DEFAULT '',
	    preferred_faction_id VARCHAR(36) DEFAULT '',
	    death_meta_reason TEXT DEFAULT '',
	    death_meta_key VARCHAR(255) DEFAULT '',
	    death_meta_val VARCHAR(255) DEFAULT '',
	    natural_death_tick INTEGER NOT NULL DEFAULT 0,
	    random INTEGER NOT NULL DEFAULT 0
	);
CREATE TABLE IF NOT EXISTS factions (
            id VARCHAR(36) PRIMARY KEY,
	    name VARCHAR(255) DEFAULT '',
	    home_area_id VARCHAR(36) NOT NULL,
	    hq_plot_id VARCHAR(36) NOT NULL,
            ethos_altruism INTEGER NOT NULL DEFAULT 0,
            ethos_ambition INTEGER NOT NULL DEFAULT 0,
            ethos_tradition INTEGER NOT NULL DEFAULT 0,
            ethos_pacifism INTEGER NOT NULL DEFAULT 0,
            ethos_piety INTEGER NOT NULL DEFAULT 0,
            ethos_caution INTEGER NOT NULL DEFAULT 0,
            action_frequency_ticks INTEGER NOT NULL DEFAULT 1,
            leadership INTEGER NOT NULL DEFAULT 0,
	    structure INTEGER NOT NULL DEFAULT 0,
            wealth INTEGER NOT NULL DEFAULT 0,
            cohesion INTEGER NOT NULL DEFAULT 0,
            corruption INTEGER NOT NULL DEFAULT 0,
            is_covert BOOLEAN NOT NULL DEFAULT FALSE,
            government_id VARCHAR(36) NOT NULL,
            is_government BOOLEAN NOT NULL DEFAULT FALSE,
            religion_id VARCHAR(36) DEFAULT '',
            is_religion BOOLEAN NOT NULL DEFAULT FALSE,
            is_member_by_birth BOOLEAN NOT NULL DEFAULT FALSE,
	    espionage_offense INTEGER NOT NULL DEFAULT 0,
	    espionage_defense INTEGER NOT NULL DEFAULT 0,
	    military_offense INTEGER NOT NULL DEFAULT 0,
	    military_defense INTEGER NOT NULL DEFAULT 0,
            parent_faction_id VARCHAR(36) DEFAULT '',
            parent_faction_relation INTEGER NOT NULL DEFAULT 0,
	    members INTEGER NOT NULL DEFAULT 0,
	    vassals INTEGER NOT NULL DEFAULT 0,
	    plots INTEGER NOT NULL DEFAULT 0,
	    areas INTEGER NOT NULL DEFAULT 0
        );
CREATE TABLE IF NOT EXISTS events (
	    id VARCHAR(36) PRIMARY KEY,
	    type VARCHAR(255) DEFAULT '',
	    tick INTEGER NOT NULL DEFAULT 0,
	    message TEXT DEFAULT '',
	    subject_meta_key VARCHAR(255) DEFAULT '',
	    subject_meta_val VARCHAR(255) DEFAULT '',
	    cause_meta_key VARCHAR(255) DEFAULT '',
	    cause_meta_val VARCHAR(255) DEFAULT ''
	);
CREATE TABLE IF NOT EXISTS tuples_trust_faction_to_faction (
	    subject VARCHAR(255) NOT NULL, 
	    object VARCHAR(255) NOT NULL,
	    value INTEGER NOT NULL DEFAULT 0,
	    UNIQUE(subject, object)
	);
CREATE TABLE IF NOT EXISTS modifiers_trust_faction_to_faction (
	    subject VARCHAR(255) NOT NULL,
	    object VARCHAR(255) NOT NULL,
	    value INTEGER NOT NULL DEFAULT 0,
	    tick_expires INTEGER NOT NULL DEFAULT 0,
	    meta_key VARCHAR(255) DEFAULT '',
	    meta_val VARCHAR(255) DEFAULT '',
	    meta_reason TEXT DEFAULT ''
	);
CREATE TABLE IF NOT EXISTS tuples_import_faction_to_commodity (
	    subject VARCHAR(255) NOT NULL, 
	    object VARCHAR(255) NOT NULL,
	    value INTEGER NOT NULL DEFAULT 0,
	    UNIQUE(subject, object)
	);
CREATE TABLE IF NOT EXISTS tuples_export_faction_to_commodity (
	    subject VARCHAR(255) NOT NULL, 
	    object VARCHAR(255) NOT NULL,
	    value INTEGER NOT NULL DEFAULT 0,
	    UNIQUE(subject, object)
	);
CREATE TABLE IF NOT EXISTS tuples_research_faction_to_topic (
	    subject VARCHAR(255) NOT NULL, 
	    object VARCHAR(255) NOT NULL,
	    value INTEGER NOT NULL DEFAULT 0,
	    UNIQUE(subject, object)
	);
CREATE TABLE IF NOT EXISTS tuples_weight_faction_to_topic (
	    subject VARCHAR(255) NOT NULL, 
	    object VARCHAR(255) NOT NULL,
	    value INTEGER NOT NULL DEFAULT 0,
	    UNIQUE(subject, object)
	);
CREATE TABLE IF NOT EXISTS tuples_affiliation_person_to_faction (
	    subject VARCHAR(255) NOT NULL, 
	    object VARCHAR(255) NOT NULL,
	    value INTEGER NOT NULL DEFAULT 0,
	    UNIQUE(subject, object)
	);
CREATE TABLE IF NOT EXISTS modifiers_affiliation_person_to_faction (
	    subject VARCHAR(255) NOT NULL,
	    object VARCHAR(255) NOT NULL,
	    value INTEGER NOT NULL DEFAULT 0,
	    tick_expires INTEGER NOT NULL DEFAULT 0,
	    meta_key VARCHAR(255) DEFAULT '',
	    meta_val VARCHAR(255) DEFAULT '',
	    meta_reason TEXT DEFAULT ''
	);
CREATE TABLE IF NOT EXISTS tuples_skill_person_to_profession (
	    subject VARCHAR(255) NOT NULL, 
	    object VARCHAR(255) NOT NULL,
	    value INTEGER NOT NULL DEFAULT 0,
	    UNIQUE(subject, object)
	);
CREATE TABLE IF NOT EXISTS tuples_relationship_person_to_person (
	    subject VARCHAR(255) NOT NULL, 
	    object VARCHAR(255) NOT NULL,
	    value INTEGER NOT NULL DEFAULT 0,
	    UNIQUE(subject, object)
	);
CREATE TABLE IF NOT EXISTS tuples_faith_person_to_religion (
	    subject VARCHAR(255) NOT NULL, 
	    object VARCHAR(255) NOT NULL,
	    value INTEGER NOT NULL DEFAULT 0,
	    UNIQUE(subject, object)
	);
CREATE TABLE IF NOT EXISTS tuples_trust_person_to_person (
	    subject VARCHAR(255) NOT NULL, 
	    object VARCHAR(255) NOT NULL,
	    value INTEGER NOT NULL DEFAULT 0,
	    UNIQUE(subject, object)
	);
CREATE TABLE IF NOT EXISTS modifiers_trust_person_to_person (
	    subject VARCHAR(255) NOT NULL,
	    object VARCHAR(255) NOT NULL,
	    value INTEGER NOT NULL DEFAULT 0,
	    tick_expires INTEGER NOT NULL DEFAULT 0,
	    meta_key VARCHAR(255) DEFAULT '',
	    meta_val VARCHAR(255) DEFAULT '',
	    meta_reason TEXT DEFAULT ''
	);
CREATE TABLE IF NOT EXISTS tuples_weight_faction_to_action_type (
	    subject VARCHAR(255) NOT NULL, 
	    object VARCHAR(255) NOT NULL,
	    value INTEGER NOT NULL DEFAULT 0,
	    UNIQUE(subject, object)
	);
CREATE TABLE IF NOT EXISTS modifiers_weight_faction_to_action_type (
	    subject VARCHAR(255) NOT NULL,
	    object VARCHAR(255) NOT NULL,
	    value INTEGER NOT NULL DEFAULT 0,
	    tick_expires INTEGER NOT NULL DEFAULT 0,
	    meta_key VARCHAR(255) DEFAULT '',
	    meta_val VARCHAR(255) DEFAULT '',
	    meta_reason TEXT DEFAULT ''
	);
CREATE TABLE IF NOT EXISTS tuples_weight_faction_to_profession (
	    subject VARCHAR(255) NOT NULL, 
	    object VARCHAR(255) NOT NULL,
	    value INTEGER NOT NULL DEFAULT 0,
	    UNIQUE(subject, object)
	);
CREATE TABLE IF NOT EXISTS tuples_rank_person_to_faction (
	    subject VARCHAR(255) NOT NULL, 
	    object VARCHAR(255) NOT NULL,
	    value INTEGER NOT NULL DEFAULT 0,
	    UNIQUE(subject, object)
	);
CREATE TABLE IF NOT EXISTS tuples_intelligence_faction_to_faction (
	    subject VARCHAR(255) NOT NULL, 
	    object VARCHAR(255) NOT NULL,
	    value INTEGER NOT NULL DEFAULT 0,
	    UNIQUE(subject, object)
	);

--- Create indexes
CREATE INDEX IF NOT EXISTS laws_source_id ON laws (source_id);
CREATE INDEX IF NOT EXISTS plot_owner ON plots (faction_id);
CREATE INDEX IF NOT EXISTS jobs_target_area ON jobs (target_area_id, state);
CREATE INDEX IF NOT EXISTS event_tick_type ON events (tick, type);
CREATE INDEX IF NOT EXISTS fam_child_bearing ON families (area_id, is_child_bearing);
CREATE INDEX IF NOT EXISTS peo_area ON people (area_id);
CREATE INDEX IF NOT EXISTS peo_prof ON people (area_id, preferred_profession);
CREATE INDEX IF NOT EXISTS peo_job ON people (job_id);
CREATE INDEX IF NOT EXISTS peo_fam ON people (birth_family_id);
CREATE INDEX IF NOT EXISTS peo_rand ON people (random);
CREATE INDEX IF NOT EXISTS fact_government ON factions (government_id);
CREATE INDEX IF NOT EXISTS fact_action_fq ON factions (action_frequency_ticks);
CREATE INDEX IF NOT EXISTS subject_tuples_trust_faction_to_faction ON tuples_trust_faction_to_faction (subject);
CREATE INDEX IF NOT EXISTS subject_modifiers_trust_faction_to_faction ON modifiers_trust_faction_to_faction (subject);
CREATE INDEX IF NOT EXISTS subject_tuples_import_faction_to_commodity ON tuples_import_faction_to_commodity (subject);
CREATE INDEX IF NOT EXISTS subject_tuples_export_faction_to_commodity ON tuples_export_faction_to_commodity (subject);
CREATE INDEX IF NOT EXISTS subject_tuples_research_faction_to_topic ON tuples_research_faction_to_topic (subject);
CREATE INDEX IF NOT EXISTS subject_tuples_weight_faction_to_topic ON tuples_weight_faction_to_topic (subject);
CREATE INDEX IF NOT EXISTS subject_tuples_affiliation_person_to_faction ON tuples_affiliation_person_to_faction (subject);
CREATE INDEX IF NOT EXISTS subject_modifiers_affiliation_person_to_faction ON modifiers_affiliation_person_to_faction (subject);
CREATE INDEX IF NOT EXISTS subject_tuples_skill_person_to_profession ON tuples_skill_person_to_profession (subject);
CREATE INDEX IF NOT EXISTS subject_tuples_relationship_person_to_person ON tuples_relationship_person_to_person (subject);
CREATE INDEX IF NOT EXISTS subject_tuples_faith_person_to_religion ON tuples_faith_person_to_religion (subject);
CREATE INDEX IF NOT EXISTS subject_tuples_trust_person_to_person ON tuples_trust_person_to_person (subject);
CREATE INDEX IF NOT EXISTS subject_modifiers_trust_person_to_person ON modifiers_trust_person_to_person (subject);
CREATE INDEX IF NOT EXISTS subject_tuples_weight_faction_to_action_type ON tuples_weight_faction_to_action_type (subject);
CREATE INDEX IF NOT EXISTS subject_modifiers_weight_faction_to_action_type ON modifiers_weight_faction_to_action_type (subject);
CREATE INDEX IF NOT EXISTS subject_tuples_weight_faction_to_profession ON tuples_weight_faction_to_profession (subject);
CREATE INDEX IF NOT EXISTS subject_tuples_rank_person_to_faction ON tuples_rank_person_to_faction (subject);
CREATE INDEX IF NOT EXISTS subject_tuples_intelligence_faction_to_faction ON tuples_intelligence_faction_to_faction (subject);

--- Hand out permissions
GRANT SELECT ON ALL TABLES IN SCHEMA faction TO factionreadonly;
GRANT SELECT ON ALL SEQUENCES IN SCHEMA faction TO factionreadonly;
GRANT ALL ON ALL TABLES IN SCHEMA faction TO factionreadwrite;
GRANT ALL ON ALL SEQUENCES IN SCHEMA faction TO factionreadwrite;
