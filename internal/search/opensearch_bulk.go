package search

import (
	"bytes"
	"context"
	"fmt"

	"github.com/voidshard/faction/pkg/structs"
	"github.com/voidshard/faction/pkg/util/log"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/opensearch-project/opensearch-go/v4/opensearchutil"
)

type opensearchBulk struct {
	l     log.Logger
	index string
	bulk  opensearchutil.BulkIndexer

	work chan *indexObject
}

type indexObject struct {
	ctx    context.Context
	obj    structs.Object
	result func(error)
}

func newOpensearchBulk(index string, api *opensearchapi.Client, cfg *OpensearchConfig) (*opensearchBulk, error) {
	l := log.Sublogger("opensearch-bulk", map[string]interface{}{"index": index})

	me := &opensearchBulk{
		l:     l,
		index: index,
		work:  make(chan *indexObject),
	}

	bulk, err := opensearchutil.NewBulkIndexer(opensearchutil.BulkIndexerConfig{
		Client:        api,
		FlushInterval: cfg.FlushInterval,
		NumWorkers:    cfg.WriteRoutines,
		OnFlushStart: func(ctx context.Context) context.Context {
			l.Debug().Str("index", index).Msg("opensearch flush start")
			return ctx
		},
		OnFlushEnd: func(ctx context.Context) {
			l.Debug().Str("index", index).Msg("opensearch flush end")
		},
	})
	if err != nil {
		return nil, err
	}

	me.bulk = bulk

	for i := 0; i < cfg.WriteRoutines*2; i++ {
		// start two workers for each opensearch write routine, OS only has to write the
		// buffer out, we have to encode everything from the object and do a bunch of other
		// book keeping allocs
		go me.start()
	}

	return me, nil
}

func (b *opensearchBulk) indexObject(work *indexObject) {
	// https://opensearch.org/docs/latest/api-reference/document-apis/update-document/
	// https://github.com/opensearch-project/opensearch-go/blob/main/opensearchutil/bulk_indexer.go

	b.l.Debug().Str("id", work.obj.GetId()).Str("index", b.index).Msg("indexing object started")
	defer b.l.Debug().Str("id", work.obj.GetId()).Str("index", b.index).Msg("indexing object completed")

	pan := log.NewSpan(work.ctx, "opensearch.bulkindexer", map[string]interface{}{"index": b.index})
	defer func() {
		stats := b.bulk.Stats()
		pan.SetAttributes(map[string]interface{}{
			"added":    stats.NumAdded,
			"flushed":  stats.NumFlushed,
			"failed":   stats.NumFailed,
			"indexed":  stats.NumIndexed,
			"created":  stats.NumCreated,
			"updated":  stats.NumUpdated,
			"deleted":  stats.NumDeleted,
			"requests": stats.NumRequests,
		})
		pan.End()
	}()

	data, err := work.obj.MarshalJSON()
	if err != nil {
		b.l.Error().Err(err).Str("id", work.obj.GetId()).Str("index", b.index).Msg("failed to marshal object")
		pan.Err(err)
		return
	}

	buf := bytes.Buffer{}
	buf.WriteString(`{"doc_as_upsert":true,"doc":`)
	buf.Write(data)
	buf.WriteString("}")

	err = b.bulk.Add(context.Background(), opensearchutil.BulkIndexerItem{
		Action:     "update",
		Index:      b.index,
		DocumentID: work.obj.GetId(),
		Body:       bytes.NewReader(buf.Bytes()),
		OnSuccess: func(ctx context.Context, i opensearchutil.BulkIndexerItem, r opensearchapi.BulkRespItem) {
			if work.result != nil {
				work.result(nil)
			}
		},
		OnFailure: func(ctx context.Context, i opensearchutil.BulkIndexerItem, r opensearchapi.BulkRespItem, err error) {
			errMsg := fmt.Sprintf("error indexing item %s", i.DocumentID)
			if err != nil {
				errMsg = fmt.Sprintf("%s: %s", errMsg, err.Error())
			}
			if r.Error != nil {
				errMsg = fmt.Sprintf(
					"%s %s: %s (cause: %s %s), full %v",
					errMsg, r.Error.Type, r.Error.Reason, r.Error.Cause.Type, r.Error.Cause.Reason, r,
				)
			}
			osErr := fmt.Errorf("failed to index document bulk: %s", errMsg)
			b.l.Error().Err(osErr).Str("id", i.DocumentID).Str("index", i.Index).Msg("failed to index document bulk")
			pan.Err(osErr)
			if work.result != nil {
				work.result(osErr)
			}
		},
	})
	if err != nil {
		b.l.Error().Err(err).Str("id", work.obj.GetId()).Str("index", b.index).Msg("failed to add document to bulk")
		pan.Err(err)
		return
	}
}

func (b *opensearchBulk) Close() {
	close(b.work)
	b.bulk.Close(context.Background())
}

func (b *opensearchBulk) start() {
	b.l.Debug().Msg("starting opensearch bulk indexer worker")
	defer b.l.Debug().Msg("stopping opensearch bulk indexer worker")
	for work := range b.work {
		b.indexObject(work)
	}
}

func (b *opensearchBulk) Index(ctx context.Context, obj structs.Object, result func(error)) {
	b.work <- &indexObject{ctx: ctx, obj: obj, result: result}
}
