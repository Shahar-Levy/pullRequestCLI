package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/google/go-github/github"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {

	app := cli.NewApp()

	app.Name = "Github Pull Request Lookup CLI"
	app.Usage = "Let's you look up for pull requests within your organizations"

	myFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "organization",
			Value: "",
			Usage: "Determines which organization will be checked for pull requests",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "pulls",
			Usage: "Looks up the pull requests in your organizations",
			Flags: myFlags,

			Action: func(c *cli.Context) error {

				fmt.Println(c.String("organization"))
				userLogin, client := authenticateUser()

				if len(c.String("organization")) > 1 {
					printRequestedOrgPR(userLogin, client, strings.TrimSpace(c.String("organization")))
				} else {
					printAllOrgPR(userLogin, client)
				}
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// printRequestedOrgPR prints all pull requests only from the requested organization
func printRequestedOrgPR(userLogin string, client *github.Client, requestedOrg string) {
	repos := getOrgRepos(client, requestedOrg, userLogin)
	for _, repo := range repos {
		repoName := trimQuote(github.Stringify(repo.Name))

		prs := getPullRequests(requestedOrg, repoName, client)
		for _, pr := range prs {
			fmt.Printf("%v  has made a pull request in %v titled: %v\n", github.Stringify(pr.User.Login), repoName,
				github.Stringify(pr.Title))
		}
	}
}

// printAllOrgPR prints all pull requests from every organization the user is a member of
func printAllOrgPR(userLogin string, client *github.Client) {
	orgs := getUserOrgs(userLogin, client)
	orgsRepos := makeOrgMap(orgs, client, userLogin)

	for org, repoArr := range orgsRepos {
		orgName := trimQuote(github.Stringify(org))

		for _, repo := range repoArr {
			repoName := trimQuote(github.Stringify(repo.Name))

			prs := getPullRequests(orgName, repoName, client)
			for _, pr := range prs {
				fmt.Printf("%v  has made a pull request in %v titled: %v\n", github.Stringify(pr.User.Login), repoName,
					github.Stringify(pr.Title))
			}
		}

	}
}

// authenticateUser uses Github's auth to authenticate with user credentials
func authenticateUser() (string, *github.Client) {
	r := bufio.NewReader(os.Stdin)
	fmt.Print("GitHub Username: ")
	username, _ := r.ReadString('\n')

	fmt.Print("GitHub Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)

	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}

	client := github.NewClient(tp.Client())
	ctx := context.Background()
	user, _, err := client.Users.Get(ctx, "")

	if _, ok := err.(*github.TwoFactorAuthError); ok {
		fmt.Print("\nGitHub OTP: ")
		otp, _ := r.ReadString('\n')
		tp.OTP = strings.TrimSpace(otp)
		user, _, err = client.Users.Get(ctx, "")
	}
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
	}
	fmt.Println("")
	userLogin := trimQuote(github.Stringify(user.Login))
	return userLogin, client

}

// getUserOrgs gets all user's public organgizations
func getUserOrgs(userLogin string, client *github.Client) []*github.Organization {
	orgOpt := &github.ListOptions{}
	orgs, _, err := client.Organizations.List(context.Background(), userLogin, orgOpt)
	if err != nil {
		log.Fatal(err)
	}
	return orgs
}

// makeOrgMap maps a user's organizations to the repositories in it
func makeOrgMap(orgs []*github.Organization, client *github.Client, userLogin string) map[string][]*github.Repository {
	orgsRepos := make(map[string][]*github.Repository)
	for _, org := range orgs {
		currentOrg := trimQuote(github.Stringify(org.Login))
		orgsRepos[currentOrg] = getOrgRepos(client, trimQuote(github.Stringify(org.GetLogin())), userLogin)
	}
	return orgsRepos
}

// getOrgRepos gets all of an organization's repositories
func getOrgRepos(client *github.Client, org string, userLogin string) []*github.Repository {
	repoOpt := &github.RepositoryListByOrgOptions{}
	// repos, _, err := client.Repositories.List(context.Background(), userLogin, repoOpt)
	repos, _, err := client.Repositories.ListByOrg(context.Background(), org, repoOpt)
	if err != nil {
		log.Fatal(err)
	}
	return repos
}

// getPullRequests gets pull requests in a repository
func getPullRequests(org string, repo string, client *github.Client) []*github.PullRequest {

	prOpt := &github.PullRequestListOptions{}
	prs, _, err := client.PullRequests.List(context.Background(), org, repo, prOpt) //"testingGithubAPI"
	if err != nil {
		log.Fatal(err)
	}
	return prs
}

func trimQuote(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}
