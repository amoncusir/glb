package unres

import (
	"fmt"
	"strings"
)

func ExtractFromUri(uri string) (*Uri, error) {
	u := &Uri{}
	err := fmt.Errorf("invalid uri: %s must follow the format <scheme>:[//[<userinfo>@]<host>[]:<port>]]<path>[?<query>][#<fragment>]", uri)

	if chunks := strings.SplitN(uri, ":", 2); len(chunks) != 2 {
		return nil, err
	} else {
		u.Scheme = chunks[0]
		uri = chunks[1]
	}

	if strings.HasPrefix(uri, "//") {
		uri = uri[2:]
		var auth string

		if i := strings.Index(uri, "/"); i >= 0 {
			// Has Path
			auth = uri[:i]
			uri = uri[i:]
		} else if i := strings.Index(uri, "?"); i >= 0 {
			// No path, has query
			auth = uri[:i]
			uri = uri[i:]
		} else if i := strings.Index(uri, "#"); i >= 0 {
			// No path, has fragment
			auth = uri[:i]
			uri = uri[i:]
		} else {
			// Nothing more
			auth = uri
			uri = ""
		}

		if b, a, found := strings.Cut(auth, "@"); found {
			u.Autority.Userinfo = b
			auth = a
		}

		if len(auth) > 0 && auth[0] == '[' {
			// Ipv6 host
			i := strings.Index(auth, "]")

			if i == -1 {
				return nil, err
			}

			u.Autority.Host = auth[0 : i+1]

			if len(auth) > i+1 {
				u.Autority.Port = auth[i+2:]
			}

		} else {
			if b, a, found := strings.Cut(auth, ":"); found {
				u.Autority.Port = a
				auth = b
			}

			u.Autority.Host = auth
		}

	}

	if b, a, found := strings.Cut(uri, "#"); found {
		u.Fragment = a
		uri = b
	}

	if b, a, found := strings.Cut(uri, "?"); found {
		u.Query = a
		uri = b
	}

	u.Path = uri

	return u, nil
}
