package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/queue"
	"github.com/voidshard/faction/internal/search"
	"github.com/voidshard/faction/pkg/kind"
	"github.com/voidshard/faction/pkg/structs/api"
	"github.com/voidshard/faction/pkg/util/log"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	apiVersion = "v1"
)

// Server encapsulates serving our service as an HTTP API.
type Server struct {
	cfg *Config
	log log.Logger

	srv    *http.Server
	router *mux.Router

	shuttingDown bool

	svc *Service
}

func NewServer(cfg *Config, db db.Database, qu queue.Queue, sb search.Search) (*Server, error) {
	if cfg == nil {
		cfg = &Config{}
	}
	cfg.setDefaults()

	svc, err := newService(cfg, db, qu, sb)
	if err != nil {
		return nil, err
	}

	me := &Server{
		cfg:    cfg,
		log:    log.Sublogger("api.server"),
		router: mux.NewRouter(),
		svc:    svc,
	}

	me.router.HandleFunc(fmt.Sprintf("/_health"), me.health).Methods("GET")
	me.router.HandleFunc(fmt.Sprintf("/%s/event", apiVersion), me.deferEvent).Methods("POST")
	me.router.HandleFunc(fmt.Sprintf("/%s/event", apiVersion), me.onChangeEvent).Methods("GET") // Websocket
	me.router.HandleFunc(fmt.Sprintf("/%s/{world}/search", apiVersion), me.search).Methods("GET")
	me.router.HandleFunc(fmt.Sprintf("/%s/{kind}", apiVersion), me.getKind).Methods("GET")
	me.router.HandleFunc(fmt.Sprintf("/%s/{kind}", apiVersion), me.setKind).Methods("POST")
	me.router.HandleFunc(fmt.Sprintf("/%s/{kind}", apiVersion), me.delKind).Methods("DELETE")

	return me, nil
}

func (s *Server) writeResp(w http.ResponseWriter, code int, resp interface{}) {
	if resp == nil {
		resp = struct {
			Error *api.ErrorResponse
		}{
			Error: &api.ErrorResponse{},
		}
	}

	data, err := json.Marshal(resp)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)

	if code >= 500 {
		log.Error().Int("code", code).Msg(string(data))
	} else if code >= 400 {
		log.Warn().Int("code", code).Msg(string(data))
	}
}

func (s *Server) onChangeEvent(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), s.cfg.TimeoutRead)
	defer cancel()

	pan := log.NewSpan(ctx, "api.onChangeEvent", map[string]interface{}{"url": r.URL.String()})
	ctx = pan.Context // make sure the root span is in the context
	defer pan.End()

	resp := &api.ErrorResponse{}

	// parse request from URL query params
	qvars := r.URL.Query()
	req := &api.StreamEvents{
		World:      qvars.Get("world"),
		Kind:       qvars.Get("kind"),
		Id:         qvars.Get("id"),
		Controller: qvars.Get("controller"),
		Queue:      qvars.Get("queue"),
	}

	// upgrade to a websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		pan.Err(err)
		resp.Code = http.StatusInternalServerError
		resp.Message = "failed to upgrade connection to websocket"
		s.writeResp(w, http.StatusInternalServerError, resp)
		return
	}
	sock := newWebSocket(s.svc, conn)

	// usual validations
	if req.Kind != "" && !kind.IsValid(req.Kind) {
		pan.Err(fmt.Errorf("kind %s not found", req.Kind))
		resp.Code = http.StatusBadRequest
		resp.Message = "invalid kind"
		sock.Close(resp)
		return
	}
	err = kind.Validate(req.Kind, req)
	if err != nil {
		pan.Err(err)
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		sock.Close(resp)
		return
	}

	// subscribe to events
	events, kill, err := s.svc.subscribeToEvents(req)
	if err != nil {
		pan.Err(err)
		resp.Code = http.StatusInternalServerError
		resp.Message = "failed to subscribe to events"
		sock.Close(resp)
		return
	}

	// start the pump to begin normal operation
	sock.Pump(events, kill)
	return
}

