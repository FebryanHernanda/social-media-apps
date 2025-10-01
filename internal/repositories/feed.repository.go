package repositories

import (
	"context"
	"encoding/json"

	"github.com/febryanhernanda/social-media-apps/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FeedRepository struct {
	DB *pgxpool.Pool
}

func NewFeedRepository(db *pgxpool.Pool) *FeedRepository {
	return &FeedRepository{
		DB: db,
	}
}

func (r *FeedRepository) GetUserFeed(ctx context.Context, userID int) ([]models.FeedPost, error) {
	query := `
    SELECT 
        p.id AS post_id,
        p.content,
        p.image_path,
        p.user_id AS author_id,
        u.name AS author_name,
        u.avatar_path AS author_avatar,
        p.created_at,
        COALESCE(likes_count.count, 0) AS like_count,
        COALESCE(comments_data.comments, '[]')::json AS comments
    FROM posts p
    JOIN users u ON p.user_id = u.id
    JOIN follows f ON f.followed_user_id = p.user_id
    LEFT JOIN (
        SELECT post_id, COUNT(*) AS count
        FROM likes
        GROUP BY post_id
    ) likes_count ON likes_count.post_id = p.id
    LEFT JOIN (
        SELECT 
            c.post_id,
            JSON_AGG(JSON_BUILD_OBJECT(
                'id', c.id,
				'post_id', c.post_id, 
                'user_id', c.user_id,
                'name', u.name,
                'avatar', u.avatar_path,
                'content', c.content,
                'created_at', c.created_at
            ) ORDER BY c.created_at ASC) AS comments
        FROM comments c
        JOIN users u ON c.user_id = u.id
        WHERE c.deleted_at IS NULL
        GROUP BY c.post_id
    ) comments_data ON comments_data.post_id = p.id
    WHERE f.user_id = $1
      AND p.deleted_at IS NULL
    ORDER BY p.created_at DESC
    LIMIT 10
    `

	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []models.FeedPost
	for rows.Next() {
		var post models.FeedPost
		var commentsJSON []byte

		err := rows.Scan(
			&post.ID,
			&post.Content,
			&post.ImagePath,
			&post.AuthorID,
			&post.AuthorName,
			&post.AvatarPath,
			&post.CreatedAt,
			&post.LikeCount,
			&commentsJSON,
		)
		if err != nil {
			return nil, err
		}

		post.CreatedAtStr = post.CreatedAt.Format("2006-01-02T15:04:05")

		if len(commentsJSON) > 0 {
			var comments []models.Comment
			if err := json.Unmarshal(commentsJSON, &comments); err != nil {
				return nil, err
			}
			post.Comments = comments
		} else {
			post.Comments = []models.Comment{}
		}

		feeds = append(feeds, post)
	}

	return feeds, nil
}
