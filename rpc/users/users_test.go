package users

import (
	"testing"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"io"
	"github.com/foofilers/cfhd/util"
)

var Conn *grpc.ClientConn
var createdUser *User

func TestMain(m *testing.M) {
	var err error
	Conn, err = grpc.Dial("localhost:50051", grpc.WithInsecure())
	defer Conn.Close()
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestAddUser(t *testing.T) {
	userCl := NewUsersClient(Conn)
	user := &User{
		Username:util.RandStringRunes(8),
		Email:"n3wtron@gmail.com",
	}
	insertedUser, err := userCl.Add(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}
	if (insertedUser.Username != user.Username) {
		t.Fatal("Username not match")
	}
	if len(insertedUser.Id) == 0 {
		t.Fatal("UserID should be valorized")
	}
	createdUser = insertedUser
}

func TestListUser(t *testing.T) {
	if createdUser == nil {
		t.Skip("previous dependency test [addUser] doesn't pass")
	}
	userCl := NewUsersClient(Conn)
	lstCl, err := userCl.List(context.TODO(), &UserListRequest{})
	if err != nil {
		t.Fatal(err)
	}
	for {
		user, err := lstCl.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("%v.ListFeatures(_) = _, %v", userCl, err)
		}
		if user.Id == createdUser.Id {
			return
		}
	}
	t.Fatalf("Created user {%+v} not found", createdUser)
}


func TestDeleteUser(t *testing.T) {
	if createdUser == nil {
		t.Skip("previous dependency test [addUser] doesn't pass")
	}
	userCl := NewUsersClient(Conn)
	if _, err := userCl.Delete(context.TODO(), &DeleteRequest{Id:createdUser.Id}); err != nil {
		t.Fatal(err)
	}
}