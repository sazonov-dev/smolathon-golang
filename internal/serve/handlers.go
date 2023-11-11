package serve

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/smolathon/internal/database"
	"github.com/smolathon/internal/handlers"
	"github.com/smolathon/internal/models"
	"github.com/smolathon/pkg/decoding"
	"github.com/smolathon/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
)

var _ handlers.Handler = &handler{} // check if structure satisfies the condition

// TODO constants

type handler struct {
	logger   *logging.Logger
	mStorage *database.MasterStorage
	cStorage *database.CardStorage
	ctx      context.Context
}

func NewHandler(logger *logging.Logger, ms *database.MasterStorage, cs *database.CardStorage) handlers.Handler {
	return &handler{
		logger:   logger,
		mStorage: ms,
		cStorage: cs,
	}
}

// Нет совместимости с дефолтным http router
func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/main", h.GetCards)
	router.HandlerFunc(http.MethodPost, "/create", h.CreateCard)
	router.HandlerFunc(http.MethodPost, "/delete", h.DeleteCard)

}

func (h *handler) GetCards(w http.ResponseWriter, r *http.Request) {
	var cards []models.Card
	var notEmptyCards []models.Card
	cursor, err := h.cStorage.Collection.Find(h.ctx, bson.M{}) // bson.M{} - ?, check doc for Find()

	if cursor.Err() != nil {
		h.logger.Errorf("failed to find all users due to error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	if err := cursor.All(h.ctx, &cards); err != nil {
		h.logger.Errorf("failed to read all documents from cursor. error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	for _, val := range cards {
		if val.Description != "" {
			notEmptyCards = append(notEmptyCards, val)
		}
	}

	jData, err := json.Marshal(notEmptyCards)
	if err != nil {
		h.logger.Errorf("error while marshaling json: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	h.logger.Info("got all cards")
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
	w.WriteHeader(http.StatusOK)
}

func (h *handler) CreateCard(w http.ResponseWriter, r *http.Request) {
	var card models.Card
	err := decoding.DecodeJSONBody(w, r, &card)
	if err != nil {
		h.logger.Errorf("error while decoding. err: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	id, err := h.mStorage.CreateCard(h.ctx, card, h.cStorage.Collection)
	if err != nil {
		h.logger.Errorf("error while creating a card: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	h.logger.Infof("created card. id: %v", id)
	w.WriteHeader(http.StatusOK)
}

type Id struct {
	Id string `json:"id"`
}

func (h *handler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	var i Id
	err := decoding.DecodeJSONBody(w, r, &i)
	if err != nil {
		h.logger.Errorf("error while decoding. err: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	err = h.mStorage.DeleteCard(h.ctx, i.Id, h.cStorage.Collection)
	if err != nil {
		h.logger.Errorf("error while deleting. err: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}
