package dbutils

import (
	"strconv"

	"github.com/cockroachdb/errors"
	"github.com/jmoiron/sqlx"
)

type MetaDB interface {
	Set(key string, val any) error
	GetString(key string) (out string, err error)
	GetInt64(key string) (int64, error)
	GetUint64(key string) (uint64, error)
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

func (meta *metaDB) Set(key string, val any) error {
	var err error
	switch v := val.(type) {
	case string:
		_, err = meta.setMeta.Exec(key, v)
	case int:
		_, err = meta.setMeta.Exec(key, strconv.Itoa(v))
	case int64:
		_, err = meta.setMeta.Exec(key, strconv.FormatInt(v, 10))
	case uint64:
		_, err = meta.setMeta.Exec(key, strconv.FormatUint(v, 10))
	default:
		err = errors.New("unsupported type")
	}

	return err
}

func (meta *metaDB) GetString(key string) (out string, err error) {
	return out, meta.getMeta.Get(&out, key)
}

func (meta *metaDB) GetInt64(key string) (int64, error) {
	rawVal, err := meta.GetString(key)
	if err != nil {
		return -1, err
	}
	return strconv.ParseInt(rawVal, 10, 64)
}

func (meta *metaDB) GetUint64(key string) (uint64, error) {
	rawVal, err := meta.GetString(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(rawVal, 10, 64)
}
