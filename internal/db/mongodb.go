package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/martinlindhe/base36"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	v1 "github.com/voidshard/faction/pkg/structs/v1"
	"github.com/voidshard/faction/pkg/util/log"
	"github.com/voidshard/faction/pkg/util/uuid"
)

const (
	colWorlds = "world"
)

type Mongo struct {
	opts *options.ClientOptions
	log  log.Logger
	cfg  *MongoConfig
	conn *mongo.Client
}

type MongoConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	Database string
}

func NewMongo(cfg *MongoConfig) (*Mongo, error) {
	url := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	bsonOpts := &options.BSONOptions{UseJSONStructTags: true}

	opts := options.Client().ApplyURI(url).SetBSONOptions(bsonOpts)
	// github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo/example_test.go
	// opts.Monitor = otelmongo.NewMonitor()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	me := &Mongo{
		opts: opts,
		log: log.Sublogger("mongodb", map[string]interface{}{
			"host":     cfg.Host,
			"port":     cfg.Port,
			"database": cfg.Database,
		}),
		cfg:  cfg,
		conn: client,
	}
	go me.ping()
	return me, nil
}

func (m *Mongo) connect() {
	// if connection fails, try to reconnect
	i := 0
	for {
		i += 1
		if i > 1 {
			time.Sleep(time.Second * time.Duration(i) * time.Duration(i))
			m.log.Debug().Int("attempt", i).Msg("retrying connection to mongo")
		}
		c, err := mongo.Connect(context.Background(), m.opts)
		if err != nil {
			log.Error().Err(err).Msg("reconnect to mongo failed")
			time.Sleep(time.Second * time.Duration(i) * time.Duration(i))
			continue
		}
		m.conn = c
		return
	}

}

func (m *Mongo) ping() {
	for {
		// ping the server every minute to keep the connection alive
		time.Sleep(60 * time.Second)
		if m.conn == nil {
			continue
		}
		err := m.conn.Ping(context.Background(), nil)
		log.Debug().Err(err).Msg("mongodb ping")
		if err == nil {
			continue
		}

		// if we can't ping the server, try to reconnect
		m.connect()
	}
}

func (m *Mongo) Get(c context.Context, world, kind string, id []string, out interface{}) error {
	pan := log.NewSpan(c, "db.Find", map[string]interface{}{"world": world, "kind": kind, "ids": len(id)})
	defer pan.End()
	return m.objects(c, world_collection(world, kind), bson.M{"_id": bson.M{"$in": id}}, out)
}

func (m *Mongo) List(c context.Context, world, kind string, labels map[string]string, limit, offset int64, out interface{}) error {
	pan := log.NewSpan(c, "db.Find", map[string]interface{}{"world": world, "kind": kind, "labels": len(labels), "limit": limit, "offset": offset})
	defer pan.End()
	return m.listObjects(c, world_collection(world, kind), labels, limit, offset, out)
}

func (m *Mongo) Delete(c context.Context, world, kind string, id string) error {
	pan := log.NewSpan(c, "db.Delete", map[string]interface{}{"world": world, "kind": kind, "id": id})
	defer pan.End()
	return m.deleteObject(c, world_collection(world, kind), id)
}

func (m *Mongo) Set(c context.Context, world, etag string, in []v1.Object) (*Result, error) {
	pan := log.NewSpan(c, "db.Set", map[string]interface{}{"world": world, "etag": etag, "count": len(in)})
	defer pan.End()
	if len(in) == 0 {
		return nil, fmt.Errorf("%w no objects to set", ErrInvalid)
	}
	models, err := prepareSet(world, etag, in)
	if err != nil {
		return nil, err
	}
	return m.setObjects(c, world_collection(world, in[0].GetKind()), models)
}

func (m *Mongo) Close() {
	m.conn.Disconnect(context.Background())
}

func (m *Mongo) collection(name string) *mongo.Collection {
	return m.conn.Database(m.cfg.Database).Collection(name)
}

