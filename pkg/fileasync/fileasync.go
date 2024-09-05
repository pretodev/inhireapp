package fileasync

import (
	"context"
	"io/fs"
	"os"

	"github.com/dgraph-io/badger/v4"
)

var badgerManager = NewConnectionManager[*badger.DB]()

func BadgerPool(ctx context.Context, options badger.Options, opts ...poolOption) (*Connection[*badger.DB], error) {
	return badgerManager.NewPool(
		ctx,
		options.Dir,
		func(path string) (*badger.DB, error) { return badger.Open(options) },
		opts...,
	)
}

var fileManager = NewConnectionManager[*os.File]()

func FilePool(ctx context.Context, path string, flag int, perm fs.FileMode, opts ...poolOption) (*Connection[*os.File], error) {
	return fileManager.NewPool(
		ctx,
		path,
		func(path string) (*os.File, error) { return os.OpenFile(path, flag, perm) },
		opts...,
	)
}
