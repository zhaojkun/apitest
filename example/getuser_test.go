package main

import (
	"time"

	"github.com/octokit/go-octokit/octokit"

	"github.com/seesawlabs/testilla"
)

// TODO: mock Github API call
type GetUserTest struct {
}

func (t *GetUserTest) Method() string      { return "GET" }
func (t *GetUserTest) Description() string { return "Test for GetUser API handler" }
func (t *GetUserTest) Path() string        { return "user/{username}" }

func (t *GetUserTest) TestCases() []testilla.ApiTestCase {
	elgrisCreatedAt := time.Date(2012, time.June, 29, 11, 57, 38, 0, time.UTC)
	elgrisUpdatedAt := time.Date(2015, time.December, 27, 19, 33, 41, 0, time.UTC)

	return []testilla.ApiTestCase{
		{
			Description: "Successful getting of user details",
			PathParams: testilla.ParamMap{
				"username": testilla.Param{Value: "elgris"},
			},

			ExpectedHttpCode: 200,
			ExpectedData: octokit.User{
				AvatarURL:         "https://avatars.githubusercontent.com/u/1905821?v=3",
				Blog:              "http://elgris-blog.blogspot.com/",
				CreatedAt:         &elgrisCreatedAt,
				UpdatedAt:         &elgrisUpdatedAt,
				EventsURL:         "https://api.github.com/users/elgris/events{/privacy}",
				Followers:         10,
				FollowersURL:      "https://api.github.com/users/elgris/followers",
				Following:         3,
				FollowingURL:      "https://api.github.com/users/elgris/following{/other_user}",
				GistsURL:          "https://api.github.com/users/elgris/gists{/gist_id}",
				Hireable:          true,
				HTMLURL:           "https://github.com/elgris",
				ID:                1905821,
				Location:          "Saint Petersburg, Russia",
				Login:             "elgris",
				Name:              "elgris",
				OrganizationsURL:  "https://api.github.com/users/elgris/orgs",
				PublicRepos:       24,
				ReceivedEventsURL: "https://api.github.com/users/elgris/received_events",
				ReposURL:          "https://api.github.com/users/elgris/repos",
				StarredURL:        "https://api.github.com/users/elgris/starred{/owner}{/repo}",
				SubscriptionsURL:  "https://api.github.com/users/elgris/subscriptions",
				Type:              "User",
				URL:               "https://api.github.com/users/elgris",
			},
		},
		{
			Description: "404 error in case user not found",
			PathParams: testilla.ParamMap{
				"username": testilla.Param{Value: "someveryunknown"},
			},

			ExpectedHttpCode: 404,
			ExpectedData:     []byte("user someveryunknown not found"),
		},
		{
			Description: "500 error in case something bad happens",
			PathParams: testilla.ParamMap{
				"username": testilla.Param{Value: "BadGuy"},
			},

			ExpectedHttpCode: 500,
			ExpectedData:     []byte("BadGuy failed me :("),
		},
	}
}
