package reposter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	GithubApiVersion = "2022-11-28"
	GithubApiUrl     = "https://api.github.com"
	TimeoutSeconds   = 30
	DismissMessage   = "Review reposted by automated tool."
)

var JulesUsernames = map[string]bool{
	"google-labs-jules":      true,
	"google-labs-jules[bot]": true,
}

const GraphqlMutationBatchSize = 50

const getThreadsQuery = `
query($owner: String!, $name: String!, $pullNumber: Int!, $after: String) {
  repository(owner: $owner, name: $name) {
    pullRequest(number: $pullNumber) {
      reviewThreads(first: 100, after: $after) {
        pageInfo {
          hasNextPage
          endCursor
        }
        nodes {
          id
          isResolved
          comments(first: 1) {
            nodes {
              pullRequestReview {
                databaseId
              }
            }
          }
          lastComments: comments(last: 1) {
            nodes {
              author {
                login
              }
            }
          }
        }
      }
    }
  }
}
`

type ReviewReposter struct {
	Token      string
	Repo       string
	PullNumber string
}

func NewReviewReposter(token, repo, pullNumber string) *ReviewReposter {
	return &ReviewReposter{
		Token:      token,
		Repo:       repo,
		PullNumber: pullNumber,
	}
}

func (r *ReviewReposter) getHeaders() map[string]string {
	return map[string]string{
		"Accept":               "application/vnd.github+json",
		"Authorization":        fmt.Sprintf("Bearer %s", r.Token),
		"X-GitHub-Api-Version": GithubApiVersion,
		"Content-Type":         "application/json",
	}
}

func (r *ReviewReposter) executeRequest(method, url string, body io.Reader) (interface{}, string, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, "", err
	}

	for k, v := range r.getHeaders() {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: TimeoutSeconds * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return nil, "", fmt.Errorf("GitHub API error: %d %s\n%s", resp.StatusCode, resp.Status, buf.String())
	}

	var data interface{}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	if len(respBody) > 0 {
		if err := json.Unmarshal(respBody, &data); err != nil {
			return nil, "", err
		}
	}

	nextURL := ""
	linkHeader := resp.Header.Get("Link")
	if linkHeader != "" {
		re := regexp.MustCompile(`<([^>]+)>;\s*rel="next"`)
		match := re.FindStringSubmatch(linkHeader)
		if len(match) > 1 {
			nextURL = match[1]
		}
	}

	return data, nextURL, nil
}

func (r *ReviewReposter) fetchJSON(url string) (interface{}, error) {
	data, nextURL, err := r.executeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	listData, ok := data.([]interface{})
	if !ok {
		return data, nil
	}

	allData := listData
	for nextURL != "" {
		pageData, next, err := r.executeRequest("GET", nextURL, nil)
		if err != nil {
			return nil, err
		}
		if pList, ok := pageData.([]interface{}); ok {
			allData = append(allData, pList...)
		} else {
			break
		}
		nextURL = next
	}
	return allData, nil
}

func (r *ReviewReposter) postJSON(url string, payload interface{}, method string) (interface{}, error) {
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	data, _, err := r.executeRequest(method, url, bytes.NewBuffer(jsonBody))
	return data, err
}

func (r *ReviewReposter) executeGraphQL(query string, variables map[string]interface{}) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}
	data, err := r.postJSON(fmt.Sprintf("%s/graphql", GithubApiUrl), payload, "POST")
	if err != nil {
		return nil, err
	}
	if m, ok := data.(map[string]interface{}); ok {
		return m, nil
	}
	return nil, fmt.Errorf("unexpected GraphQL response format")
}

func (r *ReviewReposter) dismissReview(reviewID string) {
	url := fmt.Sprintf("%s/reviews/%s/dismissals", r.getPullApiURL(), reviewID)
	payload := map[string]string{"message": DismissMessage}
	_, err := r.postJSON(url, payload, "PUT")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to dismiss review %s: %v\n", reviewID, err)
	} else {
		fmt.Fprintf(os.Stderr, "Dismissed review %s.\n", reviewID)
	}
}

