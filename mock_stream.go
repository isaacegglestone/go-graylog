package graylog

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

// AddStream adds a stream to the MockServer.
func (ms *MockServer) AddStream(stream *Stream) {
	if stream.Id == "" {
		stream.Id = randStringBytesMaskImprSrc(24)
	}
	ms.Streams[stream.Id] = *stream
	ms.safeSave()
}

// DeleteStream removes a stream from the MockServer.
func (ms *MockServer) DeleteStream(id string) {
	delete(ms.Streams, id)
	ms.safeSave()
}

// StreamList returns a list of all streams.
func (ms *MockServer) StreamList() []Stream {
	if ms.Streams == nil {
		return []Stream{}
	}
	arr := make([]Stream, len(ms.Streams))
	i := 0
	for _, index := range ms.Streams {
		arr[i] = index
		i++
	}
	return arr
}

// EnabledStreamList returns all enabled streams.
func (ms *MockServer) EnabledStreamList() []Stream {
	if ms.Streams == nil {
		return []Stream{}
	}
	arr := []Stream{}
	for _, index := range ms.Streams {
		if index.Disabled {
			continue
		}
		arr = append(arr, index)
	}
	return arr
}

func validateCreateStream(stream *Stream) (int, []byte) {
	key := ""
	switch {
	case stream.Id != "":
		key = "id"
	case stream.CreatorUserId != "":
		key = "creator_user_id"
	case stream.Outputs != nil && len(stream.Outputs) != 0:
		key = "outputs"
	case stream.CreatedAt != "":
		key = "created_at"
	case stream.Disabled:
		key = "disabled"
	case stream.AlertConditions != nil && len(stream.AlertConditions) != 0:
		key = "alert_conditions"
	case stream.AlertReceivers != nil:
		key = "alert_receivers"
	case stream.IsDefault:
		key = "is_default"
	}
	if key != "" {
		return 400, []byte(fmt.Sprintf(`{"type": "ApiError", "message": "Unable to map property %s.\nKnown properties include: index_set_id, rules, title, description, content_pack, matching_type, remove_matches_from_default_stream"}`, key))
	}
	if stream.Title == "" {
		return 400, []byte(`{"type": "ApiError", "message": "Can not construct instance of org.graylog2.rest.resources.streams.requests.CreateStreamRequest, problem: Null title\n at [Source: org.glassfish.jersey.message.internal.ReaderInterceptorExecutor$UnCloseableInputStream@53a6a093; line: 1, column: 2]" }`)
	}
	if stream.IndexSetId == "" {
		return 400, []byte(`{"type": "ApiError", "message": "Can not construct instance of org.graylog2.rest.resources.streams.requests.CreateStreamRequest, problem: Null indexSetId\n at [Source: org.glassfish.jersey.message.internal.ReaderInterceptorExecutor$UnCloseableInputStream@3b7194f4; line: 1, column: 17]"}`)
	}
	// 500, {"type": "ApiError", "message": "invalid hexadecimal representation of an ObjectId: [%s]"}

	return 200, []byte("")
}

// GET /streams Get a list of all streams
func (ms *MockServer) handleGetStreams(
	w http.ResponseWriter, r *http.Request, _ httprouter.Params,
) {
	ms.handleInit(w, r, false)
	arr := ms.StreamList()
	streams := &streamsBody{Streams: arr, Total: len(arr)}
	writeOr500Error(w, streams)
}

// POST /streams Create index set
func (ms *MockServer) handleCreateStream(
	w http.ResponseWriter, r *http.Request, _ httprouter.Params,
) {
	b, err := ms.handleInit(w, r, true)
	if err != nil {
		write500Error(w)
		return
	}

	requiredFields := []string{"title", "index_set_id"}
	allowedFields := []string{
		"title", "index_set_id", "rules", "description", "content_pack",
		"matching_type", "remove_matches_from_default_stream"}
	sc, msg, body := validateRequestBody(b, requiredFields, allowedFields, nil)
	if sc != 200 {
		w.WriteHeader(sc)
		w.Write([]byte(msg))
		return
	}

	stream := &Stream{}
	if err := mapstructure.Decode(body, stream); err != nil {
		ms.Logger.WithFields(log.Fields{
			"body": string(b), "error": err,
		}).Info("Failed to parse request body as stream")
		writeApiError(w, 400, "400 Bad Request")
		return
	}

	if err := CreateValidator.Struct(stream); err != nil {
		writeApiError(w, 400, err.Error())
		return
	}
	ms.AddStream(stream)
	ret := map[string]string{"stream_id": stream.Id}
	writeOr500Error(w, ret)
}

