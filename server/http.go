package main

import (
	"encoding/json"
	"net/http"
)

type ApiServer struct {
}

func (s *ApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Info("--> %s %s\n", r.Method, r.URL.Path)

	switch {
	default:
		s.handleNotFound(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/clients":
		s.serverClients(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/server_status":
		s.serveServerStatus(w, r)
	}
}

func (s *ApiServer) serveServerStatus(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"number_of_clients":              len(clients),
		"number_of_operators":            -1,
		"number_of_client_connections":   -1,
		"number_of_operator_connections": -1,
	}

	jsonString, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

func (s *ApiServer) serverClients(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{}

	for key, value := range clients {
		cl := map[string]interface{}{
			"id":                  value.cl.id,
			"client_connected":    value.cl.connected,
			"client_localAddr":    value.cl.conn.LocalAddr(),
			"client_remoteAddr":   value.cl.conn.RemoteAddr(),
			"operators_connected": len(value.op),
			"open_streams":        value.cl.session.NumStreams(),
			"closed":              value.cl.session.IsClosed(),
		}
		var operators = make([]map[string]interface{}, len(value.op))
		for idx, o := range value.op {
			operators[idx] = map[string]interface{}{
				"localAddr":    o.conn.LocalAddr(),
				"remoteAddr":   o.conn.RemoteAddr(),
				"open_streams": o.session.NumStreams(),
				"closed":       o.session.IsClosed(),
			}
		}
		cl["operators"] = operators
		response[key] = cl
	}

	jsonString, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

func (s *ApiServer) handleNotFound(w http.ResponseWriter, r *http.Request) {
	log.Debug("handleNotFound\n")
	w.WriteHeader(http.StatusNotFound)
}

func (s *ApiServer) handleError(w http.ResponseWriter, r *http.Request, err error) {
	log.Debug("handleError\n")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
