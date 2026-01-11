package events

import (
	"database/sql"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"middleware/example/internal/helpers"
	"middleware/example/internal/models"
)


type scanner interface {
	Scan(dest ...any) error
}

func splitUUIDsCSV(csv string) ([]uuid.UUID, error) {
	csv = strings.TrimSpace(csv)
	if csv == "" {
		return []uuid.UUID{}, nil
	}
	parts := strings.Split(csv, ",")
	out := make([]uuid.UUID, 0, len(parts))
	for _, p := range parts {
		u, err := uuid.FromString(strings.TrimSpace(p))
		if err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, nil
}

func stringFromNull(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func parseRFC3339(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func parseNullableRFC3339(s string) (*time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func scanEvent(s scanner) (models.Event, error) {
	var (
		idStr, uid, startStr, endStr string
		descNull, nameNull, locNull, lastUpdNull, csvIDsNull sql.NullString
	)

	if err := s.Scan( &idStr, &uid, &descNull, &nameNull, &startStr,
					  &endStr, &locNull, &lastUpdNull, &csvIDsNull,
	); err != nil {
		return models.Event{}, err
	}

	evID, err := uuid.FromString(idStr)
	if err != nil {
		return models.Event{}, err
	}

	start, err := parseRFC3339(startStr)
	if err != nil {
		return models.Event{}, err
	}

	end, err := parseRFC3339(endStr)
	if err != nil {
		return models.Event{}, err
	}

	lastUpdate, err := parseNullableRFC3339(stringFromNull(lastUpdNull))
	if err != nil {
		return models.Event{}, err
	}

	agendaIDs, err := splitUUIDsCSV(stringFromNull(csvIDsNull))
	if err != nil {
		return models.Event{}, err
	}

	return models.Event{
		ID:          evID,
		AgendaIDs:   agendaIDs,
		UID:         uid,
		Description: stringFromNull(descNull),
		Name:        stringFromNull(nameNull),
		Start:       start,
		End:         end,
		Location:    stringFromNull(locNull),
		LastUpdate:  lastUpdate,
	}, nil
}


func GetEventsByAgenda(agendaID uuid.UUID) ([]models.Event, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	const q = `
		SELECT e.id, e.uid, e.description, e.name, e.start, e."end", e.location, e.last_update, GROUP_CONCAT(ea_all.agenda_id) AS agenda_ids
		FROM events e
		JOIN events_agendas ea_filter
			ON ea_filter.event_id = e.id
			AND ea_filter.agenda_id = ?
		LEFT JOIN events_agendas ea_all
			ON ea_all.event_id = e.id
		GROUP BY e.id, e.uid, e.description, e.name, e.start, e."end", e.location, e.last_update
		ORDER BY e.start;
	`

	rows, err := db.Query(q, agendaID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Event
	for rows.Next() {
		ev, err := scanEvent(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, ev)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func GetAllEvents() ([]models.Event, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	const q = `
		SELECT e.id, e.uid, e.description, e.name, e.start, e."end", e.location, e.last_update, GROUP_CONCAT(ea.agenda_id) AS agenda_ids
		FROM events e
		LEFT JOIN events_agendas ea
			ON ea.event_id = e.id
		GROUP BY e.id, e.uid, e.description, e.name, e.start, e."end", e.location, e.last_update
		ORDER BY e.start;
	`

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Event
	for rows.Next() {
		ev, err := scanEvent(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, ev)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func GetEventById(id uuid.UUID) (*models.Event, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	const q = `
		SELECT e.id, e.uid, e.description, e.name, e.start, e."end", e.location, e.last_update, GROUP_CONCAT(ea_all.agenda_id) AS agenda_ids
		FROM events e
		LEFT JOIN events_agendas ea_all
			ON ea_all.event_id = e.id
		WHERE e.id = ?
		GROUP BY e.id, e.uid, e.description, e.name, e.start, e."end", e.location, e.last_update;
	`

	row := db.QueryRow(q, id.String())

	ev, err := scanEvent(row)
	if err != nil {
		return nil, err
	}

	return &ev, nil
}

func GetEventByUID(uid string) (*models.Event, error) {
    db, err := helpers.OpenDB()
    if err != nil {
        return nil, err
    }

	
    defer helpers.CloseDB(db)

    const q = `
        SELECT e.id, e.uid, e.description, e.name, e.start, e."end", e.location, e.last_update, GROUP_CONCAT(ea.agenda_id) AS agenda_ids
        FROM events e
        LEFT JOIN events_agendas ea ON ea.event_id = e.id
        WHERE e.uid = ?
        GROUP BY e.id, e.uid, e.description, e.name, e.start, e."end", e.location, e.last_update;
    `

    row := db.QueryRow(q, uid)
	ev, err := scanEvent(row)
	if err != nil {
		return nil, err
	}

	return &ev, nil
}

func UpdateEventByUID(UID string, agendaIDs []uuid.UUID, description string, name string, Start time.Time, End time.Time, location string, LastUpdate *time.Time) (*models.Event, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}

	var lastUpd any
    if LastUpdate == nil {
        lastUpd = nil
    } else {
        lastUpd = LastUpdate.UTC().Format(time.RFC3339)
    }

	q := `
		UPDATE events
		SET description = ?, name = ?, start = ?, "end" = ?, location = ?, last_update = ?
		WHERE uid = ?
		`

	_, err = db.Exec(
		q,
		description,
		name,
		Start.UTC().Format(time.RFC3339),
		End.UTC().Format(time.RFC3339),
		location,
		lastUpd,
		UID,
	)
	helpers.CloseDB(db)

	if err != nil {
		return nil, err
	}

	updatedEvent := &models.Event{
		AgendaIDs: agendaIDs,
		Description: description,
		Name:  name,
		Start: Start,
		End: End,
		Location: location,
		LastUpdate: LastUpdate,
		UID: UID,
	}

	return updatedEvent, nil
}

func EventAgendasLink(eventID uuid.UUID, agendaIDs []uuid.UUID) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	const qDelete = `
		DELETE FROM events_agendas
		WHERE event_id = ?;
	`

	_, err = tx.Exec(qDelete, eventID.String())
	if err != nil {
		return err
	}	
	
	const qInsert = `
		INSERT INTO events_agendas (event_id, agenda_id)
		VALUES (?, ?);
	`

	for _, agendaID := range agendaIDs {
		_, err = tx.Exec(qInsert, eventID.String(), agendaID.String())
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func InsertEvent(event models.Event) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	const insertEventQ = `
		INSERT INTO events (id, uid, description, name, start, "end", location, last_update)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?);
	`
	if event.ID == uuid.Nil {
		event.ID = uuid.Must(uuid.NewV4())
	}

	var lastUpd any
    if event.LastUpdate == nil {
        lastUpd = nil
    } else {
        lastUpd = event.LastUpdate.UTC().Format(time.RFC3339)
    }

	_, err = tx.Exec(
		insertEventQ,
		event.ID.String(),
		event.UID,
		event.Description,
		event.Name,
		event.Start.UTC().Format(time.RFC3339),
		event.End.UTC().Format(time.RFC3339),
		event.Location,
		lastUpd,
	)
	if err != nil {
		return err
	}

	const insertLinkQ = `
		INSERT OR IGNORE INTO events_agendas (event_id, agenda_id)
		VALUES (?, ?);
	`
	for _, agendaID := range event.AgendaIDs {
		_, err = tx.Exec(insertLinkQ, event.ID.String(), agendaID.String())
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}