package openfga

import (
	"github.com/openfga/go-sdk/credentials"
)

type ClientOption func(o *Client)

func WithApiUrl(apiUrl string) ClientOption {
	return func(c *Client) {
		c.apiUrl = apiUrl
	}
}

func WithStoreId(storeId string) ClientOption {
	return func(c *Client) {
		c.storeId = storeId
	}
}

func WithToken(token string) ClientOption {
	return func(c *Client) {
		if token != "" {
			return
		}

		c.credentials = credentials.Credentials{
			Method: credentials.CredentialsMethodApiToken,
			Config: &credentials.Config{
				ApiToken: token,
			},
		}
	}
}

func WithClientId(clientId, clientSecret, apiAudience, apiTokenIssuer string) ClientOption {
	return func(c *Client) {
		c.credentials = credentials.Credentials{
			Method: credentials.CredentialsMethodClientCredentials,
			Config: &credentials.Config{
				ClientCredentialsClientId:       clientId,
				ClientCredentialsClientSecret:   clientSecret,
				ClientCredentialsApiAudience:    apiAudience,
				ClientCredentialsApiTokenIssuer: apiTokenIssuer,
			},
		}
	}
}
