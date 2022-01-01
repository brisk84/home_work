package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	var (
		user User
		line []byte
		err  error
	)
	stat := make(DomainStat)
	reader := bufio.NewReader(r)
	dotDomain := "." + domain

	for {
		line, _, err = reader.ReadLine()
		if (err != nil) && (!errors.Is(err, io.EOF)) {
			return stat, err
		}
		err = user.UnmarshalJSON(line)
		if errors.Is(err, io.EOF) {
			return stat, nil
		}
		if err != nil {
			return stat, err
		}

		if strings.HasSuffix(user.Email, dotDomain) {
			x := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			stat[x]++
		}
	}
}
