package helpers

import "middleware/example/internal/models"

type ChangeLine struct {
	FieldLabel string
	Before     string
	After      string
}

func FieldLabel(f string) string {
	switch f {
	case "start":
		return "Début"
	case "end":
		return "Fin"
	case "location":
		return "Salle"
	case "name":
		return "Cours"
	case "description":
		return "Description"
	default:
		return f
	}
}

func MapChanges(changes []models.FieldChange) []ChangeLine {
	out := make([]ChangeLine, 0, len(changes))
	for _, c := range changes {
		out = append(out, ChangeLine{
			FieldLabel: FieldLabel(c.Field),
			Before:     c.Before,
			After:      c.After,
		})
	}
	return out
}
