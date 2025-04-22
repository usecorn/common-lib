package dbutils

import (
	"context"
	"strconv"

	"github.com/cockroachdb/errors"
	"github.com/jmoiron/sqlx"
)

type MetaDB interface {
	Set(ctx context.Context, key string, val any) error
	GetString(ctx context.Context, key string) (out string, err error)
	GetInt64(ctx context.Context, key string) (int64, error)
	GetUint64(ctx context.Context, key string) (uint64, error)
}

type metaDB struct {
	db      *sqlx.DB
	getMeta *sqlx.Stmt
	setMeta *sqlx.Stmt
}

func NewMetaDB(sdb *sqlx.DB) (MetaDB, error) {
	getMeta, err := sdb.Preparex("SELECT value FROM meta WHERE key = $1")
	if err != nil {
		return nil, err
	}
	setMeta, err := sdb.Preparex("INSERT INTO meta (key,value) VALUES($1,$2) ON CONFLICT (key) DO UPDATE SET value = $2")
	if err != nil {
		return nil, err
	}

	return &metaDB{
		db:      sdb,
		getMeta: getMeta,
		setMeta: setMeta,
	}, nil
}

func (meta *metaDB) Set(ctx context.Context, key string, val any) error {
	var err error
	switch v := val.(type) {
	case string:
		_, err = meta.setMeta.ExecContext(ctx, key, v)
	case int:
		_, err = meta.setMeta.ExecContext(ctx, key, strconv.Itoa(v))
	case int64:
		_, err = meta.setMeta.ExecContext(ctx, key, strconv.FormatInt(v, 10))
	case uint64:
		_, err = meta.setMeta.ExecContext(ctx, key, strconv.FormatUint(v, 10))
	default:
		err = errors.New("unsupported type")
	}

	return err
}

func (meta *metaDB) GetString(ctx context.Context, key string) (out string, err error) {
	return out, meta.getMeta.GetContext(ctx, &out, key)
}

func (meta *metaDB) GetInt64(ctx context.Context, key string) (int64, error) {
	rawVal, err := meta.GetString(ctx, key)
	if err != nil {
		return -1, err
	}
	return strconv.ParseInt(rawVal, 10, 64)
}

func (meta *metaDB) GetUint64(ctx context.Context, key string) (uint64, error) {
	rawVal, err := meta.GetString(ctx, key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(rawVal, 10, 64)
}
