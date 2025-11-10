package events

import (
	"database/sql"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"middleware/example/internal/helpers"
	"middleware/example/internal/models"
)

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

// GetEventsByAgenda retourne tous les évènements liés à un agenda donné.
// Filtrage en SQL, et agrégation de TOUS les agendaIds liés à chaque event.
func GetEventsByAgenda(agendaID uuid.UUID) ([]models.Event, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	// On filtre avec un JOIN "ea_filter" (agenda demandé),
	// et on agrège tous les agendas liés via "ea_all".
	const q = `
			SELECT e.id, e.uid, e.description, e.name, e.start, e."end", e.location, e.last_update,
			GROUP_CONCAT(ea_all.agenda_id) AS agenda_ids
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
		var (
			idStr, uid, startStr, endStr string
			descNull, nameNull, locNull, lastUpdNull, csvIDsNull sql.NullString
		)
		if err := rows.Scan(&idStr, &uid, &descNull, &nameNull, &startStr, &endStr, &locNull, &lastUpdNull, &csvIDsNull); err != nil {
			return nil, err
		}

		// safe conversions from sql.NullString
		desc := ""
		if descNull.Valid {
			desc = descNull.String
		}
		name := ""
		if nameNull.Valid {
			name = nameNull.String
		}
		loc := ""
		if locNull.Valid {
			loc = locNull.String
		}
		lastUpdStr := ""
		if lastUpdNull.Valid {
			lastUpdStr = lastUpdNull.String
		}
		csvIDs := ""
		if csvIDsNull.Valid {
			csvIDs = csvIDsNull.String
		}

		// Conversions
		evID, err := uuid.FromString(idStr)
		if err != nil {
			return nil, err
		}

		start, err := time.Parse(time.RFC3339, startStr)
		if err != nil {
			return nil, err
		}
		end, err := time.Parse(time.RFC3339, endStr)
		if err != nil {
			return nil, err
		}

		var lastUpdate *time.Time
		if s := strings.TrimSpace(lastUpdStr); s != "" {
			t, err := time.Parse(time.RFC3339, s)
			if err != nil {
				return nil, err
			}
			lastUpdate = &t
		}

		agendaIDs, err := splitUUIDsCSV(csvIDs)
		if err != nil {
			return nil, err
		}

		out = append(out, models.Event{
			ID:          evID,
			AgendaIDs:   agendaIDs,
			UID:         uid,
			Description: desc,
			Name:        name,
			Start:       start,
			End:         end,
			Location:    loc,
			LastUpdate:  lastUpdate,
		})
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
		SELECT e.id, e.uid, e.description, e.name, e.start, e."end", e.location, e.last_update,
		GROUP_CONCAT(ea_all.agenda_id) AS agenda_ids
		FROM events e
		LEFT JOIN events_agendas ea_all
		ON ea_all.event_id = e.id
		WHERE e.id = ?
		GROUP BY e.id, e.uid, e.description, e.name, e.start, e."end", e.location, e.last_update;
	`

	row := db.QueryRow(q, id.String())

	var (
		idStr, uid, startStr, endStr string
		descNull, nameNull, locNull, lastUpdNull, csvIDsNull sql.NullString
	)

	if err := row.Scan(&idStr, &uid, &descNull, &nameNull, &startStr, &endStr, &locNull, &lastUpdNull, &csvIDsNull); err != nil {
		return nil, err
	}

	// safe conversions from sql.NullString
	desc := ""
	if descNull.Valid {
		desc = descNull.String
	}
	name := ""
	if nameNull.Valid {
		name = nameNull.String
	}
	loc := ""
	if locNull.Valid {
		loc = locNull.String
	}
	lastUpdStr := ""
	if lastUpdNull.Valid {
		lastUpdStr = lastUpdNull.String
	}
	csvIDs := ""
	if csvIDsNull.Valid {
		csvIDs = csvIDsNull.String
	}

	// Conversions
	evID, err := uuid.FromString(idStr)
	if err != nil {
		return nil, err
	}

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		return nil, err
	}

	var lastUpdate *time.Time
	if s := strings.TrimSpace(lastUpdStr); s != "" {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return nil, err
		}
		lastUpdate = &t
	}

	agendaIDs, err := splitUUIDsCSV(csvIDs)
	if err != nil {
		return nil, err
	}

	return &models.Event{
		ID:         evID,
		AgendaIDs:  agendaIDs,
		UID:        uid,
		Description: desc,
		Name:       name,
		Start:      start,
		End:        end,
		Location:   loc,
		LastUpdate: lastUpdate,
	}, nil
}