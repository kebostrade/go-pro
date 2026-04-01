package loader

import (
	"context"

	"basic/projects/graphql/pkg/models"
)

type Loaders struct {
	UserLoader *UserLoader
	PostLoader *PostLoader
}

func NewLoaders(db *models.DB) *Loaders {
	return &Loaders{
		UserLoader: newUserLoader(db),
		PostLoader: newPostLoader(db),
	}
}

type loadersKey struct{}

func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey{}).(*Loaders)
}

type userResult struct {
	user *models.User
	err  error
}

type UserLoader struct {
	load func(keys []string) []userResult
}

func newUserLoader(db *models.DB) *UserLoader {
	return &UserLoader{
		load: func(keys []string) []userResult {
			results := make([]userResult, len(keys))
			for i, id := range keys {
				if u, err := db.UserByID(id); err != nil {
					results[i] = userResult{err: err}
				} else {
					results[i] = userResult{user: u}
				}
			}
			return results
		},
	}
}

func (l *UserLoader) Load(ctx context.Context, id string) (*models.User, error) {
	results := l.load([]string{id})
	return results[0].user, results[0].err
}

type postResult struct {
	posts []*models.Post
	err   error
}

type PostLoader struct {
	load func(authorID string) postResult
}

func newPostLoader(db *models.DB) *PostLoader {
	return &PostLoader{
		load: func(authorID string) postResult {
			return postResult{posts: db.PostsByAuthorID(authorID)}
		},
	}
}

func (l *PostLoader) LoadByAuthor(ctx context.Context, authorID string) []*models.Post {
	return l.load(authorID).posts
}
