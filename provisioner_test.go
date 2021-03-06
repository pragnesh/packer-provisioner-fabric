package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/packer"
)

// Be sure to remove the Fabric stub file in each test with:
//   defer os.Remove(config["command"].(string))
func testConfig(t *testing.T) map[string]interface{} {
	m := make(map[string]interface{})
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	fabric_stub := path.Join(wd, "packer-fabric-stub.sh")

	err = ioutil.WriteFile(fabric_stub, []byte("#!/usr/bin/env bash\necho -e \"Fabric 1.10.2\nParamiko 1.16.0\""), 0777)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	m["command"] = fabric_stub

	return m
}

func TestProvisioner_Impl(t *testing.T) {
	var raw interface{}
	raw = &Provisioner{}
	if _, ok := raw.(packer.Provisioner); !ok {
		t.Fatalf("must be a Provisioner")
	}
}

func TestProvisionerPrepare_Defaults(t *testing.T) {
	var p Provisioner
	config := testConfig(t)
	defer os.Remove(config["command"].(string))

	err := p.Prepare(config)
	if err == nil {
		t.Fatalf("should have error")
	}

	hostkey_file, err := ioutil.TempFile("", "hostkey")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(hostkey_file.Name())

	publickey_file, err := ioutil.TempFile("", "publickey")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(publickey_file.Name())

	fab_file, err := ioutil.TempFile("", "fabfile")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(fab_file.Name())

	config["ssh_host_key_file"] = hostkey_file.Name()
	config["ssh_authorized_key_file"] = publickey_file.Name()
	config["fab_file"] = fab_file.Name()
	config["fab_tasks"] = "test"
	err = p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvisionerPrepare_FabFile(t *testing.T) {
	var p Provisioner
	config := testConfig(t)
	defer os.Remove(config["command"].(string))

	hostkey_file, err := ioutil.TempFile("", "hostkey")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(hostkey_file.Name())

	publickey_file, err := ioutil.TempFile("", "publickey")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(publickey_file.Name())

	config["ssh_host_key_file"] = hostkey_file.Name()
	config["ssh_authorized_key_file"] = publickey_file.Name()

	err = p.Prepare(config)
	if err == nil {
		t.Fatal("should have error")
	}

	fab_file, err := ioutil.TempFile("", "fabfile")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(fab_file.Name())

	config["fab_file"] = fab_file.Name()
	config["fab_tasks"] = "test"
	err = p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvisionerPrepare_HostKeyFile(t *testing.T) {
	var p Provisioner
	config := testConfig(t)
	defer os.Remove(config["command"].(string))

	publickey_file, err := ioutil.TempFile("", "publickey")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(publickey_file.Name())

	fab_file, err := ioutil.TempFile("", "fabfile")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(fab_file.Name())

	filename := make([]byte, 10)
	n, err := io.ReadFull(rand.Reader, filename)
	if n != len(filename) || err != nil {
		t.Fatal("could not create random file name")
	}

	config["ssh_host_key_file"] = fmt.Sprintf("%x", filename)
	config["ssh_authorized_key_file"] = publickey_file.Name()
	config["fab_file"] = fab_file.Name()
	config["fab_tasks"] = "test"

	err = p.Prepare(config)
	if err == nil {
		t.Fatal("should error if ssh_host_key_file does not exist")
	}

	hostkey_file, err := ioutil.TempFile("", "hostkey")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(hostkey_file.Name())

	config["ssh_host_key_file"] = hostkey_file.Name()
	err = p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvisionerPrepare_AuthorizedKeyFile(t *testing.T) {
	var p Provisioner
	config := testConfig(t)
	defer os.Remove(config["command"].(string))

	hostkey_file, err := ioutil.TempFile("", "hostkey")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(hostkey_file.Name())

	fab_file, err := ioutil.TempFile("", "fabfile")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(fab_file.Name())

	filename := make([]byte, 10)
	n, err := io.ReadFull(rand.Reader, filename)
	if n != len(filename) || err != nil {
		t.Fatal("could not create random file name")
	}

	config["ssh_host_key_file"] = hostkey_file.Name()
	config["fab_file"] = fab_file.Name()
	config["fab_tasks"] = "test"
	config["ssh_authorized_key_file"] = fmt.Sprintf("%x", filename)

	err = p.Prepare(config)
	if err == nil {
		t.Errorf("should error if ssh_authorized_key_file does not exist")
	}

	publickey_file, err := ioutil.TempFile("", "publickey")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(publickey_file.Name())

	config["ssh_authorized_key_file"] = publickey_file.Name()
	err = p.Prepare(config)
	if err != nil {
		t.Errorf("err: %s", err)
	}
}

func TestProvisionerPrepare_LocalPort(t *testing.T) {
	var p Provisioner
	config := testConfig(t)
	defer os.Remove(config["command"].(string))

	hostkey_file, err := ioutil.TempFile("", "hostkey")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(hostkey_file.Name())

	publickey_file, err := ioutil.TempFile("", "publickey")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(publickey_file.Name())

	fab_file, err := ioutil.TempFile("", "fabfile")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(fab_file.Name())

	config["ssh_host_key_file"] = hostkey_file.Name()
	config["ssh_authorized_key_file"] = publickey_file.Name()
	config["fab_file"] = fab_file.Name()
	config["fab_tasks"] = "test"

	config["local_port"] = "65537"
	err = p.Prepare(config)
	if err == nil {
		t.Fatal("should have error")
	}

	config["local_port"] = "22222"
	err = p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}
