package zanzibar

import (
	"github.com/tx7do/kratos-authz/engine/zanzibar/keto"
	"github.com/tx7do/kratos-authz/engine/zanzibar/openfga"
)

type OptFunc func(*State)

func WithKeto(readUrl, writeUrl string, useGRPC bool) OptFunc {
	return func(s *State) {
		s.ketoClient = keto.NewClient(readUrl, writeUrl, useGRPC)
	}
}

func WithOpenFga(
	apiUrl string,
	storeId string,
	token *string,
	clientId, clientSecret, apiAudience, apiTokenIssuer *string,
) OptFunc {
	return func(s *State) {
		var opts []openfga.ClientOption

		opts = append(opts, openfga.WithApiUrl(apiUrl), openfga.WithStoreId(storeId))

		if clientId != nil && clientSecret != nil && apiAudience != nil && apiTokenIssuer != nil {
			opts = append(opts, openfga.WithClientId(*clientId, *clientSecret, *apiAudience, *apiTokenIssuer))
		}
		if token != nil {
			opts = append(opts, openfga.WithToken(*token))
		}

		s.openfgaClient = openfga.NewClient(opts...)
	}
}
