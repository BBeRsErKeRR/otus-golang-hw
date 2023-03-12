package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/valyala/fastjson"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type usersEmails [100_000]string

func getUsers(r io.Reader) (result usersEmails, err error) {
	fileScanner := bufio.NewScanner(r)
	fileScanner.Split(bufio.ScanLines)
	var (
		i int
	)
	var p fastjson.Parser
	for fileScanner.Scan() {
		var v *fastjson.Value
		if v, err = p.ParseBytes(fileScanner.Bytes()); err != nil {
			return
		}
		result[i] = string(v.GetStringBytes("Email"))
		i++
	}
	return
}

func countDomains(u usersEmails, domain string) (DomainStat, error) {
	result := make(DomainStat)
	re, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}

	for _, uE := range u {
		if re.Match([]byte(uE)) {
			result[strings.ToLower(strings.SplitN(uE, "@", 2)[1])]++
		}
	}
	return result, nil
}
