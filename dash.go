package main

import "io/ioutil"
import "sort"
import "github.com/google/go-github/github"
import "golang.org/x/oauth2"

type ByName []github.Repository

func (a ByName) Len()          int  {return len(a)}
func (a ByName) Swap(i, j int)      {a[i], a[j] = a[j], a[i]}
func (a ByName) Less(i, j int) bool {return *a[i].Name < *a[j].Name}

func getAllRepos(client *github.Client, org *github.Organization) []github.Repository {
    opt := &github.RepositoryListByOrgOptions{ListOptions: github.ListOptions{PerPage: 10},}
    var allRepos []github.Repository
    for {
        repos, resp, err := client.Repositories.ListByOrg(*org.Login, opt)
        if err != nil {
            println("An error was reveiced while fetching")
            break
        }
        allRepos = append(allRepos, repos...)
        if resp.NextPage == 0 {
            break
        }
        opt.ListOptions.Page = resp.NextPage
    }
    sort.Sort(ByName(allRepos))
    return allRepos
}

func readTokenFromFile() string {
    dat, _ := ioutil.ReadFile("token")
    return string(dat)
}

func main() {
    var token string = readTokenFromFile()
    ts := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: token},
    )
    tc := oauth2.NewClient(oauth2.NoContext, ts)
    client := github.NewClient(tc)
    org, _, err := client.Organizations.Get("immobilienscout24")
    if err != nil {
        println("Failed to fetch organization from github")
        println(err.Error())
        return
    }
    println("Listing repositories for organization: ", *org.Login)

    var allRepos []github.Repository = getAllRepos(client, org)
    for i := 0; i < len(allRepos); i++ {
        println(*allRepos[i].Name)
    }
}