func (s *Server) deferEvent(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), s.cfg.TimeoutRead)
	defer cancel()

	pan := log.NewSpan(ctx, "api.DeferEvent")
	ctx = pan.Context // make sure the root span is in the context
	defer pan.End()

	resp := &api.DeferEventResponse{Error: &api.ErrorResponse{}}
	req := &api.DeferEventRequest{}

	err := readJson(r, req)
	if err != nil {
		pan.Err(err)
		resp.Error.Code = http.StatusBadRequest
		resp.Error.Message = "invalid request json"
		s.writeResp(w, http.StatusBadRequest, resp)
		return
	}

	if !kind.IsValid(req.Kind) {
		err = fmt.Errorf("kind %s not found", req.Kind)
		pan.Err(err)
		resp.Error.Code = http.StatusBadRequest
		resp.Error.Message = "invalid kind"
		s.writeResp(w, http.StatusBadRequest, resp)
		return
	}

	pan.SetAttributes(map[string]interface{}{"kind": req.Kind, "world": req.World, "controller": req.Controller, "id": req.Id, "to-tick": req.ToTick, "by-tick": req.ByTick})

	err = kind.Validate(req.Kind, req)
	if err != nil {
		pan.Err(err)
		resp.Error.Code = http.StatusBadRequest
		resp.Error.Message = err.Error()
		s.writeResp(w, http.StatusBadRequest, resp)
		return
	}

	err = s.svc.deferEvent(ctx, req, resp)
	if err != nil {
		pan.Err(err)
		resp.Error.Code = errorCodeHTTP(err)
		resp.Error.Message = err.Error()
		s.writeResp(w, errorCodeHTTP(err), resp)
		return
	}

	s.writeResp(w, http.StatusOK, resp)
	return
}

func (s *Server) search(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), s.cfg.TimeoutRead)
	defer cancel()

	pan := log.NewSpan(ctx, "api.Search")
	ctx = pan.Context // make sure the root span is in the context
	defer pan.End()

	req := &api.SearchRequest{}
	resp := &api.SearchResponse{Error: &api.ErrorResponse{}}

	err := readJson(r, req)
	if err != nil {
		pan.Err(err)
		resp.Error.Code = http.StatusBadRequest
		resp.Error.Message = "invalid request json"
		s.writeResp(w, http.StatusBadRequest, resp)
		return
	}

	if !kind.IsValid(req.Kind) {
		pan.Err(fmt.Errorf("kind %s not found", req.Kind))
		resp.Error.Code = http.StatusNotFound
		resp.Error.Message = "invalid kind"
		s.writeResp(w, http.StatusNotFound, resp)
		return
	}

	vars := mux.Vars(r)
	world, ok := vars["world"]
	if !ok || world == "" {
		pan.Err(fmt.Errorf("world id invalid"))
		resp.Error.Code = http.StatusBadRequest
		resp.Error.Message = "invalid world"
		s.writeResp(w, http.StatusBadRequest, resp)
		return
	}

	err = kind.Validate(req.Kind, req)
	if err != nil {
		pan.Err(err)
		resp.Error.Code = http.StatusBadRequest
		resp.Error.Message = err.Error()
		s.writeResp(w, http.StatusBadRequest, resp)
		return
	}

	pan.SetAttributes(map[string]interface{}{
		"kind":          req.Kind,
		"world":         world,
		"limit":         req.Limit,
		"random-weight": req.RandomWeight,
		"all":           len(req.All),
		"any":           len(req.Any),
		"not":           len(req.Not),
		"score":         len(req.Score),
	})

	err = s.svc.searchKind(ctx, world, req, resp)
	if err != nil {
		pan.Err(err)
		resp.Error.Code = errorCodeHTTP(err)
		resp.Error.Message = err.Error()
		s.writeResp(w, errorCodeHTTP(err), resp)
		return
	}

	s.writeResp(w, http.StatusOK, resp)
	return
}

// nb. technically dictating our reply based on a GET body is considered anti-html best practice
// but opensearch / elasticsearch do this because it makes more sense than forcing users to use
// a POST to get data OR forcing a boatload of query params .. so .. eh.
func (s *Server) getKind(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), s.cfg.TimeoutRead)
	defer cancel()

	pan := log.NewSpan(ctx, "api.Get")
	ctx = pan.Context // make sure the root span is in the context
	defer pan.End()

	resp := &api.GetResponse{Error: &api.ErrorResponse{}}

	// read out the kind
	vars := mux.Vars(r)
	k, ok := vars["kind"]
	if !ok || !kind.IsValid(k) {
		pan.Err(fmt.Errorf("kind %s not found", k))
		resp.Error.Code = http.StatusNotFound
		resp.Error.Message = "not found"
		s.writeResp(w, http.StatusNotFound, resp)
		return
	}
	pan.SetAttributes(map[string]interface{}{"kind": k})

	// parse the body of the request
	body := &api.GetRequest{}
	err := readJson(r, body)
	if err != nil {
		pan.Err(err)
		resp.Error.Code = http.StatusBadRequest
		resp.Error.Message = "invalid request json"
		s.writeResp(w, http.StatusBadRequest, resp)
		return
	}

	// we could just ignore this, but it might confuse a caller if we ignore inputs
	if len(body.Ids) > 0 && len(body.Labels) > 0 {
		pan.Err(fmt.Errorf("cannot specify both IDs and labels in get"))
		resp.Error.Code = http.StatusBadRequest
		resp.Error.Message = "cannot specify both IDs and labels in get"
		s.writeResp(w, http.StatusBadRequest, resp)
		return
	}
	if body.Limit <= 0 {
		body.Limit = 100 // set default if user doesn't specify
	}

	// set some attributes for the span
	pan.SetAttributes(map[string]interface{}{
		"ids":    len(body.Ids),
		"limit":  body.Limit,
		"offset": body.Offset,
		"labels": len(body.Labels),
		"world":  body.World,
	})

	// validate the GetRequest Fields
	err = kind.Validate(k, body)
	if err != nil {
		pan.Err(err)
		resp.Error.Code = http.StatusBadRequest
		resp.Error.Message = err.Error()
		s.writeResp(w, http.StatusBadRequest, resp)
		return
	}

	err = s.svc.getKind(ctx, k, body, resp)
	if err != nil {
		pan.Err(err)
		resp.Error.Code = errorCodeHTTP(err)
		resp.Error.Message = err.Error()
		s.writeResp(w, errorCodeHTTP(err), resp)
		return
	}

	s.writeResp(w, http.StatusOK, resp)
	return
}

