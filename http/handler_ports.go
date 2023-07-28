package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/allegro/bigcache/v3"
	"github.com/fir1/port/internal/port/model"
)

// savePorts example
//
//	@Summary		It will create ports and save into the DB by default ports.json file will be used
//	@Description	It will create ports and save into the DB by default ports.json file will be used
//	@Tags Ports
//	@ID				save-ports
//	@Accept			json
//	@Produce		json
//
// @Success      201
//
//	@Failure      400
//
// @Failure      500
// @Router			/ports [post].
func (s *Service) savePorts(w http.ResponseWriter, r *http.Request) {
	err := s.portService.SavePortsFromFile(r.Context(), "ports.json", nil)
	if err != nil {
		s.respond(w, err, http.StatusInternalServerError)
		return
	}

	// we have updated list on DB so we have to clear cache
	// so our API's must refetch the list
	err = s.cacheClient.Reset()
	if err != nil {
		s.respond(w, err, http.StatusInternalServerError)
		return
	}

	s.respond(w, nil, http.StatusCreated)
}

// savePortsFromFile example
//
//	@Summary		You are able to provide json file the service will parse and save into the DB
//	@Description	You are able to provide json file the service will parse and save into the DB
//	@Tags Ports
//	@ID				save-ports-from-file
//	@Accept			mpfd
//	@Produce		json
//
// @Param file formData file true "File"
// @Success      201
//
//	@Failure      400
//
// @Failure      500
// @Router			/ports/from-file [post].
func (s *Service) savePortsFromFile(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file not found", http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = s.portService.SavePortsFromFile(r.Context(), "", file)
	if err != nil {
		s.respond(w, err, http.StatusInternalServerError)
		return
	}

	// we have updated list on DB so we have to clear cache
	// so our API's must refetch the list
	err = s.cacheClient.Reset()
	if err != nil {
		s.respond(w, err, http.StatusInternalServerError)
		return
	}

	s.respond(w, nil, http.StatusCreated)
}

// listPorts example
//
//		@Summary		It will return all the available ports from the DB
//		@Description	It will return all the available ports from the DB. We will use API caching for this purpose
//	 	@Description so we don't have to get all data over again from DB, which is useful in real world applications
//	 	@Description where we are connected to the real database such as PostgresSQL it saves a lot of latency.
//		@Description in case if the list empty, it means that you should call API endpoint `POST /ports` it
//		@Description will parse `ports.json` file and saves into the DB, then you can make a call to `GET /ports`
//		@Description to get all the available ports from the DB
//		@Tags Ports
//		@ID			list-ports
//		@Accept			json
//		@Produce		json
//
// @Success      201
//
//	@Failure      400
//
// @Failure      500
// @Router			/ports [get].
func (s *Service) listPorts(w http.ResponseWriter, r *http.Request) {
	cacheResponse, err := s.cacheClient.Get(r.RequestURI)
	switch {
	case err == nil:
		response := map[string]model.Port{}
		err = json.Unmarshal(cacheResponse, &response)
		if err != nil {
			s.respond(w, err, http.StatusInternalServerError)
			return
		}
		s.respond(w, response, http.StatusOK)
		return
	case errors.Is(err, bigcache.ErrEntryNotFound):
	default:
		s.respond(w, err, http.StatusInternalServerError)
		return
	}

	ports, err := s.portService.ListPorts(r.Context())
	if err != nil {
		s.respond(w, err, http.StatusInternalServerError)
		return
	}

	responseBytes, err := json.Marshal(&ports)
	if err != nil {
		s.respond(w, err, http.StatusInternalServerError)
		return
	}

	err = s.cacheClient.Set(r.RequestURI, responseBytes)
	if err != nil {
		s.respond(w, err, http.StatusInternalServerError)
		return
	}

	s.respond(w, ports, http.StatusOK)
}
