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

// splitUUIDsCSV convertit "AgendaId1,AgendaId2,AgendaId3" -> []uuid.UUID (tolère vide)
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

// scanEvent scanne une ligne SQL et la convertit en models.Event.
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


// GetEventsByAgenda retourne tous les évènements liés à un agenda donné.
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

// GetEventById returns a single event (including its agenda IDs) by event id.
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
