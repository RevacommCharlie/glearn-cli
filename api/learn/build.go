package learn

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// BlockPost represents the shape of the data needed to POST to learn for
// creating a new block
type BlockPost struct {
	Block Block `json:"block"`
}

// Block holds information yielded from the Learn Block API
type Block struct {
	ID           int      `json:"id"`
	RepoName     string   `json:"repo_name"`
	SyncErrors   []string `json:"sync_errors"`
	Title        string   `json:"title"`
	CohortsUsing []int    `json:"cohorts_using"`
}

// blockReponse represents the shape of our Learn API block responses
type blockResponse struct {
	Blocks []Block `json:"blocks"`
}

// ReleaseResponse holds the release id of a fetched or created release
type ReleaseResponse struct {
	ReleaseID int `json:"release_id"`
}

// Exists reports if a Block struct has a nonzero id value
func (b Block) Exists() bool {
	return b.ID != 0
}

// GetBlockByRepoName takes a string repo name and requests a block from Learn. Returns
// either the Block or an error
func (api *APIClient) GetBlockByRepoName(repoName string) (Block, error) {
	u, err := url.Parse(fmt.Sprintf("%s/api/v1/blocks", api.baseURL))
	if err != nil {
		return Block{}, errors.New("unable to parse Learn remote")
	}
	v := url.Values{}
	v.Set("repo_name", repoName)
	u.RawQuery = v.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return Block{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.token))

	res, err := api.client.Do(req)
	if err != nil {
		return Block{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Block{}, fmt.Errorf("Error: response status: %d", res.StatusCode)
	}

	var blockResp blockResponse
	err = json.NewDecoder(res.Body).Decode(&blockResp)
	if err != nil {
		return Block{}, err
	}

	if len(blockResp.Blocks) == 1 {
		return blockResp.Blocks[0], nil
	}
	return Block{}, nil
}

// CreateBlockByRepoName takes a string repo name and makes a POST to the Learn API to create the block
func (api *APIClient) CreateBlockByRepoName(repoName string) (Block, error) {
	payload := BlockPost{Block: Block{RepoName: repoName}}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return Block{}, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/blocks", api.baseURL), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return Block{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.token))

	res, err := api.client.Do(req)
	if err != nil {
		return Block{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Block{}, fmt.Errorf("Error: response status: %d", res.StatusCode)
	}

	var blockResp blockResponse

	err = json.NewDecoder(res.Body).Decode(&blockResp)
	if err != nil {
		return Block{}, err
	}

	if len(blockResp.Blocks) == 1 {
		return blockResp.Blocks[0], nil
	}

	return Block{}, nil
}

// CreateMasterRelease takes a block ID and creates a master release from it by POSTing to the Learn API
func (api *APIClient) CreateMasterRelease(blockID int) (int, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/blocks/%d/releases", api.baseURL, blockID), nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.token))

	res, err := api.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Error: response status: %d", res.StatusCode)
	}

	var r ReleaseResponse

	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return 0, err
	}

	return r.ReleaseID, nil
}