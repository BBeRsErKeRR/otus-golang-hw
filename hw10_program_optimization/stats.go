package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/buger/jsonparser"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type UsersWithErr struct {
	Users chan User
	Err   chan error
}

func parseJSONUser(input []byte) (User, error) {
	var (
		user User
		err  error
	)
	err = jsonparser.ObjectEach(input,
		func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			switch string(key) {
			case "Id":
				v, err := jsonparser.ParseInt(value)
				if err != nil {
					return err
				}
				user.ID = int(v)
			case "Name":
				user.Name, err = jsonparser.ParseString(value)
				if err != nil {
					return err
				}
			case "Username":
				user.Username, err = jsonparser.ParseString(value)
				if err != nil {
					return err
				}
			case "Email":
				user.Email, err = jsonparser.ParseString(value)
				if err != nil {
					return err
				}
			case "Phone":
				user.Phone, err = jsonparser.ParseString(value)
				if err != nil {
					return err
				}
			case "Password":
				user.Password, err = jsonparser.ParseString(value)
				if err != nil {
					return err
				}
			case "Address":
				user.Address, err = jsonparser.ParseString(value)
				if err != nil {
					return err
				}
			}
			return nil
		})
	return user, err
}

func getUsers(r io.Reader) (result UsersWithErr, err error) { //nolint:unparam
	result.Users = make(chan User)
	result.Err = make(chan error, 1)
	go func() {
		defer close(result.Users)
		fileScanner := bufio.NewScanner(r)
		fileScanner.Split(bufio.ScanLines)
		for fileScanner.Scan() {
			user, err := parseJSONUser(fileScanner.Bytes())
			if err != nil {
				result.Err <- err
				return
			}
			result.Users <- user
		}
	}()
	return result, nil
}

func countDomains(usersWithErr UsersWithErr, domain string) (DomainStat, error) {
	result := make(DomainStat)
	re, err := regexp.Compile("\\." + domain)
	if err != nil {
		return result, err
	}

	for {
		select {
		case err := <-usersWithErr.Err:
			return nil, err
		case u, ok := <-usersWithErr.Users:
			select {
			case err := <-usersWithErr.Err:
				return nil, err
			default:
				if re.Match([]byte(u.Email)) {
					result[strings.ToLower(strings.SplitN(u.Email, "@", 2)[1])]++
				}
			}
			if !ok {
				return result, nil
			}
		}
	}
}
