package models

import (
    "time"

    "github.com/google/uuid"
)

type Subscription struct {
    ID          uuid.UUID  `json:"id" db:"id"`
    ServiceName string     `json:"service_name" db:"service_name"`
    Price       int        `json:"price" db:"price"`
    UserID      uuid.UUID  `json:"user_id" db:"user_id"`
    StartDate   time.Time  `json:"start_date" db:"start_date"`
    EndDate     *time.Time `json:"end_date,omitempty" db:"end_date"`
    CreatedAt   time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateSubscriptionRequest struct {
    ServiceName string    `json:"service_name" binding:"required"`
    Price       int       `json:"price" binding:"required,min=0"`
    UserID      string    `json:"user_id" binding:"required"`
    StartDate   string    `json:"start_date" binding:"required"`
    EndDate     *string   `json:"end_date,omitempty"`
}

type UpdateSubscriptionRequest struct {
    ServiceName *string `json:"service_name,omitempty"`
    Price       *int    `json:"price,omitempty"`
    StartDate   *string `json:"start_date,omitempty"`
    EndDate     *string `json:"end_date,omitempty"`
}

type SubscriptionResponse struct {
    ID          string  `json:"id"`
    ServiceName string  `json:"service_name"`
    Price       int     `json:"price"`
    UserID      string  `json:"user_id"`
    StartDate   string  `json:"start_date"`
    EndDate     *string `json:"end_date,omitempty"`
    CreatedAt   string  `json:"created_at"`
    UpdatedAt   string  `json:"updated_at"`
}

type TotalCostRequest struct {
    StartDate   string  `json:"start_date" binding:"required"`
    EndDate     string  `json:"end_date" binding:"required"`
    UserID      *string `json:"user_id,omitempty"`
    ServiceName *string `json:"service_name,omitempty"`
}

type TotalCostResponse struct {
    TotalCost int `json:"total_cost"`
}