func (s *Server) setKind(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), s.cfg.TimeoutWrite)
	defer cancel()

	pan := log.NewSpan(ctx, "api.setKind")
	ctx = pan.Context // make sure the root span is in the context
	defer pan.End()

	resp := &api.SetResponse{Error: &api.ErrorResponse{}}

	// read out the kind
	vars := mux.Vars(r)
	k, ok := vars["kind"]
	if !ok || !kind.IsValid(k) {
		pan.Err(fmt.Errorf("kind %s not found", k))
		resp.Error.Code = http.StatusNotFound
		resp.Error.Message = "not found"
		s.writeResp(w, http.StatusNotFound, resp)
		return
	}
	pan.SetAttributes(map[string]interface{}{"kind": k})

	// parse the body of the request
	body := &api.SetRequest{}
	err := readJson(r, body)
	if err != nil {
		pan.Err(err)
		resp.Error.Code = http.StatusBadRequest
		resp.Error.Message = "invalid request json"
		s.writeResp(w, http.StatusBadRequest, resp)
		return
	}

	// validate request Fields
	err = kind.Validate(k, body)
	if err != nil {
		pan.Err(err)
		resp.Error.Code = http.StatusBadRequest
		resp.Error.Message = err.Error()
		s.writeResp(w, http.StatusBadRequest, resp)
		return
	}

	err = s.svc.setKind(ctx, k, body, resp)
	if err != nil {
		pan.Err(err)
		resp.Error.Code = errorCodeHTTP(err)
		resp.Error.Message = err.Error()
		s.writeResp(w, errorCodeHTTP(err), resp)
		return
	}

	s.writeResp(w, http.StatusOK, resp)
	return
}

func (s *Server) delKind(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), s.cfg.TimeoutWrite)
	defer cancel()

	pan := log.NewSpan(ctx, "api.delKind")
	ctx = pan.Context // make sure the root span is in the context
	defer pan.End()

	resp := &api.DeleteResponse{Error: &api.ErrorResponse{}}

	// read out the kind
	vars := mux.Vars(r)
	k, ok := vars["kind"]
	if !ok || !kind.IsValid(k) {
		pan.Err(fmt.Errorf("kind %s not found", k))
		resp.Error.Code = http.StatusNotFound
		resp.Error.Message = "not found"
		s.writeResp(w, http.StatusNotFound, resp)
		return
	}
	pan.SetAttributes(map[string]interface{}{"kind": k})

	// parse the body of the request
	body := &api.DeleteRequest{}
	err := readJson(r, body)
	if err != nil {
		pan.Err(err)
		resp.Error.Code = http.StatusBadRequest
		resp.Error.Message = "invalid request json"
		s.writeResp(w, http.StatusBadRequest, resp)
		return
	}

	// validate request Fields
	err = kind.Validate(k, body)
	if err != nil {
		pan.Err(err)
		resp.Error.Code = http.StatusBadRequest
		resp.Error.Message = err.Error()
		s.writeResp(w, http.StatusBadRequest, resp)
		return
	}

	pan.SetAttributes(map[string]interface{}{"ids": len(body.Ids), "world": body.World})

	err = s.svc.deleteKind(ctx, k, body, resp)
	if err != nil {
		pan.Err(err)
		resp.Error.Code = errorCodeHTTP(err)
		resp.Error.Message = err.Error()
		s.writeResp(w, errorCodeHTTP(err), resp)
		return
	}

	s.writeResp(w, http.StatusOK, resp)
	return
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	if s.shuttingDown {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) Stop() {
	s.shuttingDown = true
	s.svc.Shutdown()
	s.srv.Shutdown(context.Background())
}

func (s *Server) Serve(port int) error {
	s.srv = &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      handlers.CompressHandler(s.router),
		ReadTimeout:  s.cfg.TimeoutRead,
		WriteTimeout: s.cfg.TimeoutWrite,
	}
	return s.srv.ListenAndServe()
}
