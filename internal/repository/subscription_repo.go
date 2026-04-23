package repository

import (
    "database/sql"
    "fmt"
    "time"

    "subscription-service/internal/models"

    "github.com/google/uuid"
    log "github.com/sirupsen/logrus"
)

type SubscriptionRepository struct {
    db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
    return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(sub *models.Subscription) error {
    query := `
        INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `
    
    sub.ID = uuid.New()

    log.WithFields(log.Fields{
        "id":           sub.ID,
        "service_name": sub.ServiceName,
        "user_id":      sub.UserID,
        "price":        sub.Price,
    }).Info("Creating new subscription")

    err := r.db.QueryRow(query, sub.ID, sub.ServiceName, sub.Price, sub.UserID, 
        sub.StartDate, sub.EndDate).Scan(&sub.ID, &sub.CreatedAt, &sub.UpdatedAt)
    
    if err != nil {
        log.WithError(err).Error("Failed to create subscription")
        return fmt.Errorf("failed to create subscription: %w", err)
    }

    log.WithField("id", sub.ID).Info("Subscription created successfully")
    return nil
}

func (r *SubscriptionRepository) GetByID(id uuid.UUID) (*models.Subscription, error) {
    query := `
        SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
        FROM subscriptions
        WHERE id = $1 AND (end_date IS NULL OR end_date > NOW())
    `

    var sub models.Subscription
    err := r.db.QueryRow(query, id).Scan(
        &sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID,
        &sub.StartDate, &sub.EndDate, &sub.CreatedAt, &sub.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        log.WithField("id", id).Warn("Subscription not found")
        return nil, nil
    }
    if err != nil {
        log.WithError(err).Error("Failed to get subscription")
        return nil, fmt.Errorf("failed to get subscription: %w", err)
    }

    log.WithField("id", id).Info("Subscription retrieved successfully")
    return &sub, nil
}

func (r *SubscriptionRepository) Update(id uuid.UUID, req *models.UpdateSubscriptionRequest) error {
    // Начинаем построение запроса
    query := "UPDATE subscriptions SET updated_at = NOW()"
    args := []interface{}{}
    argPos := 1

    if req.ServiceName != nil {
        query += fmt.Sprintf(", service_name = $%d", argPos)
        args = append(args, *req.ServiceName)
        argPos++
    }

    if req.Price != nil {
        query += fmt.Sprintf(", price = $%d", argPos)
        args = append(args, *req.Price)
        argPos++
    }

    if req.StartDate != nil {
        startDate, err := time.Parse("01-2006", *req.StartDate)
        if err != nil {
            return fmt.Errorf("invalid start_date format: %w", err)
        }
        query += fmt.Sprintf(", start_date = $%d", argPos)
        args = append(args, startDate)
        argPos++
    }

    if req.EndDate != nil {
        var endDate *time.Time
        if *req.EndDate != "" {
            parsed, err := time.Parse("01-2006", *req.EndDate)
            if err != nil {
                return fmt.Errorf("invalid end_date format: %w", err)
            }
            endDate = &parsed
        }
        query += fmt.Sprintf(", end_date = $%d", argPos)
        args = append(args, endDate)
        argPos++
    }

    query += fmt.Sprintf(" WHERE id = $%d AND (end_date IS NULL OR end_date > NOW())", argPos)
    args = append(args, id)

    log.WithFields(log.Fields{
        "id":           id,
        "fields_count": argPos - 1,
    }).Info("Updating subscription")

    result, err := r.db.Exec(query, args...)
    if err != nil {
        log.WithError(err).Error("Failed to update subscription")
        return fmt.Errorf("failed to update subscription: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        log.WithField("id", id).Warn("Subscription not found for update")
        return fmt.Errorf("subscription not found")
    }

    log.WithField("id", id).Info("Subscription updated successfully")
    return nil
}

func (r *SubscriptionRepository) Delete(id uuid.UUID) error {
    query := "DELETE FROM subscriptions WHERE id = $1"

    log.WithField("id", id).Info("Deleting subscription")

    result, err := r.db.Exec(query, id)
    if err != nil {
        log.WithError(err).Error("Failed to delete subscription")
        return fmt.Errorf("failed to delete subscription: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        log.WithField("id", id).Warn("Subscription not found for deletion")
        return fmt.Errorf("subscription not found")
    }

    log.WithField("id", id).Info("Subscription deleted successfully")
    return nil
}

func (r *SubscriptionRepository) GetTotalCost(startDate, endDate time.Time, userID, serviceName *string) (int, error) {
    query := `
        SELECT COALESCE(SUM(price), 0)
        FROM subscriptions
        WHERE start_date >= $1 
        AND start_date <= $2
        AND (end_date IS NULL OR end_date >= start_date)
    `
    args := []interface{}{startDate, endDate}
    argPos := 3

    if userID != nil && *userID != "" {
        parsedUserID, err := uuid.Parse(*userID)
        if err == nil {
            query += fmt.Sprintf(" AND user_id = $%d", argPos)
            args = append(args, parsedUserID)
            argPos++
        }
    }

    if serviceName != nil && *serviceName != "" {
        query += fmt.Sprintf(" AND service_name = $%d", argPos)
        args = append(args, *serviceName)
        argPos++
    }

    var total int
    err := r.db.QueryRow(query, args...).Scan(&total)
    if err != nil {
        log.WithError(err).Error("Failed to calculate total cost")
        return 0, fmt.Errorf("failed to calculate total cost: %w", err)
    }

    log.WithFields(log.Fields{
        "total_cost":   total,
        "start_date":   startDate.Format("01-2006"),
        "end_date":     endDate.Format("01-2006"),
        "user_id":      userID,
        "service_name": serviceName,
    }).Info("Total cost calculated successfully")

    return total, nil
}