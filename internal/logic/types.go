package logic

import (
	"fmt"
	"time"

	"github.com/swayrider/grpcclients/authclient"
)

type ServiceClient struct {
	name        string
	description string
	clientId    string
	scopes      []string
}

func (s ServiceClient) Name() string {
	return s.name
}

func (s ServiceClient) Description() string {
	return s.description
}

func (s ServiceClient) ClientId() string {
	return s.clientId
}

func (s ServiceClient) Scopes() []string {
	return s.scopes
}

func NewServiceClient(
	clientId string,
	name string,
	description string,
	scopes ...string,
) authclient.ServiceClient {
	return &ServiceClient{
		name:        name,
		description: description,
		clientId:    clientId,
		scopes:      scopes,
	}
}

type User struct {
	userId      string
	email       string
	isVerified  bool
	isAdmin     bool
	accountType string
}

func (u User) UserId() string {
	return u.userId
}

func (u User) Email() string {
	return u.email
}

func (u User) IsVerified() bool {
	return u.isVerified
}

func (u User) IsAdmin() bool {
	return u.isAdmin
}

func (u User) AccountType() string {
	return u.accountType
}

func (u User) Display() {
	fmt.Printf(
		"\tEmail: %s\n\tUserID: %s\n\tAccount Type: %s\n\tVerified: %t\n\tAdmin: %t\n",
		u.email, u.userId, u.accountType, u.isVerified, u.isAdmin)
}

type Invite struct {
	id        string
	email     string
	createdAt time.Time
}

func (i Invite) Id() string {
	return i.id
}

func (i Invite) Email() string {
	return i.email
}

func (i Invite) CreatedAt() time.Time {
	return i.createdAt
}

func NewInvite(id string, email string, createdAt time.Time) authclient.Invite {
	return &Invite{id: id, email: email, createdAt: createdAt}
}

func NewUser(
	userId string,
	email string,
	isVerified bool,
	isAdmin bool,
	accountType string,
) authclient.User {
	return &User{
		userId:      userId,
		email:       email,
		isVerified:  isVerified,
		isAdmin:     isAdmin,
		accountType: accountType,
	}
}
