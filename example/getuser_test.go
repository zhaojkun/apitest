package main

import (
	"github.com/octokit/go-octokit/octokit"

	"github.com/seesawlabs/apitest"
)

// TODO: mock Github API call
type GetUserTest struct {
}

func (t *GetUserTest) Method() string      { return "GET" }
func (t *GetUserTest) Description() string { return "Test for GetUser API handler" }
func (t *GetUserTest) Path() string        { return "user/{username}" }

func (t *GetUserTest) TestCases() []apitest.ApiTestCase {
	return []apitest.ApiTestCase{
		{
			Description: "Successful getting of user details",
			PathParams: apitest.ParamMap{
				"username": apitest.Param{Value: "octocat"},
			},

			ExpectedHttpCode: 200,
			ExpectedData: octokit.User{
				EventsURL:         "https://api.github.com/users/octocat/events{/privacy}",
				Followers:         20,
				FollowersURL:      "https://api.github.com/users/octocat/followers",
				Following:         0,
				FollowingURL:      "https://api.github.com/users/octocat/following{/other_user}",
				GistsURL:          "https://api.github.com/users/octocat/gists{/gist_id}",
				Hireable:          false,
				HTMLURL:           "https://github.com/octocat",
				Location:          "San Francisco",
				Login:             "octocat",
				Name:              "monalisa octocat",
				OrganizationsURL:  "https://api.github.com/users/octocat/orgs",
				PublicRepos:       2,
				ReceivedEventsURL: "https://api.github.com/users/octocat/received_events",
				ReposURL:          "https://api.github.com/users/octocat/repos",
				StarredURL:        "https://api.github.com/users/octocat/starred{/owner}{/repo}",
				SubscriptionsURL:  "https://api.github.com/users/octocat/subscriptions",
				Type:              "User",
				URL:               "https://api.github.com/users/octocat",
			},
		},
		{
			Description: "404 error in case user not found",
			PathParams: apitest.ParamMap{
				"username": apitest.Param{Value: "someveryunknown"},
			},

			ExpectedHttpCode: 404,
			ExpectedData:     []byte("user someveryunknown not found"),
		},
		{
			Description: "500 error in case something bad happens",
			PathParams: apitest.ParamMap{
				"username": apitest.Param{Value: "BadGuy"},
			},

			ExpectedHttpCode: 500,
			ExpectedData:     []byte("BadGuy failed me :("),
		},
	}
}
