package main

import (
	"fmt"
	"log"

	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
	"time"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	/*if _, err = cli.RoleAdd(context.TODO(), "root"); err != nil {
		log.Fatal(err)
	}
	if _, err = cli.UserAdd(context.TODO(), "root", "123"); err != nil {
		log.Fatal(err)
	}
	if _, err = cli.UserGrantRole(context.TODO(), "root", "root"); err != nil {
		log.Fatal(err)
	}

	if _, err = cli.RoleAdd(context.TODO(), "r"); err != nil {
		log.Fatal(err)
	}

	if _, err = cli.RoleGrantPermission(
		context.TODO(),
		"r", // role name
		"foo", // key
		"zoo", // range end
		clientv3.PermissionType(clientv3.PermReadWrite),
	); err != nil {
		log.Fatal(err)
	}
	if _, err = cli.UserAdd(context.TODO(), "u", "123"); err != nil {
		log.Fatal(err)
	}
	if _, err = cli.UserGrantRole(context.TODO(), "u", "r"); err != nil {
		log.Fatal(err)
	}*/
	if _, err = cli.AuthEnable(context.TODO()); err != nil {
		log.Fatal(err)
	}

	cliAuth, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
		Username:    "u",
		Password:    "123",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cliAuth.Close()

	if _, err = cliAuth.Put(context.TODO(), "foo1", "bar"); err != nil {
		log.Fatal(err)
	}

	_, err = cliAuth.Txn(context.TODO()).
			If(clientv3.Compare(clientv3.Value("zoo1"), ">", "abc")).
			Then(clientv3.OpPut("zoo1", "XYZ")).
			Else(clientv3.OpPut("zoo1", "ABC")).
			Commit()
	fmt.Println(err)

	// now check the permission with the root account
	rootCli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
		Username:    "root",
		Password:    "123",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer rootCli.Close()

	resp, err := rootCli.RoleGet(context.TODO(), "r")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("user u permission: key %q, range end %q\n", resp.Perm[0].Key, resp.Perm[0].RangeEnd)

	/*if _, err = rootCli.AuthDisable(context.TODO()); err != nil {
		log.Fatal(err)
	}*/
	// Output: etcdserver: permission denied
	// user u permission: key "foo", range end "zoo"
}