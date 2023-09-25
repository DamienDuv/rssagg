// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0

package database

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error)
	CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) (FeedFollow, error)
	CreatePost(ctx context.Context, arg CreatePostParams) (Post, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteFeedFollow(ctx context.Context, arg DeleteFeedFollowParams) error
	GetFeedFollows(ctx context.Context, userID uuid.UUID) ([]FeedFollow, error)
	GetFeeds(ctx context.Context) ([]Feed, error)
	GetNextFeedsToFetch(ctx context.Context, limit int32) ([]Feed, error)
	GetUserByAPIKey(ctx context.Context, apiKey string) (User, error)
	MarkFeedAsFetched(ctx context.Context, id uuid.UUID) (Feed, error)
}

var _ Querier = (*Queries)(nil)