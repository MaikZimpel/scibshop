// +build mongo

package filestore

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"scib-svr/configuration"
	"scib-svr/logging"
)

type filestore struct {
	client *mongo.Client
	db     string
	log    logging.Logger
}

func (f filestore) Upload(c context.Context, file []byte, fileExt string) (id string, err error) {
	bucket, err := gridfs.NewBucket(f.client.Database(f.db))
	if err == nil {
		id = uuid.New().String()
		fn := id + fileExt
		err := bucket.UploadFromStreamWithID(id, id+fileExt, bytes.NewBuffer(file))
		if err == nil {
			f.log.Debug(c, "new file [%s] uploaded to filestore with id %v [%d Byte]", fn, id, len(file))
		}
	}
	if err != nil {
		f.log.Error(c, "upload of file [%s] failed because %v", id, err)
	}
	return
}

func (f filestore) Download(c context.Context, fileId string) (size int64, file []byte, err error) {
	bucket, err := gridfs.NewBucket(f.client.Database(f.db))
	if err == nil {
		var buf bytes.Buffer
		size, err = bucket.DownloadToStream(fileId, &buf)
		file = buf.Bytes()
	} else {
		f.log.Error(c, "download of file [%s] failed because %v", fileId, err)
	}
	return
}

func (f filestore) Remove(c context.Context, fileId string) error {
	bucket, err := gridfs.NewBucket(f.client.Database(f.db))
	if err == nil {
		return bucket.Delete(fileId)
	} else {
		f.log.Error(c, "removal of file %s failed because %v", fileId, err)
		return err
	}
}

func NewMFilestore() Filestore {

	clientOptions := options.Client().ApplyURI(
		fmt.Sprintf("mongodb://%s:%d", configuration.MongoDbHost, configuration.MongoDbPort))
	clientOptions.SetMaxPoolSize(5)
	client, err := mongo.Connect(context.Background(), clientOptions)
	log := logging.New()
	if err != nil {
		log.Critical(context.Background(), "unable to connect to mongodb %v", err)
	} else {
		log.Info(context.Background(), "connected to database at %s",
			fmt.Sprintf("mongodb://%s:%d", configuration.MongoDbHost, configuration.MongoDbPort))
	}

	return &filestore{
		client,
		configuration.MongoDbDatabase,
		log,
	}
}

func init() {
	New = NewMFilestore
}
