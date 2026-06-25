package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/degoke/health-ai-stack/pkg/store"
)

// SearchStore persists typed search index rows in SQLite.
type SearchStore struct {
	exec searchExec
}

type searchExec interface {
	queryExec
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type searchTable string

const (
	searchTableToken     searchTable = "search_token"
	searchTableString    searchTable = "search_string"
	searchTableDate      searchTable = "search_date"
	searchTableNumber    searchTable = "search_number"
	searchTableReference searchTable = "search_reference"
)

func newSearchStore(db *sql.DB) *SearchStore {
	return &SearchStore{exec: db}
}

func newSearchStoreTx(tx *sql.Tx) *SearchStore {
	return &SearchStore{exec: tx}
}

func parseSearchFieldKey(key string) (searchTable, string, error) {
	parts := strings.SplitN(key, ".", 2)
	if len(parts) == 1 {
		return searchTableString, key, nil
	}
	switch parts[0] {
	case "token":
		return searchTableToken, parts[1], nil
	case "string":
		return searchTableString, parts[1], nil
	case "date":
		return searchTableDate, parts[1], nil
	case "number":
		return searchTableNumber, parts[1], nil
	case "reference", "ref":
		return searchTableReference, parts[1], nil
	default:
		return searchTableString, key, nil
	}
}

func (s *SearchStore) Index(ctx context.Context, entry store.SearchIndexEntry) error {
	for fieldKey, value := range entry.Fields {
		table, normalizedKey, err := parseSearchFieldKey(fieldKey)
		if err != nil {
			return err
		}
		query := fmt.Sprintf(`
			INSERT OR IGNORE INTO %s (resource_type, resource_id, field_key, value)
			VALUES (?, ?, ?, ?)`, table)
		if _, err := s.exec.ExecContext(ctx, query, entry.ResourceType, entry.ID, normalizedKey, value); err != nil {
			return fmt.Errorf("index search field %s: %w", fieldKey, err)
		}
	}
	return nil
}

func (s *SearchStore) RemoveIndex(ctx context.Context, resourceType, id string) error {
	tables := []searchTable{
		searchTableToken,
		searchTableString,
		searchTableDate,
		searchTableNumber,
		searchTableReference,
	}
	for _, table := range tables {
		query := fmt.Sprintf(`DELETE FROM %s WHERE resource_type = ? AND resource_id = ?`, table)
		if _, err := s.exec.ExecContext(ctx, query, resourceType, id); err != nil {
			return fmt.Errorf("remove search index from %s: %w", table, err)
		}
	}
	return nil
}

func (s *SearchStore) Lookup(ctx context.Context, key, value string) ([]string, error) {
	table, fieldKey, err := parseSearchFieldKey(key)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
		SELECT resource_id FROM %s
		WHERE field_key = ? AND value = ?
		ORDER BY resource_id`, table)
	rows, err := s.exec.QueryContext(ctx, query, fieldKey, value)
	if err != nil {
		return nil, fmt.Errorf("lookup search index: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan search lookup: %w", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate search lookup: %w", err)
	}
	return ids, nil
}

func (s *SearchStore) QueryPrepared(ctx context.Context, query store.PreparedQuery, args map[string]string) ([]string, error) {
	if query.Name == "by-field" {
		return s.Lookup(ctx, args["key"], args["value"])
	}
	return nil, nil
}
