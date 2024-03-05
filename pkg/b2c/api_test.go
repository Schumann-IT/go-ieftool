package b2c

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"com.schumann-it.go-ieftool/pkg/b2c/environment"
	"github.com/stretchr/testify/assert"
)

var (
	testBaseDir           = "../../test"
	testBuildTargetDir    = "/tmp/b2ctests/build"
	testFixturesSourceDir = "fixtures/source"
	testFixturesConfigDir = "fixtures/config"
	testSourceDir         = ""
)

func setup(t *testing.T, env string) *Api {
	_ = os.RemoveAll(testBuildTargetDir)

	r, err := filepath.Abs(testBaseDir)
	cp := fmt.Sprintf("%s.yaml", path.Join(testBaseDir, testFixturesConfigDir, env))
	testSourceDir = path.Join(r, testFixturesSourceDir, env)
	a, err := NewApi(cp, testSourceDir, testBuildTargetDir)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	return a
}

func Test_NewApi(t *testing.T) {
	expected := environment.Config{
		Name:   "simple",
		Tenant: "test.onmicrosoft.com",
		Settings: map[string]string{
			"Tenant": "test.onmicrosoft.com",
		},
	}

	a := setup(t, "simple")
	actual := a.FindConfig(expected.Name)

	assert.Equal(t, expected, *actual)
}

func Test_BuildPolicies(t *testing.T) {
	a := setup(t, "simple")
	err := a.BuildPolicies("simple")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	expected := countFiles(testSourceDir)
	actual := countFiles(path.Join(testBuildTargetDir, "simple"))

	assert.Equal(t, expected, actual)
}

func Test_CreateDeployBatch(t *testing.T) {
	a := setup(t, "simple")
	_ = a.BuildPolicies("simple")

	_, err := a.Batch("simple")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	assert.Nil(t, err)
}

func countFiles(p string) int {
	c := 0
	_ = filepath.Walk(p, func(_ string, i os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if i.IsDir() {
			return nil
		}
		c++
		return nil
	})

	return c
}
