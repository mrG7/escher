package leveldb

import (
	"testing"

	. "github.com/gocircuit/escher/kit/github.com/onsi/ginkgo"
	. "github.com/gocircuit/escher/kit/github.com/onsi/gomega"

	"github.com/gocircuit/escher/kit/github.com/syndtr/goleveldb/leveldb/testutil"
)

func TestLeveldb(t *testing.T) {
	testutil.RunDefer()

	RegisterFailHandler(Fail)
	RunSpecs(t, "Leveldb Suite")

	RegisterTestingT(t)
	testutil.RunDefer("teardown")
}
