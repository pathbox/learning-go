package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
)

type DialHandler struct {
	*httprouter.Router
	DialService wtf.DialService
	Logger      *log.Logger
}

func NewDialHandler() *DialHandler {
	h := &DialHandler{
		Router: httprouter.New(),
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}

	h.POST("/api/dials", h.handlePostDial)
	h.GET("/api/dials/:id", h.handleGetDial)
	h.PATCH("/api/dials/:id", h.handlePatchDial)
	return h

}

// handleGetDial handles requests to create a new dial.
func (h *DialHandler) handlePostDial(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Decode request

	var req postDialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, ErrinvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}

	d := req.Dial
	d.Token = req.Token
	d.ModTime = time.Time{}

	// Create dial.
	switch err := h.DialService.CreateDial(d); err {
	case nil:
		encodeJSON(w, &postDialResponse{Dial: d}, h.Logger)
	case wtf.ErrDialRequired, wtf.ErrDialIDRequired:
		Error(w, err, http.StatusBadRequest, h.Logger)
	case wtf.ErrDialExists:
		Error(w, err, http.StatusConflict, h.Logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

type postDialRequest struct {
	Dial  *wtf.Dial `json:"dial,omitempty"`
	Token string    `json:"token,omitempty"`
}

type postDialResponse struct {
	Dial *wtf.Dial `json:"dial,omitempty"`
	Err  string    `json:"err,omitempty"`
}

func (h *DialHandler) handleGetDial(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	id := ps.ByName("id")

	d, err := h.DialService.Dial(wtf.DialID(id))
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.Logger)
	} else if d == nil {
		NotFound(w)
	} else {
		encodeJSON(w, &getDialResponse{Dial: d}, h.Logger)
	}
}

type getDialResponse struct {
	Dial *wtf.Dial `json:"dial,omitempty"`
	Err  string    `json:"err,omitempty"`
}

func (h *DialHandler) handlePatchDial(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var req patchDialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}
	// Create dial.
	switch err := h.DialService.SetLevel(req.ID, req.Token, req.Level); err {
	case nil:
		encodeJSON(w, &patchDialResponse{}, h.Logger)
	case wtf.ErrDialNotFound:
		Error(w, err, http.StatusNotFound, h.Logger)
	case wtf.ErrUnauthorized:
		Error(w, err, http.StatusUnauthorized, h.Logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

type patchDialRequest struct {
	ID    wtf.DialID `json:"id"`
	Token string     `json:"token"`
	Level float64    `json:"level"`
}

type patchDialResponse struct {
	Err string `json:"err,omitempty"`
}

func Error(w http.ResponseWriter, err error, code int, logger *log.Logger) {
	logger.Printf("http error: %s (code=%d)", err, code)

	if code == http.StatusInternalServerError {
		err = wtf.ErrInternal
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&errorResponse{Err: err.Error()})
}

// errorResponse is a generic response for sending a error.
type errorResponse struct {
	Err string `json:"err,omitempty"`
}

func NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{}` + "\n"))
}

// encodeJSON encodes v to w in JSON format. Error() is called if encoding fails.
func encodeJSON(w http.ResponseWriter, v interface{}, logger *log.Logger) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		Error(w, err, http.StatusInternalServerError, logger)
	}
}

type Client struct {
	URL         url.URL
	dialService DialService
}

func NewClient() *Client {
	c := &Client{}
	c.dialService.URL = &c.URL
	return c
}

// DialService returns the service for managing dials.
func (c *Client) DialService() wtf.DialService {
	return &c.dialService
}

// DialService represents an HTTP implementation of wtf.DialService.
type DialService struct {
	URL *url.URL
}

/ CreateDial creates a new dial.
func (s *DialService) CreateDial(d *wtf.Dial) error {
  // Validate arguments.
  if d == nil {
    return wtf.ErrDialRequired
  }

  u := *s.URL
  u.Path = "/api/dials"

  // Save token.
  token := d.Token

  // Encode request body.
  reqBody, err := json.Marshal(postDialRequest{Dial: d, Token: token})
  if err != nil {
    return err
  }

  // Execute request.
  resp, err := http.Post(u.String(), "application/json", bytes.NewReader(reqBody))
  if err != nil {
    return err
  }
  defer resp.Body.Close()

  // Decode response into JSON.
  var respBody postDialResponse
  if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
    return err
  } else if respBody.Err != "" {
    return wtf.Error(respBody.Err)
  }

  // Copy returned dial.
  *d = *respBody.Dial
  d.Token = token

  return nil
}

// Dial returns a dial by id.
func (s *DialService) Dial(id wtf.DialID) (*wtf.Dial, error) {
  u := *s.URL
  u.Path = "/api/dials/" + url.QueryEscape(string(id))

  // Execute request.
  resp, err := http.Get(u.String())
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()

  // Decode response into JSON.
  var respBody getDialResponse
  if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
    return nil, err
  } else if respBody.Err != "" {
    return nil, wtf.Error(respBody.Err)
  }
  return respBody.Dial, nil
}

/ SetLevel sets the level of an existing dial.
func (s *DialService) SetLevel(id wtf.DialID, token string, level float64) error {
  u := *s.URL
  u.Path = "/api/dials/" + url.QueryEscape(string(id))

  // Encode request body.
  reqBody, err := json.Marshal(patchDialRequest{ID: id, Token: token, Level: level})
  if err != nil {
    return err
  }

  // Create request.
  req, err := http.NewRequest("PATCH", u.String(), bytes.NewReader(reqBody))
  if err != nil {
    return err
  }

  // Execute request.
  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    return err
  }
  defer resp.Body.Close()

  // Decode response into JSON.
  var respBody postDialResponse
  if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
    return err
  } else if respBody.Err != "" {
    return wtf.Error(respBody.Err)
  }

  return nil
}