func (r *ReviewReposter) fetchAllReviewThreads() ([]interface{}, error) {
	repoParts := strings.SplitN(r.Repo, "/", 2)
	if len(repoParts) != 2 {
		return nil, fmt.Errorf("invalid repo format: %s", r.Repo)
	}
	owner, name := repoParts[0], repoParts[1]

	var pullNumber int
	fmt.Sscanf(r.PullNumber, "%d", &pullNumber)

	var allThreads []interface{}
	cursor := interface{}(nil)
	hasNextPage := true

	for hasNextPage {
		variables := map[string]interface{}{
			"owner":      owner,
			"name":       name,
			"pullNumber": pullNumber,
			"after":      cursor,
		}
		resp, err := r.executeGraphQL(getThreadsQuery, variables)
		if err != nil {
			return nil, err
		}

		data, ok := resp["data"].(map[string]interface{})
		if !ok {
			break
		}
		repo, ok := data["repository"].(map[string]interface{})
		if !ok {
			break
		}
		pr, ok := repo["pullRequest"].(map[string]interface{})
		if !ok {
			break
		}
		threads, ok := pr["reviewThreads"].(map[string]interface{})
		if !ok {
			break
		}

		nodes, _ := threads["nodes"].([]interface{})
		allThreads = append(allThreads, nodes...)

		pageInfo, _ := threads["pageInfo"].(map[string]interface{})
		hasNextPage, _ = pageInfo["hasNextPage"].(bool)
		cursor = pageInfo["endCursor"]
	}

	return allThreads, nil
}

func (r *ReviewReposter) filterThreadsForReview(threads []interface{}, reviewDbID float64) []string {
	var results []string
	for _, t := range threads {
		thread := t.(map[string]interface{})
		if resolved, _ := thread["isResolved"].(bool); resolved {
			continue
		}

		comments, ok := thread["comments"].(map[string]interface{})
		if !ok {
			continue
		}
		nodes, _ := comments["nodes"].([]interface{})
		if len(nodes) == 0 {
			continue
		}

		firstComment := nodes[0].(map[string]interface{})
		review, _ := firstComment["pullRequestReview"].(map[string]interface{})
		if review != nil {
			id, _ := review["databaseId"].(float64)
			if id == reviewDbID {
				results = append(results, thread["id"].(string))
			}
		}
	}
	return results
}

func (r *ReviewReposter) filterJulesThreads(threads []interface{}) []string {
	var results []string
	for _, t := range threads {
		thread := t.(map[string]interface{})
		if resolved, _ := thread["isResolved"].(bool); resolved {
			continue
		}

		lastComments, ok := thread["lastComments"].(map[string]interface{})
		if !ok {
			continue
		}
		nodes, _ := lastComments["nodes"].([]interface{})
		if len(nodes) == 0 {
			continue
		}

		authorInfo, _ := nodes[0].(map[string]interface{})["author"].(map[string]interface{})
		if authorInfo != nil {
			login, _ := authorInfo["login"].(string)
			if JulesUsernames[login] {
				results = append(results, thread["id"].(string))
			}
		}
	}
	return results
}

func (r *ReviewReposter) resolveThreadsById(threadIDs []string) {
	if len(threadIDs) == 0 {
		return
	}

	for i := 0; i < len(threadIDs); i += GraphqlMutationBatchSize {
		end := i + GraphqlMutationBatchSize
		if end > len(threadIDs) {
			end = len(threadIDs)
		}
		batch := threadIDs[i:end]

		var mutations []string
		variables := make(map[string]interface{})
		for j, id := range batch {
			alias := fmt.Sprintf("resolve%d", j)
			varName := fmt.Sprintf("threadId%d", j)
			mutations = append(mutations, fmt.Sprintf("  %s: resolveReviewThread(input: {threadId: $%s}) { thread { isResolved } }", alias, varName))
			variables[varName] = id
		}

		var varDefs []string
		for name := range variables {
			varDefs = append(varDefs, fmt.Sprintf("$%s: ID!", name))
		}

		query := fmt.Sprintf("mutation(%s) {\n%s\n}", strings.Join(varDefs, ", "), strings.Join(mutations, "\n"))
		resp, err := r.executeGraphQL(query, variables)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Batch resolution error: %v\n", err)
			continue
		}

		if errors, ok := resp["errors"]; ok {
			fmt.Fprintf(os.Stderr, "Warning: GraphQL errors: %v\n", errors)
		}

		data, _ := resp["data"].(map[string]interface{})
		for j, id := range batch {
			alias := fmt.Sprintf("resolve%d", j)
			res, _ := data[alias].(map[string]interface{})
			if res != nil {
				thread, _ := res["thread"].(map[string]interface{})
				if resolved, _ := thread["isResolved"].(bool); resolved {
					fmt.Fprintf(os.Stderr, "Resolved thread %s.\n", id)
				}
			}
		}
	}
}

