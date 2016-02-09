package main

import "sort"
import "github.com/google/go-github/github"

type ByName []github.Repository

func (a ByName) Len()          int  {return len(a)}
func (a ByName) Swap(i, j int)      {a[i], a[j] = a[j], a[i]}
func (a ByName) Less(i, j int) bool {return *a[i].Name < *a[j].Name}

func main() {
    client := github.NewClient(nil)
    org, _, _ := client.Organizations.Get("immobilienscout24")
    println("Listing repositories for organization: ", *org.Login)
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
    for i := 0; i < len(allRepos); i++ {
        println(*allRepos[i].Name)
    }
}