func (m *Mongo) objects(c context.Context, collection string, filter bson.M, out interface{}) error {
	m.log.Debug().Str("database", m.cfg.Database).Str("collection", collection).Msg("objects")
	cursor, err := m.collection(collection).Find(c, filter)
	if err != nil {
		return err
	}
	return cursor.All(c, out)
}

func (m *Mongo) listObjects(c context.Context, collection string, labels map[string]string, limit, offset int64, out interface{}) error {
	var filter map[string]string
	if labels != nil && len(labels) > 0 {
		filter = map[string]string{}
		for k, v := range labels {
			filter[fmt.Sprintf("Labels.%s", k)] = v
		}
	}

	m.log.Debug().Str("database", m.cfg.Database).Str("collection", collection).Int("limit", int(limit)).Int("offset", int(offset)).Msg("listObjects")
	cursor, err := m.collection(collection).Find(c, filter, &options.FindOptions{
		Limit: &limit,
		Skip:  &offset,
		Sort:  bson.M{"_id": 1},
	})
	if err != nil {
		return err
	}

	return cursor.All(c, out)
}

func (m *Mongo) deleteObject(c context.Context, collection, id string) error {
	m.log.Debug().Str("database", m.cfg.Database).Str("collection", collection).Str("_id", id).Msg("deleteObject")
	_, err := m.collection(collection).DeleteOne(c, bson.M{"_id": id})
	return err
}

// setObjects takes a list of objects and writes them to the database. Handles both insert and update.
// In the worst case(s) we may need to re-read the Etags of docs to determine which were successful.
func (m *Mongo) setObjects(c context.Context, collection string, models []mongo.WriteModel) (*Result, error) {
	m.log.Debug().Str("collection", collection).Int("count", len(models)).Msg("setObjects")

	results, err := m.collection(collection).BulkWrite(c, models, options.BulkWrite().SetOrdered(false))
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrDuplicate
		}
		return nil, err // Timeout / network error
	}

	res := NewResult()
	if int(results.InsertedCount+results.ModifiedCount+results.UpsertedCount) == len(models) {
		return res, nil
	}
	return res, ErrEtagMismatch
}

// prepareSet takes a list of objects and prepares them for insert or update
func prepareSet(world, newEtag string, in []v1.Object) ([]mongo.WriteModel, error) {
	models := []mongo.WriteModel{}
	for _, v := range in {
		isCreate := v.GetEtag() == ""
		if isCreate {
			v.SetEtag(newEtag)
			if v.GetId() == "" {
				v.SetId(uuid.New())
			}
		}

		// seriously mongo I should be able to just set an explicit ID on insert
		// but noooo unless I want to add the field to each struct root I have to
		// on-the-fly another struct to set it there.
		//doc := structs.Map(v)
		//doc["_id"] = v.GetId()
		//doc["etag"] = v.GetEtag()
		doc := v

		var mod mongo.WriteModel
		if isCreate { // if it doesn't have an Etag we're inserting
			mod = mongo.NewInsertOneModel().SetDocument(doc)
		} else { // otherwise we're updating (well, replacing)
			// update if the ID and either the new or the old etag match
			// (ie. if we did a partial update but the etag hasn't changed since our write operation).
			// If we update with the same etag then nothing will change.
			// If we update with the old etag (the one we've sent) then we want the update.
			mod = mongo.NewReplaceOneModel().SetReplacement(doc).SetFilter(bson.D{
				{"_id", v.GetId()},
				{"_etag", bson.M{"$in": []string{v.GetEtag(), newEtag}}},
			})
		}
		v.SetEtag(newEtag)
		models = append(models, mod)
	}
	return models, nil
}

// world_collection returns a name for a world, collection tuple.
// ie. actors_world1, actors_world2 etc.
// This forcibly divides data from different worlds and simplifies data management.
func world_collection(world, name string) string {
	if world == "" {
		// if this doesn't have a world space (ie. it's global or a world itself) then just use the name
		return strings.ToLower(name)
	}
	return strings.ToLower(fmt.Sprintf("%s_%s", name, base36.EncodeBytes([]byte(world))))
}