// GET /streams/enabled Get a list of all enabled streams
func (ms *MockServer) handleGetEnabledStreams(
	w http.ResponseWriter, r *http.Request, _ httprouter.Params,
) {
	ms.handleInit(w, r, false)
	arr := ms.EnabledStreamList()
	streams := &streamsBody{Streams: arr, Total: len(arr)}
	writeOr500Error(w, streams)
}

// GET /streams/{streamId} Get a single stream
func (ms *MockServer) handleGetStream(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params,
) {
	id := ps.ByName("streamId")
	if id == "enabled" {
		ms.handleGetEnabledStreams(w, r, ps)
		return
	}
	ms.handleInit(w, r, false)
	stream, ok := ms.Streams[id]
	if !ok {
		writeApiError(w, 404, "No stream found with id %s", id)
		return
	}
	writeOr500Error(w, &stream)
}

// PUT /streams/{streamId} Update a stream
func (ms *MockServer) handleUpdateStream(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params,
) {
	b, err := ms.handleInit(w, r, true)
	if err != nil {
		write500Error(w)
		return
	}
	id := ps.ByName("streamId")
	stream, ok := ms.Streams[id]
	if !ok {
		writeApiError(w, 404, "No stream found with id %s", id)
		return
	}
	data := map[string]interface{}{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		writeApiError(w, 400, "400 Bad Request")
		return
	}
	if title, ok := data["title"]; ok {
		t, ok := title.(string)
		if !ok {
			writeApiError(w, 400, "title must be string")
			return
		}
		stream.Title = t
	}
	if description, ok := data["description"]; ok {
		d, ok := description.(string)
		if !ok {
			writeApiError(w, 400, "description must be string")
			return
		}
		stream.Description = d
	}
	// TODO outputs
	if matchingType, ok := data["matching_type"]; ok {
		m, ok := matchingType.(string)
		if !ok {
			writeApiError(w, 400, "matching_type must be string")
			return
		}
		stream.MatchingType = m
	}
	if removeMathcesFromDefaultStream, ok := data["remove_matches_from_default_stream"]; ok {
		m, ok := removeMathcesFromDefaultStream.(bool)
		if !ok {
			writeApiError(w, 400, "remove_matches_from_default_stream must be bool")
			return
		}
		stream.RemoveMatchesFromDefaultStream = m
	}
	if indexSetId, ok := data["index_set_id"]; ok {
		m, ok := indexSetId.(string)
		if !ok {
			writeApiError(w, 400, "index_set_id must be string")
			return
		}
		stream.IndexSetId = m
	}
	stream.Id = id
	ms.AddStream(&stream)
	writeOr500Error(w, &stream)
}

// DELETE /streams/{streamId} Delete a stream
func (ms *MockServer) handleDeleteStream(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params,
) {
	ms.handleInit(w, r, false)
	id := ps.ByName("streamId")
	_, ok := ms.Streams[id]
	if !ok {
		writeApiError(w, 404, "No stream found with id %s", id)
		return
	}
	ms.DeleteStream(id)
}

// POST /streams/{streamId}/pause Pause a stream
func (ms *MockServer) handlePauseStream(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params,
) {
	ms.handleInit(w, r, false)
	id := ps.ByName("streamId")
	_, ok := ms.Streams[id]
	if !ok {
		writeApiError(w, 404, "No stream found with id %s", id)
		return
	}
	// TODO pause
}

func (ms *MockServer) handleResumeStream(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params,
) {
	ms.handleInit(w, r, false)
	id := ps.ByName("streamId")
	_, ok := ms.Streams[id]
	if !ok {
		writeApiError(w, 404, "No stream found with id %s", id)
		return
	}
	// TODO resume
}
