package logic

import (
	"errors"
	"strings"

	"github.com/swayrider/grpcclients/authclient"
)

func newAuthClient(host string, port int) (*authclient.Client, error) {
	clnt, err := authclient.New(func() (string, int) { return host, port })
	if err != nil {
		return nil, err
	}
	return clnt.(*authclient.Client), nil
}

func CheckPasswordStrength(
	authHost string,
	authPort int,
	password string,
) (ok bool, message string, err error) {
	client, err := newAuthClient(authHost, authPort)
	if err != nil {
		return
	}
	defer client.Close()

	ok, message, err = client.CheckPasswordStrength(password)
	return
}

func CreateAdmin(
	authHost string,
	authPort int,
	user string,
	password string,
	adminEmail string,
	adminPassword string,
) (newUser *User, err error) {
	client, err := newAuthClient(authHost, authPort)
	if err != nil {
		return
	}
	defer client.Close()

	adminAccessToken, _, err := client.Login(user, password, false)
	if err != nil {
		return
	}

	userId, _, err := client.CreateAdmin(
		adminAccessToken,
		adminEmail,
		adminPassword,
	)
	if err != nil {
		return
	}

	u, err := client.WhoIs(
		adminAccessToken,
		authclient.WhoIs_UserId(userId),
		NewUser)
	if err != nil {
		return
	}

	newUser, ok := u.(*User)
	if !ok {
		err = errors.New("failed to create admin user")
		return
	}
	return
}

func CreateUser(
	authHost string,
	authPort int,
	user string,
	password string,
	userEmail string,
	userPassword string,
	setVerified bool,
	accountType string,
) (newUser *User, err error) {
	client, err := newAuthClient(authHost, authPort)
	if err != nil {
		return
	}
	defer client.Close()

	adminAccessToken, _, err := client.Login(user, password, false)
	if err != nil {
		return
	}

	userId, _, err := client.Register(
		userEmail,
		userPassword,
		"", // verificationUrl not needed for CLI
	)
	if err != nil {
		return
	}

	if setVerified {
		var userAccessToken, token string

		userAccessToken, _, err = client.Login(userEmail, userPassword, false)
		if err != nil {
			return
		}

		_, token, _, err = client.CreateVerificationToken(userAccessToken)
		if err != nil {
			return
		}

		var valid bool
		valid, err = client.CheckVerificationToken(userId, token)
		if err != nil {
			return
		}

		if !valid {
			err = errors.New("failed to verify user")
			return
		}
	}

	if accountType != "" {
		_, err = client.ChangeAccountType(
			adminAccessToken, userId, accountType)
		if err != nil {
			return
		}
	}

	u, err := client.WhoIs(
		adminAccessToken,
		authclient.WhoIs_UserId(userId),
		NewUser)
	if err != nil {
		return
	}
	newUser = u.(*User)
	return
}

func CreateServiceClient(
	authHost string,
	authPort int,
	user string,
	password string,
	name string,
	description string,
	scopes []string,
) (clientId, clientSecret string, err error) {
	client, err := newAuthClient(authHost, authPort)
	if err != nil {
		return
	}
	defer client.Close()

	accessToken, _, err := client.Login(user, password, false)
	if err != nil {
		return
	}

	clientId, clientSecret, err = client.CreateServiceClient(
		accessToken,
		name,
		description,
		scopes,
	)
	return
}

func ListServiceClients(
	authHost string,
	authPort int,
	user string,
	password string,
	page int,
	pageSize int,
) (list []*ServiceClient, err error) {
	client, err := newAuthClient(authHost, authPort)
	if err != nil {
		return
	}
	defer client.Close()

	accessToken, _, err := client.Login(user, password, false)
	if err != nil {
		return
	}

	l, _, err := client.ListServiceClients(
		accessToken,
		page, pageSize,
		NewServiceClient)
	if err != nil {
		return
	}

	list = make([]*ServiceClient, 0, len(l))
	for _, sc := range l {
		list = append(list, sc.(*ServiceClient))
	}
	return
}

func WhoAmI(
	authHost string,
	authPort int,
	user string,
	password string,
) (me *User, err error) {
	client, err := newAuthClient(authHost, authPort)
	if err != nil {
		return
	}
	defer client.Close()

	accessToken, _, err := client.Login(user, password, false)
	if err != nil {
		return
	}

	u, err := client.WhoAmI(accessToken, NewUser)
	if err != nil {
		return
	}

	me, ok := u.(*User)
	if !ok {
		err = errors.New("failed to get user")
		return
	}
	return
}

func GetUser(
	authHost string,
	authPort int,
	user string,
	password string,
	identifier string,
) (found *User, err error) {
	client, err := newAuthClient(authHost, authPort)
	if err != nil {
		return
	}
	defer client.Close()

	accessToken, _, err := client.Login(user, password, false)
	if err != nil {
		return
	}

	var filter authclient.WhoIsOneOf
	if strings.Contains(identifier, "@") {
		filter = authclient.WhoIs_Email(identifier)
	} else {
		filter = authclient.WhoIs_UserId(identifier)
	}

	u, err := client.WhoIs(accessToken, filter, NewUser)
	if err != nil {
		return
	}

	found, ok := u.(*User)
	if !ok {
		err = errors.New("failed to get user")
		return
	}
	return
}

func ChangePassword(
	authHost string,
	authPort int,
	user string,
	password string,
	newPassword string,
) (message string, err error) {
	client, err := newAuthClient(authHost, authPort)
	if err != nil {
		return
	}
	defer client.Close()

	accessToken, _, err := client.Login(user, password, false)
	if err != nil {
		return
	}

	message, err = client.ChangePassword(accessToken, password, newPassword)
	return
}

func DeleteServiceClient(
	authHost string,
	authPort int,
	user string,
	password string,
	clientId string,
) (err error) {
	client, err := newAuthClient(authHost, authPort)
	if err != nil {
		return
	}
	defer client.Close()

	accessToken, _, err := client.Login(user, password, false)
	if err != nil {
		return
	}

	_, err = client.DeleteServiceClient(accessToken, clientId)
	return
}
