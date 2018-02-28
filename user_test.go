package graylog

import (
	"reflect"
	"testing"
)

func dummyAdmin() *User {
	return &User{
		Id:          "local:admin",
		Username:    "admin",
		Email:       "",
		FullName:    "Administrator",
		Permissions: []string{"*"},
		Preferences: &Preferences{
			UpdateUnfocussed:  false,
			EnableSmartSearch: true,
		},
		Timezone:         "UTC",
		SessionTimeoutMs: 28800000,
		External:         false,
		Startpage:        nil,
		Roles:            []string{"Admin"},
		ReadOnly:         true,
		SessionActive:    true,
		LastActivity:     "2018-02-21T07:35:45.926+0000",
		ClientAddress:    "172.18.0.1",
	}
}

func TestCreateUser(t *testing.T) {
	server, client, err := getServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	admin := dummyAdmin()
	if err := client.CreateUser(admin); err != nil {
		t.Fatal("Failed to CreateUser", err)
	}
	if err := client.CreateUser(admin); err == nil {
		t.Fatal("the user name must be unique ")
	}
}

func TestGetUsers(t *testing.T) {
	server, client, err := getServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	admin := dummyAdmin()
	server.Users[admin.Username] = *admin
	users, err := client.GetUsers()
	if err != nil {
		t.Fatal("Failed to GetUsers", err)
	}
	exp := []User{*admin}
	if !reflect.DeepEqual(users, exp) {
		t.Fatalf("client.GetUsers() == %v, wanted %v", users, exp)
	}
}

func TestGetUser(t *testing.T) {
	server, client, err := getServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	exp := dummyAdmin()
	server.Users[exp.Username] = *exp
	user, err := client.GetUser(exp.Username)
	if err != nil {
		t.Fatal("Failed to GetUser", err)
	}
	if !reflect.DeepEqual(*user, *exp) {
		t.Fatalf("client.GetUser() == %v, wanted %v", user, exp)
	}
	if _, err := client.GetUser(""); err == nil {
		t.Fatal("username should be required.")
	}
}

func TestUpdateUser(t *testing.T) {
	server, client, err := getServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	user := dummyAdmin()
	server.Users[user.Username] = *user
	user.FullName = "changed!"
	if err := client.UpdateUser(user.Username, user); err != nil {
		t.Fatal("Failed to UpdateUser", err)
	}
}

func TestDeleteUser(t *testing.T) {
	server, client, err := getServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	user := dummyAdmin()
	server.Users[user.Username] = *user
	err = client.DeleteUser(user.Username)
	if err != nil {
		t.Fatal("Failed to DeleteUser", err)
	}
}
