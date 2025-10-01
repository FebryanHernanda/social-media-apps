package repositories

import (
	"context"
	"fmt"

	"github.com/febryanhernanda/social-media-apps/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostRepository struct {
	DB *pgxpool.Pool
}

func NewPostRepository(db *pgxpool.Pool) *PostRepository {
	return &PostRepository{
		DB: db,
	}
}

func (r *PostRepository) CreatePost(ctx context.Context, req *models.Post) (*models.Post, error) {
	query := `
		INSERT INTO posts(content, image_path, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, content, image_path, user_id, created_at
	`
	values := []any{req.Content, req.ImagePath, req.UserID}

	var post models.Post
	err := r.DB.QueryRow(ctx, query, values...).Scan(&post.ID, &post.Content, &post.ImagePath, &post.UserID, &post.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) GetPostOwnerID(ctx context.Context, postID int) (int, error) {
	var ownerID int
	err := r.DB.QueryRow(ctx, "SELECT user_id FROM posts WHERE id=$1", postID).Scan(&ownerID)
	if err != nil {
		return 0, fmt.Errorf("post not found")
	}
	return ownerID, nil
}

/* ===================================================================================================================== LIKES */
func (r *PostRepository) LikePost(ctx context.Context, postID, userID, postOwnerID int) (*models.Like, error) {
	dbTx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin db transaction: %w", err)
	}
	defer dbTx.Rollback(ctx)

	query := `
        INSERT INTO likes (post_id, user_id)
        VALUES ($1, $2)
        ON CONFLICT DO NOTHING
        RETURNING id, liked_at
    `

	var likePost models.Like
	likePost.PostID = postID
	likePost.UserID = userID

	err = dbTx.QueryRow(ctx, query, postID, userID).Scan(&likePost.ID, &likePost.LikedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert like: %w", err)
	}

	if userID != postOwnerID {
		queryNotif := `
            INSERT INTO notifications (receiver_id, actor_id, action_type, post_id)
            VALUES ($1, $2, 'like', $3)
        `
		_, err = dbTx.Exec(ctx, queryNotif, postOwnerID, userID, postID)
		if err != nil {
			return nil, fmt.Errorf("failed to insert notification: %w", err)
		}
	}

	if err := dbTx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &likePost, nil
}

func (r *PostRepository) UnlikePost(ctx context.Context, postID, userID int) error {
	dbTx, err := r.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer dbTx.Rollback(ctx)

	query := ` 
		DELETE FROM likes
        WHERE post_id=$1 AND user_id=$2
	`
	res, err := dbTx.Exec(ctx, query, postID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete like: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("like not found")
	}

	queryNotif := `
	    DELETE FROM notifications
        WHERE actor_id=$1 AND post_id=$2 AND action_type='like'
	`
	_, err = dbTx.Exec(ctx, queryNotif, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	if err := dbTx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

/* ===================================================================================================================== COMMENT */

func (r *PostRepository) AddComment(ctx context.Context, comment *models.Comment, postOwnerID int) (*models.Comment, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	queryComment := `
        INSERT INTO comments (content, post_id, user_id)
        VALUES ($1, $2, $3)
        RETURNING id, created_at
    `
	err = tx.QueryRow(ctx, queryComment, comment.Content, comment.PostID, comment.UserID).
		Scan(&comment.ID, &comment.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert comment: %w", err)
	}

	if comment.UserID != postOwnerID {
		queryNotif := `
            INSERT INTO notifications (receiver_id, actor_id, action_type, post_id)
            VALUES ($1, $2, 'comment', $3)
        `
		_, err = tx.Exec(ctx, queryNotif, postOwnerID, comment.UserID, comment.PostID)
		if err != nil {
			return nil, fmt.Errorf("failed to insert notification: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return comment, nil
}
