package apitest

import "time"

type HelloTest struct{}

func (t *HelloTest) Method() string      { return "GET" }
func (t *HelloTest) Description() string { return "Test for HelloWorld API handler" }
func (t *HelloTest) Path() string        { return "/hello" }
func (t *HelloTest) TestCases() []ApiTestCase {
	return []ApiTestCase{
		{
			Description:      "Successful greeting of the world",
			ExpectedHttpCode: 200,
			ExpectedData:     "Hello World!",
		},
	}
}

type User struct {
	SiteAdmin         bool       `json:"site_admin,omitempty"`
	Login             string     `json:"login,omitempty"`
	ID                int        `json:"id,omitempty"`
	AvatarURL         string     `json:"avatar_url,omitempty"`
	GravatarID        string     `json:"gravatar_id,omitempty"`
	URL               string     `json:"url,omitempty"`
	Name              string     `json:"name,omitempty"`
	Company           string     `json:"company,omitempty"`
	Blog              string     `json:"blog,omitempty"`
	Location          string     `json:"location,omitempty"`
	Email             string     `json:"email,omitempty"`
	Hireable          bool       `json:"hireable,omitempty"`
	Bio               string     `json:"bio,omitempty"`
	PublicRepos       int        `json:"public_repos,omitempty"`
	Followers         int        `json:"followers,omitempty"`
	Following         int        `json:"following,omitempty"`
	HTMLURL           string     `json:"html_url,omitempty"`
	CreatedAt         *time.Time `json:"created_at,omitempty"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
	Type              string     `json:"type,omitempty"`
	FollowingURL      string     `json:"following_url,omitempty"`
	FollowersURL      string     `json:"followers_url,omitempty"`
	GistsURL          string     `json:"gists_url,omitempty"`
	StarredURL        string     `json:"starred_url,omitempty"`
	SubscriptionsURL  string     `json:"subscriptions_url,omitempty"`
	OrganizationsURL  string     `json:"organizations_url,omitempty"`
	ReposURL          string     `json:"repos_url,omitempty"`
	EventsURL         string     `json:"events_url,omitempty"`
	ReceivedEventsURL string     `json:"received_events_url,omitempty"`
}

type GetUserTest struct{}

func (t *GetUserTest) Method() string      { return "GET" }
func (t *GetUserTest) Description() string { return "Test for GetUser API handler" }
func (t *GetUserTest) Path() string        { return "/user/{username}" }

func (t *GetUserTest) TestCases() []ApiTestCase {
	return []ApiTestCase{
		{
			Description: "Successful getting of user details",
			PathParams: ParamMap{
				"username": Param{Value: "octocat"},
			},

			ExpectedHttpCode: 200,
			ExpectedData: User{
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
			PathParams: ParamMap{
				"username": Param{Value: "someveryunknown"},
			},

			ExpectedHttpCode: 404,
			ExpectedData:     "user someveryunknown not found",
		},
		{
			Description: "500 error in case something bad happens",
			PathParams: ParamMap{
				"username": Param{Value: "BadGuy"},
			},

			ExpectedHttpCode: 500,
			ExpectedData:     "BadGuy failed me :(",
		},
	}
}
