package dto

import "andreasho/scalable-ecomm/services/product/internal/db/models"

type ProductResponse struct {
	Product models.Product `json:"product"`
}
