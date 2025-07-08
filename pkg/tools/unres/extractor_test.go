package unres_test

import (
	"amoncusir/example/pkg/tools/unres"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type extractorTestCase struct {
	uri                                                 string
	Scheme, Path, Query, Fragment, Userinfo, Host, Port string
}

func TestValidExtractFromUriValues(t *testing.T) {
	testCases := []extractorTestCase{
		{"://example.test", "", "", "", "", "", "example.test", ""},
		{"tcp://127.0.0.1", "tcp", "", "", "", "", "127.0.0.1", ""},
		{"tcp://127.0.0.1?hello", "tcp", "", "hello", "", "", "127.0.0.1", ""},
		{"tcp://127.0.0.1#hello", "tcp", "", "", "hello", "", "127.0.0.1", ""},
		{"tcp://127.0.0.1:9090/", "tcp", "/", "", "", "", "127.0.0.1", "9090"},

		{"tcp://user:pwd@127.0.0.1:9090/", "tcp", "/", "", "", "user:pwd", "127.0.0.1", "9090"},
		{"https://john.doe@www.example.com:1234/forum/questions/?tag=networking&order=newest#:~:text=whatever",
			"https", "/forum/questions/", "tag=networking&order=newest", ":~:text=whatever", "john.doe", "www.example.com", "1234"},

		{"ldap://[2001:db8::7]/c=GB?objectClass?one", "ldap", "/c=GB", "objectClass?one", "", "", "[2001:db8::7]", ""},
		{"ldap://[2001:db8::7]:8080/c=GB?objectClass?one", "ldap", "/c=GB", "objectClass?one", "", "", "[2001:db8::7]", "8080"},
		{"ldap://test:one@[2001:db8::7]:8080/c=GB?objectClass?one", "ldap", "/c=GB", "objectClass?one", "", "test:one", "[2001:db8::7]", "8080"},
		{"ldap://test:one@[2001:db8::7]/c=GB?objectClass?one", "ldap", "/c=GB", "objectClass?one", "", "test:one", "[2001:db8::7]", ""},

		{"mailto:John.Doe@example.com", "mailto", "John.Doe@example.com", "", "", "", "", ""},
		{"news:comp.infosystems.www.servers.unix", "news", "comp.infosystems.www.servers.unix", "", "", "", "", ""},
		{"urn:oasis:names:specification:docbook:dtd:xml:4.1.2", "urn", "oasis:names:specification:docbook:dtd:xml:4.1.2", "", "", "", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.uri, func(t *testing.T) {
			assert := assert.New(t)
			uri, err := unres.ExtractFromUri(tc.uri)

			if err != nil {
				t.Error(err)
			}

			assert.Equal(tc.Scheme, uri.Scheme)
			assert.Equal(tc.Path, uri.Path)
			assert.Equal(tc.Query, uri.Query)
			assert.Equal(tc.Fragment, uri.Fragment)

			assert.Equal(tc.Userinfo, uri.Autority.Userinfo)
			assert.Equal(tc.Host, uri.Autority.Host)
			assert.Equal(tc.Port, uri.Autority.Port)
		})
	}
}

func TestInvalidExtractFromUriValues(t *testing.T) {
	testCases := []string{
		"/hello",
		"invalid",
		"0",
	}

	for _, v := range testCases {
		t.Run(fmt.Sprintf("Invalid URI: %s", v), func(t *testing.T) {
			uri, err := unres.ExtractFromUri(v)

			if err == nil {
				t.Errorf("invalid uri must return an error: %s", uri)
			}
		})
	}
}
