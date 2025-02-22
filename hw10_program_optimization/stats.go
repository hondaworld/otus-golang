package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type UserEmail struct {
	Email string `json:"Email"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getEmails(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type emails [100_000]UserEmail

func getEmails(r io.Reader) (result emails, err error) {
	i := 0
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		if i >= len(result) {
			err = fmt.Errorf("превышен лимит пользователей")
			return
		}

		var userEmail UserEmail
		line := scanner.Bytes()
		if err = json.Unmarshal(line, &userEmail); err != nil {
			return
		}
		result[i] = userEmail
		i++
	}

	if err = scanner.Err(); err != nil {
		return
	}

	return
}

func countDomains(e emails, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, userEmail := range e {
		email := strings.ToLower(userEmail.Email)

		if strings.HasSuffix(email, "."+domain) {
			atIndex := strings.LastIndex(email, "@")
			if atIndex == -1 {
				continue
			}

			emailDomain := email[atIndex+1:]
			result[emailDomain]++
		}
	}
	return result, nil
}
