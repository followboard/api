package mapping

import (
	"time"
)

// PR serialization structure for indexing pull requests
type PR struct {
	FollowboardUserLogin string    `json:"followboardUserLogin"`
	Number               int       `json:"number"`
	URL                  string    `json:"url"`
	Title                string    `json:"title"`
	Body                 string    `json:"body"`
	State                string    `json:"state"`
	Commits              int       `json:"commits"`
	Additions            int       `json:"additions"`
	Deletions            int       `json:"deletions"`
	ChangedFiles         int       `json:"changedFiles"`
	UserLogin            string    `json:"userLogin"`
	HeadRepoOwnerLogin   string    `json:"headRepoOwnerLogin"`
	HeadRepoName         string    `json:"headRepoName"`
	HeadRef              string    `json:"headRef"`
	BaseRepoOwnerLogin   string    `json:"baseRepoOwnerLogin"`
	BaseRepoName         string    `json:"baseRepoName"`
	BaseRef              string    `json:"baseRef"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}
