package handler

import (
	"metrics-tpl/internal/server/storage"
)

type Handler struct {
	Storage *storage.MemStorage
}

func NewHandler(storage *storage.MemStorage) *Handler {
	return &Handler{storage}
}
