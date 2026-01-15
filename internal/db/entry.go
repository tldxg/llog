package db

import (
	"context"
	"time"

	"github.com/thdxg/llog/internal/model"
	"gorm.io/gorm"
)

type entryDB struct {
	i gorm.Interface[model.Entry]
}

func (db *entryDB) Count(ctx context.Context) (int64, error) {
	return db.i.Count(ctx, "id")
}

func (db *entryDB) WithRange(from time.Time, to time.Time) gorm.ChainInterface[model.Entry] {
	if from.IsZero() {
		return db.i.Where("created_at <= ?", to)
	}
	if to.IsZero() {
		return db.i.Where("created_at >= ?", from)
	}
	return db.i.Where("created_at >= ? AND created_at <= ?", from, to)
}

func (db *entryDB) WithIds(ids []uint64) gorm.ChainInterface[model.Entry] {
	return db.i.Where("id IN ?", ids)
}

func (db *entryDB) Add(ctx context.Context, entries []model.Entry) error {
	return db.i.CreateInBatches(ctx, &entries, len(entries))
}

func (db *entryDB) GetLast(ctx context.Context) (model.Entry, error) {
	return db.i.Order("created_at desc").Last(ctx)
}

// to not apply any filter, subquery = nil
// to fetch all entries, n = -1
func (db *entryDB) Get(ctx context.Context, subquery gorm.ChainInterface[model.Entry], n int) ([]model.Entry, error) {
	var limited gorm.ChainInterface[model.Entry]
	if subquery == nil {
		limited = db.i.Order("created_at desc").Limit(n)
	} else {
		limited = subquery.Order("created_at desc").Limit(n)
	}
	return db.i.Table("(?) as limited", limited).Order("created_at").Find(ctx)
}

func (db *entryDB) Delete(ctx context.Context, subquery gorm.ChainInterface[model.Entry], n int) (int, error) {
	var limited gorm.ChainInterface[model.Entry]
	if subquery == nil {
		limited = db.i.Order("created_at desc").Limit(n)
	} else {
		limited = subquery.Order("created_at desc").Limit(n)
	}

	toDelete, err := limited.Find(ctx)
	if err != nil {
		return 0, err
	}

	ids := make([]uint64, len(toDelete))
	for i, e := range toDelete {
		ids[i] = e.ID
	}

	return db.i.Where("id IN ?", ids).Delete(ctx)
}

func (db *entryDB) Nuke(ctx context.Context) error {

	return nil
}
