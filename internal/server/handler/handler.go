package handler

import (
	"devops-tpl/internal/storage"
)

type Handler struct {
	storage.Storage
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{storage}
}
