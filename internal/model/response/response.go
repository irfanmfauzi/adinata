package response

import "adinata/internal/model/entity"

type GenericResponse struct {
	Success bool   `json:"status"`
	Message string `json:"message"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type PostResponse struct {
	GenericResponse
	Data entity.Post `json:"data"`
}

type LoginResponse struct {
	GenericResponse
	Data TokenResponse `json:"data"`
}

type GetPostByIdResponse struct {
	GenericResponse
	Data entity.Post `json:"data"`
}

type GetPostResponse struct {
	GenericResponse
	Data []entity.Post `json:"data"`
}

type GetTagResponse struct {
	GenericResponse
	Data []entity.Tag `json:"data"`
}
