package web

import (
	"net/url"
)

type URLBuilder struct {
	baseUrl *url.URL
}

func NewURLBuilder(baseUrl *url.URL) URLBuilder {
	return URLBuilder{
		baseUrl: baseUrl,
	}
}

func (b URLBuilder) BuildAuthCallbackURL(token string) *url.URL {
	params := url.Values{
		tokenQueryParam: []string{token},
	}

	newUrl := b.baseUrl.JoinPath(githubAuthPath)
	newUrl.RawQuery = params.Encode()
	return newUrl
}

func (b URLBuilder) BuildWebhookURL(token string) *url.URL {
	params := url.Values{
		tokenQueryParam: []string{token},
	}

	newUrl := b.baseUrl.JoinPath(githubWebHookPath)
	newUrl.RawQuery = params.Encode()
	return newUrl
}
