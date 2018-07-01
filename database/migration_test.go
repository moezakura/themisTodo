package main

import (
	"fmt"
	"os/exec"
	"testing"
)

func TestGetGenerationFiles(t *testing.T) {
	exec.Command("mkdir", "/tmp/gotest/").Run()
	exec.Command("touch", "/tmp/gotest/001-base.up.sql").Run()
	exec.Command("touch", "/tmp/gotest/001-base.down.sql").Run()
	exec.Command("touch", "/tmp/gotest/002-super_table.up.sql").Run()
	exec.Command("touch", "/tmp/gotest/002-super_table.down.sql").Run()
	exec.Command("touch", "/tmp/gotest/003-awesome_table.up.sql").Run()
	exec.Command("touch", "/tmp/gotest/003-awesome_table.down.sql").Run()
	defer func() {
		exec.Command("rm", "-fr", "/tmp/gotest").Run()
	}()
	upFiles := getGenerationFiles("/tmp/gotest", upFileRe)
	downFiles := getGenerationFiles("/tmp/gotest", downFileRe)
	fmt.Println("upgrade files: ", upFiles)
	fmt.Println("downgrade files: ", downFiles)
	if upFiles[1] != "/tmp/gotest/001-base.up.sql" ||
		upFiles[2] != "/tmp/gotest/002-super_table.up.sql" ||
		upFiles[3] != "/tmp/gotest/003-awesome_table.up.sql" {
		t.Fatal("Could not find upgrade migration file.")
	}
	if downFiles[1] != "/tmp/gotest/001-base.down.sql" ||
		downFiles[2] != "/tmp/gotest/002-super_table.down.sql" ||
		downFiles[3] != "/tmp/gotest/003-awesome_table.down.sql" {
		t.Fatal("Could not find downgrade migration file.")
	}
}

func TestGetGenerationName(t *testing.T) {
	filePath := "/tmp/gotest/001-base.up.sql"
	if name, _ := getGenerationName(filePath); name != "base" {
		fmt.Println("generation name: ", name)
		t.Fatal("Could not get generation name")
	}
}

func TestGetGenerations(t *testing.T) {
	gen001 := Generation{
		1,
		"/tmp/gotest/001-base.up.sql",
		"/tmp/gotest/001-base.down.sql",
		"base",
	}
	gen002 := Generation{
		2,
		"/tmp/gotest/002-super_table.up.sql",
		"/tmp/gotest/002-super_table.down.sql",
		"super_table",
	}
	gen003 := Generation{
		3,
		"/tmp/gotest/003-awesome_table.up.sql",
		"/tmp/gotest/003-awesome_table.down.sql",
		"awesome_table",
	}
	exec.Command("mkdir", "/tmp/gotest/").Run()
	exec.Command("touch", "/tmp/gotest/001-base.up.sql").Run()
	exec.Command("touch", "/tmp/gotest/001-base.down.sql").Run()
	exec.Command("touch", "/tmp/gotest/002-super_table.up.sql").Run()
	exec.Command("touch", "/tmp/gotest/002-super_table.down.sql").Run()
	exec.Command("touch", "/tmp/gotest/003-awesome_table.up.sql").Run()
	exec.Command("touch", "/tmp/gotest/003-awesome_table.down.sql").Run()
	defer func() {
		exec.Command("rm", "-fr", "/tmp/gotest").Run()
	}()
	g, err := getGenerations("/tmp/gotest")
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(g)
	if g[1] != gen001 || g[2] != gen002 || g[3] != gen003 {
		t.Fatal("Incorrect Generation")
	}
}
