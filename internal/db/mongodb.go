package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/martinlindhe/base36"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/voidshard/faction/internal/log"
	"github.com/voidshard/faction/internal/uuid"
	"github.com/voidshard/faction/pkg/structs"
)

const (
	colWorlds   = "worlds"
	colActors   = "actors"
	colFactions = "factions"
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

	opts := options.Client().ApplyURI(url)
	// github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo/example_test.go
	opts.Monitor = otelmongo.NewMonitor()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	me := &Mongo{
		opts: opts,
		log: log.Sublogger("mongodb", map[string]string{
			"host":     cfg.Host,
			"port":     fmt.Sprintf("%d", cfg.Port),
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

func (m *Mongo) Worlds(c context.Context, id []string) ([]*structs.World, error) {
	var results []*structs.World
	return results, m.objects(c, colWorlds, bson.M{"_id": bson.M{"$in": id}}, &results)
}

func (m *Mongo) ListWorlds(c context.Context, labels map[string]string, limit, offset int64) ([]*structs.World, error) {
	var results []*structs.World
	return results, m.listObjects(c, colWorlds, labels, limit, offset, &results)
}

func (m *Mongo) SetWorld(c context.Context, etag string, in *structs.World) (string, error) {
	m.log.Debug().Str("database", m.cfg.Database).Str("collection", colWorlds).Str("_id", in.Id).Str("etag", in.Etag).Msg("setWorld")
	if in.Name == "" {
		return "", fmt.Errorf("%w world name required", ErrInvalid)
	}

	// if we don't have an ID we'll generate one
	// if a name is set we'll use that for a deterministic ID
	if in.Id == "" {
		if in.Name == "" { // random ID
			in.Id = uuid.NewID().String()
		} else { // deterministic ID
			in.Id = uuid.NewID(in.Name).String()
		}
	} else if !uuid.IsValidUUID(in.Id) {
		return "", fmt.Errorf("%w invalid id: %s", ErrInvalid, in.Id)
	}

	filter := bson.D{{"_id", in.Id}, {"etag", bson.M{"$in": []string{in.Etag, etag}}}}
	in.Etag = etag

	result, err := m.collection(colWorlds).UpdateOne(
		c,
		filter,
		bson.M{"$set": in},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return "", err
	}
	if result.UpsertedCount+result.ModifiedCount == 0 {
		return "", ErrEtagMismatch
	}

	return etag, err
}

func (m *Mongo) DeleteWorld(c context.Context, id string) error {
	return m.deleteObject(c, colWorlds, id)
}

func (m *Mongo) Actors(c context.Context, world string, id []string) ([]*structs.Actor, error) {
	var results []*structs.Actor
	return results, m.objects(c, world_collection(world, colActors), bson.M{"_id": bson.M{"$in": id}}, &results)
}

func (m *Mongo) ListActors(c context.Context, world string, labels map[string]string, limit, offset int64) ([]*structs.Actor, error) {
	var results []*structs.Actor
	return results, m.listObjects(c, world_collection(world, colActors), labels, limit, offset, &results)
}

func (m *Mongo) SetActors(c context.Context, world, etag string, in []*structs.Actor) (*Result, error) {
	found, err := m.Worlds(c, []string{world})
	if err != nil {
		return nil, err
	}
	if len(found) == 0 {
		return nil, fmt.Errorf("%w world not found: %s", ErrNotFound, world)
	}
	models, err := prepareSet(world, etag, in)
	if err != nil {
		return nil, err
	}
	return m.setObjects(c, world_collection(world, colActors), models)
}

func (m *Mongo) DeleteActor(c context.Context, world, id string) error {
	return m.deleteObject(c, world_collection(world, colActors), id)
}

func (m *Mongo) Factions(c context.Context, world string, id []string) ([]*structs.Faction, error) {
	var results []*structs.Faction
	return results, m.objects(c, world_collection(world, colFactions), bson.M{"_id": bson.M{"$in": id}}, &results)
}

func (m *Mongo) ListFactions(c context.Context, world string, labels map[string]string, limit, offset int64) ([]*structs.Faction, error) {
	var results []*structs.Faction
	return results, m.listObjects(c, world_collection(world, colFactions), labels, limit, offset, &results)
}

func (m *Mongo) SetFactions(c context.Context, world, etag string, in []*structs.Faction) (*Result, error) {
	found, err := m.Worlds(c, []string{world})
	if err != nil {
		return nil, err
	}
	if len(found) == 0 {
		return nil, fmt.Errorf("%w world not found: %s", ErrNotFound, world)
	}
	models, err := prepareSet(world, etag, in)
	if err != nil {
		return nil, err
	}
	return m.setObjects(c, world_collection(world, colFactions), models)
}

func (m *Mongo) DeleteFaction(c context.Context, world, id string) error {
	return m.deleteObject(c, world_collection(world, colFactions), id)
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
			filter[fmt.Sprintf("labels.%s", k)] = v
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
		return nil, err
	}

	res := NewResult()
	if int(results.InsertedCount+results.ModifiedCount+results.UpsertedCount) == len(models) {
		return res, nil
	}
	return res, ErrEtagMismatch
}

// prepareSet takes a list of objects and prepares them for insert or update
func prepareSet[o structs.Object](world, newEtag string, in []o) ([]mongo.WriteModel, error) {
	models := []mongo.WriteModel{}
	for _, v := range in {
		var mod mongo.WriteModel
		if v.GetId() == "" { // if it doesn't have an id we're inserting
			v.SetId(uuid.NewID().String())
			mod = mongo.NewInsertOneModel().SetDocument(v)
		} else { // otherwise we're updating (well, replacing)
			if !uuid.IsValidUUID(v.GetId()) {
				return nil, fmt.Errorf("%w invalid id: %s", ErrInvalid, v.GetId())
			}
			// update if the ID and either the new or the old etag match
			// (ie. if we did a partial update but the etag hasn't changed since our write operation).
			// If we update with the same etag then nothing will change.
			// If we update with the old etag (the one we've sent) then we want the update.
			mod = mongo.NewReplaceOneModel().SetReplacement(v).SetFilter(bson.D{
				{"_id", v.GetId()},
				{"etag", bson.M{"$in": []string{v.GetEtag(), newEtag}}},
			})
		}
		v.SetEtag(newEtag)
		v.SetWorld(world)
		models = append(models, mod)
	}
	return models, nil
}

// world_collection returns a name for a world, collection tuple.
// ie. actors_world1, actors_world2 etc.
// This forcibly divides data from different worlds and simplifies data management.
func world_collection(world, name string) string {
	return fmt.Sprintf("%s_%s", name, strings.ToLower(base36.EncodeBytes([]byte(world))))
}
