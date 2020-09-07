package retriever_git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/joshcarp/pb-mod/app"
	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
)

type Retriever struct {
	AppConfig app.AppConfig
}

func New(appConfig app.AppConfig) Retriever {
	return Retriever{AppConfig: appConfig}
}

func (a Retriever) Retrieve(res *pbmod.Object) error {
	var auth *http.BasicAuth
	store := memory.NewStorage()
	if a.AppConfig.Username != "" {
		auth = &http.BasicAuth{
			Username: a.AppConfig.Username,
			Password: a.AppConfig.Token,
		}
	}
	r, err := git.Clone(store, nil, &git.CloneOptions{
		URL:   "https://" + res.Repo + ".git",
		Depth: 1,
		Auth:  auth,
	})
	if err != nil {
		return err
	}
	commit, err := r.CommitObject(plumbing.NewHash(res.Version))
	if err != nil {
		return err
	}
	f, err := commit.File(res.Resource)
	if err != nil {
		return err
	}
	reader, err := f.Reader()
	if err != nil {
		return err
	}
	return app.ScanIntoString(&res.Content, reader)
}