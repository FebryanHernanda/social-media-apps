package repositories

import (
	"context"
	"fmt"

	"github.com/febryanhernanda/social-media-apps/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) GetAllUser(ctx context.Context) ([]models.AllUser, error) {
	query := `
		SELECT id, name, email, avatar_path, biography, created_at
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY id
	`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.AllUser
	for rows.Next() {
		var u models.AllUser
		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.AvatarPath,
			&u.Biography,
			&u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

/* ===================================================================================================================== NOTIFICATIONS */
func (r *UserRepository) GetNotifications(ctx context.Context, userID int) ([]models.Notifications, error) {
	query := `
	SELECT n.id, n.actor_id, u.name, u.avatar_path, n.action_type, n.post_id, n.is_read, n.created_at
	FROM notifications n
	JOIN users u ON n.actor_id = u.id
	WHERE n.receiver_id = $1
	ORDER BY n.created_at DESC
	LIMIT 5
	`
	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notifications
	for rows.Next() {
		var n models.Notifications
		err := rows.Scan(
			&n.ID,
			&n.ActorID,
			&n.ActorName,
			&n.ActorAvatar,
			&n.Action,
			&n.PostID,
			&n.IsRead,
			&n.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

func (r *UserRepository) ReadNotification(ctx context.Context, notificationID, userID int) (bool, error) {
	query := `
		UPDATE notifications
		SET is_read = TRUE 
		WHERE id = $1 AND receiver_id = $2
	`

	res, err := r.DB.Exec(ctx, query, notificationID, userID)
	if err != nil {
		return false, err
	}

	if res.RowsAffected() == 0 {
		return false, nil
	}

	return true, nil
}

/* ===================================================================================================================== FOLLOWS */
func (r *UserRepository) FollowRequest(ctx context.Context, req *models.Follows) (*models.Follows, error) {
	dbTx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin db transaction: %w", err)
	}
	defer dbTx.Rollback(ctx)

	query := `
		INSERT INTO follows(user_id, followed_user_id)
		VALUES ($1,$2)
		RETURNING id, user_id, followed_user_id, created_at
	`

	values := []any{req.UserID, req.FollowedUserID}

	var userFollow models.Follows
	err = dbTx.QueryRow(ctx, query, values...).Scan(&userFollow.ID, &userFollow.UserID, &userFollow.FollowedUserID, &req.CreatedAt)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, fmt.Errorf("already following this user")
		}
		return nil, err
	}

	queryNotif := `
		INSERT INTO notifications(receiver_id, actor_id, action_type)
		VALUES ($1, $2, 'follow')

	`
	_, err = dbTx.Exec(ctx, queryNotif, req.FollowedUserID, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert notification: %w", err)
	}

	if err := dbTx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &userFollow, nil
}

func (r *UserRepository) UnfollowRequest(ctx context.Context, userID, followedUserID int) error {
	dbTx, err := r.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin db transaction: %w", err)
	}
	defer dbTx.Rollback(ctx)

	query := `
		DELETE FROM follows
		WHERE user_id = $1 AND followed_user_id = $2
	`
	res, err := dbTx.Exec(ctx, query, userID, followedUserID)
	if err != nil {
		return fmt.Errorf("failed to delete like: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("not following this user")
	}

	queryNotif := `
		DELETE FROM notifications
		WHERE receiver_id = $1 AND actor_id = $2 AND action_type = 'follow'
	`
	_, err = dbTx.Exec(ctx, queryNotif, followedUserID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	if err := dbTx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
