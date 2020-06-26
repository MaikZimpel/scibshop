package filestore

import "context"

type (
	Filestore interface {
		Upload(c context.Context, file []byte, fileExt string) (string, error)
		Download(c context.Context, id string) (int64, []byte, error)
		Remove(c context.Context, id string) error
	}
	filestoreFactory func() Filestore
)

var (
	New filestoreFactory
)