package dbutils

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/cockroachdb/errors"
)

type metaDBMemory struct {
	data map[string]string
}

// NewMetaDBMemory creates a new MetaDB instance with in-memory storage.
// This is useful for testing purposes or when you don't need persistent storage.
func NewMetaDBMemory() (MetaDB, error) {

	return &metaDBMemory{
		data: map[string]string{},
	}, nil
}

func (meta *metaDBMemory) Set(ctx context.Context, key string, val any) error {
	switch v := val.(type) {
	case string:
		meta.data[key] = v
	case int:
		meta.data[key] = strconv.Itoa(v)
	case int64:
		meta.data[key] = strconv.FormatInt(v, 10)
	case uint64:
		meta.data[key] = strconv.FormatUint(v, 10)
	default:
		return errors.New("unsupported type")
	}

	return nil
}

func (meta *metaDBMemory) GetString(ctx context.Context, key string) (out string, err error) {
	val, ok := meta.data[key]
	if !ok {
		return "", sql.ErrNoRows
	}
	return val, nil
}

func (meta *metaDBMemory) GetInt64(ctx context.Context, key string) (int64, error) {
	rawVal, err := meta.GetString(ctx, key)
	if err != nil {
		return -1, err
	}
	return strconv.ParseInt(rawVal, 10, 64)
}

func (meta *metaDBMemory) GetUint64(ctx context.Context, key string) (uint64, error) {
	rawVal, err := meta.GetString(ctx, key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(rawVal, 10, 64)
}
