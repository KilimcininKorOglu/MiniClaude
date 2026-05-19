package db

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type migrationFile struct {
	Version int64
	Name    string
	Path    string
	Up      bool
}

var migrationNamePattern = regexp.MustCompile(`^(\d+)_(.+)\.(up|down)\.sql$`)

func MigrateUp(ctx context.Context, pool *pgxpool.Pool, dir string) error {
	files, err := loadMigrationFiles(dir, true)
	if err != nil {
		return err
	}

	return withTx(ctx, pool, func(tx pgx.Tx) error {
		if err := ensureMigrationTable(ctx, tx); err != nil {
			return err
		}

		applied, err := appliedVersions(ctx, tx)
		if err != nil {
			return err
		}

		for _, file := range files {
			if applied[file.Version] {
				continue
			}
			if err := applyMigrationFile(ctx, tx, file); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, `insert into schema_migrations (version, name) values ($1, $2)`, file.Version, file.Name); err != nil {
				return fmt.Errorf("record migration %s: %w", file.Path, err)
			}
		}

		return nil
	})
}

func MigrateDown(ctx context.Context, pool *pgxpool.Pool, dir string) error {
	files, err := loadMigrationFiles(dir, false)
	if err != nil {
		return err
	}

	byVersion := make(map[int64]migrationFile, len(files))
	for _, file := range files {
		byVersion[file.Version] = file
	}

	return withTx(ctx, pool, func(tx pgx.Tx) error {
		if err := ensureMigrationTable(ctx, tx); err != nil {
			return err
		}

		var version int64
		if err := tx.QueryRow(ctx, `select version from schema_migrations order by version desc limit 1`).Scan(&version); err != nil {
			if err == pgx.ErrNoRows {
				return nil
			}
			return fmt.Errorf("read latest migration: %w", err)
		}

		file, ok := byVersion[version]
		if !ok {
			return fmt.Errorf("missing down migration for version %d", version)
		}

		if err := applyMigrationFile(ctx, tx, file); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `delete from schema_migrations where version = $1`, version); err != nil {
			return fmt.Errorf("delete migration record %d: %w", version, err)
		}

		return nil
	})
}

func ensureMigrationTable(ctx context.Context, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, `
		create table if not exists schema_migrations (
			version bigint primary key,
			name text not null,
			applied_at timestamptz not null default now()
		)
	`)
	if err != nil {
		return fmt.Errorf("ensure schema_migrations: %w", err)
	}
	return nil
}

func appliedVersions(ctx context.Context, tx pgx.Tx) (map[int64]bool, error) {
	rows, err := tx.Query(ctx, `select version from schema_migrations`)
	if err != nil {
		return nil, fmt.Errorf("read applied migrations: %w", err)
	}
	defer rows.Close()

	versions := map[int64]bool{}
	for rows.Next() {
		var version int64
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("scan applied migration: %w", err)
		}
		versions[version] = true
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate applied migrations: %w", err)
	}
	return versions, nil
}

func loadMigrationFiles(dir string, up bool) ([]migrationFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read migrations directory %q: %w", dir, err)
	}

	files := []migrationFile{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		matches := migrationNamePattern.FindStringSubmatch(entry.Name())
		if matches == nil {
			continue
		}

		isUp := matches[3] == "up"
		if isUp != up {
			continue
		}

		version, err := strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse migration version %q: %w", entry.Name(), err)
		}

		files = append(files, migrationFile{
			Version: version,
			Name:    strings.ReplaceAll(matches[2], "_", " "),
			Path:    filepath.Join(dir, entry.Name()),
			Up:      isUp,
		})
	}

	sort.Slice(files, func(i, j int) bool {
		if up {
			return files[i].Version < files[j].Version
		}
		return files[i].Version > files[j].Version
	})

	return files, nil
}

func applyMigrationFile(ctx context.Context, tx pgx.Tx, file migrationFile) error {
	content, err := os.ReadFile(file.Path)
	if err != nil {
		return fmt.Errorf("read migration %s: %w", file.Path, err)
	}
	if _, err := tx.Exec(ctx, string(content)); err != nil {
		return fmt.Errorf("apply migration %s: %w", file.Path, err)
	}
	return nil
}

func withTx(ctx context.Context, pool *pgxpool.Pool, fn func(pgx.Tx) error) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if err := fn(tx); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}
