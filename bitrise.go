package bitrise

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type Trigger struct {
	HookInfo    HookInfo    `json:"hook_info"`
	BuildParams BuildParams `json:"build_params"`
}

type HookInfo struct {
	Type     string `json:"type"`
	APIToken string `json:"api_token"`
}

type BuildParams struct {
	Branch     *string `json:"branch,omitempty"`
	Tag        *string `json:"tag,omitempty"`
	BranchDest *string `json:"branch_dest,omitempty"`
	CommitHash *string `json:"commit_hash"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type BitRise struct {
	slug  string
	token string
}

func New(Slug, Token string) *BitRise {
	return &BitRise{Slug, Token}
}

func (b BitRise) Trigger(ref, hash string) error {
	branch, err := branchFromRef(ref)
	if err != nil {
		return errors.Wrapf(err, "invalid ref")
	}

	t := Trigger{
		HookInfo: HookInfo{
			Type:     "bitrise",
			APIToken: b.token,
		},
		BuildParams: BuildParams{
			Branch:     &branch,
			CommitHash: &hash,
		},
	}

	api := fmt.Sprintf("https://www.bitrise.io/app/%v/build/start.json", b.slug)
	data, err := json.Marshal(t)
	if err != nil {
		return errors.Wrapf(err, "unable to marshal trigger to json format")
	}

	r, err := http.Post(api, "application/json", bytes.NewReader(data))
	if err != nil {
		return errors.Wrapf(err, "trigger api call failed")
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrapf(err, "failed to read response")
	}
	res := Response{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return errors.Wrapf(err, "failed to parse api response: [%s]", body)
	}

	if res.Status == "error" {
		return fmt.Errorf("api call error: %v", res.Message)
	}
	return nil

}

func branchFromRef(ref string) (string, error) {
	parts := strings.Split(ref, "/")

	if len(parts) < 3 {
		return "", fmt.Errorf("unable to parse ref %v, expecting at least 3 parts", ref)
	}

	if parts[0] != "refs" {
		return "", fmt.Errorf("unable to parse ref %v, first part to be 'ref', parts are %v", ref, parts)
	}

	if parts[1] != "heads" {
		return "", fmt.Errorf("ref passed is not a branch")
	}

	return parts[2], nil
}
