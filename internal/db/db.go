package db

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/thdxg/llog/internal/config"
	"github.com/thdxg/llog/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	db    *gorm.DB
	Entry *entryDB
}

func Load(cfg *config.Config, ctx context.Context, db *DB) error {
	dir := filepath.Dir(cfg.DBPath)

	if err := os.Mkdir(dir, 0755); err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to create db directory: %w", err)
	}

	gormdb, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}

	if err := gormdb.AutoMigrate(&model.Entry{}); err != nil {
		return fmt.Errorf("failed to migrate Entry: %w", err)
	}

	db.db = gormdb
	db.Entry = &entryDB{gorm.G[model.Entry](gormdb)}

	count, err := db.Entry.Count(ctx)
	if err != nil {
		return err
	}

	cfg.Internal.EntryCount = count

	if count > 0 {
		last, err := db.Entry.GetLast(ctx)
		if err != nil {
			return err
		}

		cfg.Internal.MaxEntryId = last.ID
	}

	return nil
}

func (db *DB) Nuke() error {
	return db.db.Migrator().DropTable(&model.Entry{})
}
