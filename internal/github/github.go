package githubTrack

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

var OpenAPIGithub string = "https://api.github.com/search/commits?q=author:"

type CommitAuthor struct {
	Date string `json:"date"`
}

type CommitDetail struct {
	Message string       `json:"message"`
	Author  CommitAuthor `json:"author"`
}

type Repository struct {
	Name string `json:"name"`
}

type CommitItem struct {
	Commit     CommitDetail `json:"commit"`
	Repository Repository   `json:"repository"`
}

type SearchResult struct {
	TotalCount int          `json:"total_count"`
	Items      []CommitItem `json:"items"`
}

func FetchCommits(username string) ([]CommitItem, error) {
	request, err := http.NewRequest("GET", OpenAPIGithub+username, nil)
	if err != nil {
		log.Fatal("Error when requesting Github API")
		return nil, err
	}
	request.Header.Set("Accept", "application/vnd.github.cloak-preview")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal("Error when getting answer")
	}
	defer response.Body.Close()
	content, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error when reading request's body")
		return nil, err
	}
	var result SearchResult
	erreur := json.Unmarshal(content, &result)
	if erreur != nil {
		log.Printf("Unmarshal error: %v", erreur)
	}
	return result.Items, nil
}

func CountCommitsByProject(commits []CommitItem) map[string]int {
	m := make(map[string]int)

	for _, commit := range commits {
		m[commit.Repository.Name]++
	}
	return m
}
