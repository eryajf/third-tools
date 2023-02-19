package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-github/v47/github"
	"golang.org/x/oauth2"
)

var (
	client      *github.Client
	GtihubToken string = "xxxxxxxxxxxxxxxxxxxxxxxxx"
)

func init() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GtihubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client = github.NewClient(tc)
}

type ArtalkComment struct {
	ID            string `json:"id"`
	Rid           string `json:"rid"`
	Content       string `json:"content"`
	Ua            string `json:"ua"`
	IP            string `json:"ip"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	IsCollapsed   string `json:"is_collapsed"`
	IsPending     string `json:"is_pending"`
	VoteUp        string `json:"vote_up"`
	VoteDown      string `json:"vote_down"`
	Nick          string `json:"nick"`
	Email         string `json:"email"`
	Link          string `json:"link"`
	Password      string `json:"password"`
	BadgeName     string `json:"badge_name"`
	BadgeColor    string `json:"badge_color"`
	PageKey       string `json:"page_key"`
	PageTitle     string `json:"page_title"`
	PageAdminOnly string `json:"page_admin_only"`
	SiteName      string `json:"site_name"`
	SiteUrls      string `json:"site_urls"`
}

func main() {
	githubUser := "eryajf"
	githubRepo := "eryajf.github.io"

	issues, err := GetAllIssue(githubUser, githubRepo)
	if err != nil {
		fmt.Println(err)
	}
	var artalks []ArtalkComment
	for _, repo := range issues {
		if *repo.Comments > 0 {
			comments, err := GetAllComment(githubUser, githubRepo, repo.GetNumber())
			if err != nil {
				fmt.Printf("get all comments failed: %v\n", err)
			}
			for _, comment := range comments {

				user, err := GetUser(comment.User.GetLogin())
				if err != nil {
					fmt.Printf("get user failed: %v\n", err)
				}
				var email string
				if user.GetEmail() == "" {
					email = "empyt@github.com"
				} else {
					email = user.GetEmail()
				}

				artalks = append(artalks, ArtalkComment{
					ID:            fmt.Sprintf("%v", repo.GetNumber()),
					Rid:           "0",
					Content:       comment.GetBody(),
					Ua:            "",
					IP:            "",
					CreatedAt:     comment.GetCreatedAt().String(),
					UpdatedAt:     comment.GetUpdatedAt().String(),
					IsCollapsed:   "false",
					IsPending:     "false",
					VoteUp:        "0", // 赞成
					VoteDown:      "0",
					Nick:          comment.User.GetLogin(),   // 评论者的昵称
					Email:         email,                     // 评论者的邮箱
					Link:          comment.User.GetHTMLURL(), // 评论者的网站
					Password:      "",
					BadgeName:     "",
					BadgeColor:    "",
					PageKey:       strings.Split(repo.GetBody(), "https://wiki.eryajf.net")[1], // 页面的url,只取域名后的uri
					PageTitle:     "数据迁移",
					PageAdminOnly: "false",
					SiteName:      "二丫讲梵",
					SiteUrls:      "https://comment.eryajf.net",
				})
			}
		}
	}
	str, err := json.Marshal(artalks)
	if err != nil {
		fmt.Printf("marshal failed: %v\n", err)
	}
	fmt.Println(string(str))
}

// 获取用户信息
func GetUser(name string) (*github.User, error) {
	ctx := context.Background()
	user, _, err := client.Users.Get(ctx, name)
	if err != nil {
		fmt.Printf("get user info failed: %v\n", err)
	}
	return user, nil
}

// 获取仓库所有的issue
func GetAllIssue(owner, repoName string) ([]*github.Issue, error) {
	ctx := context.Background()
	opt := &github.IssueListByRepoOptions{
		State:       "open",
		Labels:      []string{"Vssue"},
		ListOptions: github.ListOptions{PerPage: 10},
	}
	// get all pages of results
	var allIssues []*github.Issue
	for {
		repos, resp, err := client.Issues.ListByRepo(ctx, owner, repoName, opt)
		if err != nil {
			return nil, err
		}
		allIssues = append(allIssues, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allIssues, nil
}

// 获取对应issue的所有对话
func GetAllComment(owner, repoName string, number int) ([]*github.IssueComment, error) {
	ctx := context.Background()
	opt := &github.IssueListCommentsOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	// get all pages of results
	var allRepos []*github.IssueComment
	for {
		repos, resp, err := client.Issues.ListComments(ctx, owner, repoName, number, opt)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allRepos, nil
}

func UpdateIssue(issue *github.Issue) {
	ctx := context.Background()
	newTitle := fmt.Sprintf("[Vssue]-[Comment]-%s", strings.Split(*issue.Title, "」")[1])
	newBody := strings.Split(*issue.Body, "页面：")[1]
	newissue := github.IssueRequest{
		Title:  github.String(newTitle),
		Body:   github.String(newBody),
		Labels: &[]string{"Vssue"},
	}
	_, rsp, err := client.Issues.Edit(ctx, *issue.GetUser().Login, *issue.GetRepository().Name, issue.GetNumber(), &newissue)
	if err != nil {
		fmt.Println(err)
	}
	if rsp.StatusCode == 200 {
		fmt.Println("修改成功")
	}
}

func GetAllRepo() ([]*github.Repository, error) {
	ctx := context.Background()
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	// get all pages of results
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.List(ctx, "", opt)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allRepos, nil
}
