package main

/*
Trigger to track changes to people table
- Nice that it's automattic, discarded because events allow us to keep *why* a change happened
- .. unless we add a reason column to the person (or whatever) table ..
 person
    value_a X
    value_a_reasonkey metakey
    value_a_reasonval metaval
    [...]
which seems like way more work than just using events


CREATE TABLE IF NOT EXISTS changes (
    source TEXT NOT NULL,
    tick INTEGER NOT NULL DEFAULT 0,
    old_val TEXT,
    new_val TEXT
)

DROP TRIGGER updates_people;

CREATE TRIGGER IF NOT EXISTS updates_people AFTER UPDATE ON people FOR EACH ROW
BEGIN
	INSERT INTO changes (source, tick, old_val, new_val)
	VALUES (
		"people",
		IFNULL((SELECT int FROM meta where id='tick' LIMIT 1), 0),
		json_object(
			'id', OLD.id
		),
		json_object(
			'id', NEW.id
		)
	);
END

^ Don't think we can dynamically json-fy the whole "new" or "old" row, so each table would need it's own
trigger with all the column names :(
*/
