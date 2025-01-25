package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/shekshuev/gophertalk-backend/internal/models"
	"github.com/shekshuev/gophertalk-backend/internal/utils"
)

func (h *Handler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		offset = 0
	}
	replyToID, err := strconv.ParseUint(r.URL.Query().Get("reply_to_id"), 10, 64)
	if err != nil {
		replyToID = 0
	}
	ownerID, err := strconv.ParseUint(r.URL.Query().Get("owner_id"), 10, 64)
	if err != nil {
		ownerID = 0
	}

	claims, ok := utils.GetClaimsFromContext(r.Context())
	if !ok {
		h.JSONError(w, http.StatusUnauthorized, ErrInvalidToken.Error())
		return
	}
	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		h.JSONError(w, http.StatusNotFound, ErrInvalidToken.Error())
		return
	}
	filterDTO := models.FilterPostDTO{
		OwnerID:   ownerID,
		UserID:    userID,
		Limit:     limit,
		Offset:    offset,
		ReplyToID: replyToID,
		Search:    r.URL.Query().Get("search"),
	}
	readDTOs, err := h.posts.GetAllPosts(filterDTO)
	if err != nil {
		h.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := json.Marshal(readDTOs)
	if err != nil {
		h.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		h.JSONError(w, http.StatusBadRequest, err.Error())
	}
}

func (h *Handler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.JSONError(w, http.StatusNotFound, ErrInvalidID.Error())
		return
	}
	claims, ok := utils.GetClaimsFromContext(r.Context())
	if !ok {
		h.JSONError(w, http.StatusUnauthorized, ErrInvalidToken.Error())
		return
	}
	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		h.JSONError(w, http.StatusNotFound, ErrInvalidToken.Error())
		return
	}
	readDTO, err := h.posts.GetPostByID(id, userID)
	if err != nil {
		h.JSONError(w, http.StatusNotFound, err.Error())
		return
	}
	resp, err := json.Marshal(readDTO)
	if err != nil {
		h.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		h.JSONError(w, http.StatusBadRequest, err.Error())
	}
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	var createDTO models.CreatePostDTO
	if err = json.Unmarshal(body, &createDTO); err != nil {
		h.JSONError(w, http.StatusNotFound, ErrInvalidID.Error())
		return
	}
	claims, ok := utils.GetClaimsFromContext(r.Context())
	if !ok {
		h.JSONError(w, http.StatusUnauthorized, ErrInvalidToken.Error())
		return
	}
	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		h.JSONError(w, http.StatusNotFound, ErrInvalidToken.Error())
		return
	}
	createDTO.UserID = userID
	err = h.validate.Struct(createDTO)
	if err != nil {
		h.JSONError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	readDTO, err := h.posts.CreatePost(createDTO)
	if err != nil {
		h.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := json.Marshal(readDTO)
	if err != nil {
		h.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(resp)
	if err != nil {
		h.JSONError(w, http.StatusBadRequest, err.Error())
	}
}

func (h *Handler) DeletePostByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.JSONError(w, http.StatusNotFound, ErrInvalidID.Error())
		return
	}
	claims, ok := utils.GetClaimsFromContext(r.Context())
	if !ok {
		h.JSONError(w, http.StatusUnauthorized, ErrInvalidToken.Error())
		return
	}
	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		h.JSONError(w, http.StatusNotFound, ErrInvalidToken.Error())
		return
	}
	err = h.posts.DeletePost(id, userID)
	if err != nil {
		h.JSONError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ViewPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.JSONError(w, http.StatusNotFound, ErrInvalidID.Error())
		return
	}
	claims, ok := utils.GetClaimsFromContext(r.Context())
	if !ok {
		h.JSONError(w, http.StatusUnauthorized, ErrInvalidToken.Error())
		return
	}
	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		h.JSONError(w, http.StatusNotFound, ErrInvalidToken.Error())
		return
	}
	err = h.posts.ViewPost(id, userID)
	if err != nil {
		h.JSONError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) LikePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.JSONError(w, http.StatusNotFound, ErrInvalidID.Error())
		return
	}
	claims, ok := utils.GetClaimsFromContext(r.Context())
	if !ok {
		h.JSONError(w, http.StatusUnauthorized, ErrInvalidToken.Error())
		return
	}
	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		h.JSONError(w, http.StatusNotFound, ErrInvalidToken.Error())
		return
	}
	err = h.posts.LikePost(id, userID)
	if err != nil {
		h.JSONError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) DislikePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.JSONError(w, http.StatusNotFound, ErrInvalidID.Error())
		return
	}
	claims, ok := utils.GetClaimsFromContext(r.Context())
	if !ok {
		h.JSONError(w, http.StatusUnauthorized, ErrInvalidToken.Error())
		return
	}
	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		h.JSONError(w, http.StatusNotFound, ErrInvalidToken.Error())
		return
	}
	err = h.posts.DislikePost(id, userID)
	if err != nil {
		h.JSONError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
