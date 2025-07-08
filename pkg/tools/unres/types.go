package unres

import (
	"fmt"
	"strings"
)

type Uri struct {
	Scheme, Path, Query, Fragment string
	Autority                      Authority
}

func (u *Uri) String() string {
	b := strings.Builder{}

	if u.Scheme != "" {
		if _, err := b.WriteString(fmt.Sprintf("%s:", u.Scheme)); err != nil {
			return ""
		}
	}

	if u.Autority.String() != "" {

		if _, err := b.WriteString(fmt.Sprintf("//%s", u.Autority)); err != nil {
			return ""
		}
	}

	if _, err := b.WriteString(u.Path); err != nil {
		return ""
	}

	if u.Query != "" {

		if _, err := b.WriteString(fmt.Sprintf("?%s", u.Query)); err != nil {
			return ""
		}
	}

	if u.Fragment != "" {

		if _, err := b.WriteString(fmt.Sprintf("#%s", u.Fragment)); err != nil {
			return ""
		}
	}

	return b.String()
}

type Authority struct {
	Userinfo, Host, Port string
}

func (a *Authority) String() string {
	b := strings.Builder{}

	if a.Userinfo != "" {

		if _, err := b.WriteString(fmt.Sprintf("%s@", a.Userinfo)); err != nil {
			return ""
		}
	}

	if _, err := b.WriteString(a.Host); err != nil {
		return ""
	}

	if a.Port != "" {

		if _, err := b.WriteString(fmt.Sprintf(":%s", a.Port)); err != nil {
			return ""
		}
	}

	return b.String()
}