func (r *ReviewReposter) getPullApiURL() string {
	return fmt.Sprintf("%s/repos/%s/pulls/%s", GithubApiUrl, r.Repo, r.PullNumber)
}

func (r *ReviewReposter) fetchAndPrepareReviewMeta(reviewID string, mentionUser string) (string, string, error) {
	url := fmt.Sprintf("%s/reviews/%s", r.getPullApiURL(), reviewID)
	data, err := r.fetchJSON(url)
	if err != nil {
		return "", "", err
	}

	reviewData := data.(map[string]interface{})
	body, _ := reviewData["body"].(string)
	newBody := fmt.Sprintf("%s\n\n%s", mentionUser, body)
	state, _ := reviewData["state"].(string)
	eventMap := map[string]string{
		"APPROVED":          "APPROVE",
		"CHANGES_REQUESTED": "REQUEST_CHANGES",
	}
	event, ok := eventMap[state]
	if !ok {
		event = "COMMENT"
	}
	return newBody, event, nil
}

func (r *ReviewReposter) processSingleComment(comment map[string]interface{}, mentionUser string) map[string]interface{} {
	body, _ := comment["body"].(string)
	newBody := fmt.Sprintf("%s\n\n%s", mentionUser, body)
	newComment := map[string]interface{}{
		"path": comment["path"],
		"body": newBody,
	}

	line, ok := comment["line"].(float64)
	if !ok {
		line, _ = comment["original_line"].(float64)
	}

	position, ok := comment["position"].(float64)
	if !ok {
		position, _ = comment["original_position"].(float64)
	}

	if line != 0 {
		newComment["line"] = int(line)
		side, _ := comment["side"].(string)
		if side == "" {
			side = "RIGHT"
		}
		newComment["side"] = side
	} else if position != 0 {
		newComment["position"] = int(position)
	} else {
		fmt.Fprintf(os.Stderr, "Warning: Skipping comment on path '%v' (no line or position).\n", comment["path"])
		return nil
	}

	startLine, ok := comment["start_line"].(float64)
	if !ok {
		startLine, _ = comment["original_start_line"].(float64)
	}
	if startLine != 0 {
		newComment["start_line"] = int(startLine)
		if startSide, ok := comment["start_side"].(string); ok {
			newComment["start_side"] = startSide
		}
	}

	return newComment
}

func (r *ReviewReposter) fetchAndPrepareComments(reviewID string, mentionUser string) ([]interface{}, error) {
	url := fmt.Sprintf("%s/reviews/%s/comments", r.getPullApiURL(), reviewID)
	data, err := r.fetchJSON(url)
	if err != nil {
		return nil, err
	}

	comments, ok := data.([]interface{})
	if !ok {
		return nil, nil
	}

	var newComments []interface{}
	for _, c := range comments {
		comment := c.(map[string]interface{})
		processed := r.processSingleComment(comment, mentionUser)
		if processed != nil {
			newComments = append(newComments, processed)
		}
	}
	return newComments, nil
}

func (r *ReviewReposter) RepostReview(reviewID string, mentionUser string) error {
	newBody, event, err := r.fetchAndPrepareReviewMeta(reviewID, mentionUser)
	if err != nil {
		return err
	}

	newComments, err := r.fetchAndPrepareComments(reviewID, mentionUser)
	if err != nil {
		return err
	}

	if len(newComments) == 0 && strings.TrimSpace(newBody) == mentionUser {
		return nil
	}

	payload := map[string]interface{}{
		"body":     newBody,
		"event":    event,
		"comments": newComments,
	}
	_, err = r.postJSON(fmt.Sprintf("%s/reviews", r.getPullApiURL()), payload, "POST")
	if err != nil {
		return err
	}

	// Cleanup
	allThreads, err := r.fetchAllReviewThreads()
	if err == nil {
		threadIDsSet := make(map[string]bool)
		var reviewDbID float64
		fmt.Sscanf(reviewID, "%f", &reviewDbID)

		for _, tid := range r.filterThreadsForReview(allThreads, reviewDbID) {
			threadIDsSet[tid] = true
		}
		for _, tid := range r.filterJulesThreads(allThreads) {
			threadIDsSet[tid] = true
		}

		var threadIDs []string
		for tid := range threadIDsSet {
			threadIDs = append(threadIDs, tid)
		}
		r.resolveThreadsById(threadIDs)

		if event == "APPROVE" || event == "REQUEST_CHANGES" {
			r.dismissReview(reviewID)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Warning: Failed to fetch threads for cleanup: %v\n", err)
	}

	return nil
}
