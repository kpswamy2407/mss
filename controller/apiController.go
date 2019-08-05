package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/google/go-github/github"
	"github.com/kpswamy540/verloop/config"
	"github.com/kpswamy540/verloop/helper"
	validator "gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

//APIController is a controller having the methods related API
type APIController struct {
}

// Test is used to check if the application working or not
func (a *APIController) Test(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Hello, you've requested: %s\n", req.URL.Path)
}

// GetQuestionIDRequest is structure hold the request parameters
type GetQuestionIDRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Name      string `json:"name" validate:"required"`
	AngelList string `json:"angel_list" validate:"required"`
	Github    string `json:"github" validate:"required"`
}
type questionIDResponse struct {
	QuestionID string `json:"question_id"`
}
type curlResponse struct {
	Error    string             `json:"error"`
	Response questionIDResponse `json:"response"`
}

//GetRepositoriesRequest is structure hold the request parameters
type GetRepositoriesRequest struct {
	Org string `json:"org" validate:"required"`
}

// GetQuestionIDResponse collects the response parameters for the GetQuestionID method.
type GetQuestionIDResponse struct {
	Message    string `json:"message"`
	QuestionID string `json:"question_id"`
}

//GetRepositoriesResponse collects reponse of GetRepositoriesRequest method
type GetRepositoriesResponse struct {
	Message string              `json:"message"`
	Results []github.Repository `json:"results"`
}

//GetQuestionID is used the question id for challenge
func (a *APIController) GetQuestionID(res http.ResponseWriter, req *http.Request) {
	var getQuestionIDRequest GetQuestionIDRequest
	var getQuestionIDResponse GetQuestionIDResponse
	json.NewDecoder(req.Body).Decode(&getQuestionIDRequest)
	validate = validator.New()
	err := validate.Struct(getQuestionIDRequest)
	if err != nil {

		getQuestionIDResponse.Message = err.Error()
		json.NewEncoder(res).Encode(getQuestionIDResponse)
	} else {
		body := new(bytes.Buffer)
		json.NewEncoder(body).Encode(getQuestionIDRequest)
		curlReq, err := http.NewRequest("POST", config.ChallengeURL, body)
		if err != nil {

			getQuestionIDResponse.Message = err.Error()
			json.NewEncoder(res).Encode(getQuestionIDResponse)
		} else {
			curlReq.Header.Add("x-verloop-password", helper.GenerateMd5(getQuestionIDRequest.Email))
			curlReq.Header.Add("content-type", "application/json")
			curlReq.Header.Add("cache-control", "no-cache")
			curlRes, _ := http.DefaultClient.Do(curlReq)
			body, _ := ioutil.ReadAll(curlRes.Body)
			data := curlResponse{}
			json.Unmarshal(body, &data)
			if len(data.Error) > 0 {
				getQuestionIDResponse.Message = data.Error
				json.NewEncoder(res).Encode(getQuestionIDResponse)
			} else {
				getQuestionIDResponse.Message = "Request completed successfully!"
				getQuestionIDResponse.QuestionID = data.Response.QuestionID
				json.NewEncoder(res).Encode(getQuestionIDResponse)
			}
		}
	}
}

//GetRepositories is used to get the repositories
func (a *APIController) GetRepositories(res http.ResponseWriter, req *http.Request) {

	var getRepositoriesRequest GetRepositoriesRequest
	var getRepositoriesResponse GetRepositoriesResponse
	json.NewDecoder(req.Body).Decode(&getRepositoriesRequest)
	validate = validator.New()
	err := validate.Struct(getRepositoriesRequest)
	if err != nil {
		getRepositoriesResponse.Message = err.Error()
		json.NewEncoder(res).Encode(getRepositoriesResponse)
	} else {
		client := github.NewClient(nil)
		opt := &github.SearchOptions{Sort: "stars", Order: "desc"}
		query := "org:" + getRepositoriesRequest.Org
		searchRepos, _, err := client.Search.Repositories(context.Background(), query, opt)
		repos := searchRepos.Repositories
		if err != nil {
			getRepositoriesResponse.Message = err.Error()
			json.NewEncoder(res).Encode(getRepositoriesResponse)
		} else {
			if len(repos) > 0 {
				sort.SliceStable(repos, func(i, j int) bool { return *repos[i].StargazersCount > *repos[j].StargazersCount })
				sortedRepos := make([]github.Repository, 3)
				i := 0
				for _, repo := range repos {
					sortedRepos[i].Name = repo.Name
					sortedRepos[i].StargazersCount = repo.StargazersCount
					i++
					if i == 3 {
						break
					}
				}
				getRepositoriesResponse.Message = "Request completed successfully"
				getRepositoriesResponse.Results = sortedRepos
				json.NewEncoder(res).Encode(getRepositoriesResponse)
			} else {
				getRepositoriesResponse.Message = "No repositories found for the org:" + getRepositoriesRequest.Org
				json.NewEncoder(res).Encode(getRepositoriesResponse)
			}

		}
	}
}
