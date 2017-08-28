package dotaccess

import (
	. "gopkg.in/check.v1"

	"log"
	// "net/url"
	// "github.com/oleiade/reflections"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type TestSuite struct{}

var _ = Suite(&TestSuite{})

type NullWriter int

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func (s *TestSuite) SetUpTest(c *C) {

	if !testing.Verbose() {
		log.SetOutput(new(NullWriter))
	}

}

type TestStruct struct {
	Foo      int
	Child    ChildStruct
	ChildRef *ChildStruct
}

type ChildStruct struct {
	Baz string
}

func (s *TestSuite) TestGetStruct(c *C) {
	test := &TestStruct{
		Foo: 5,
		Child: ChildStruct{
			Baz: "baz",
		},
		ChildRef: &ChildStruct{
			Baz: "bing",
		},
	}

	getFoo, _ := Get(test, "foo")
	getChildBaz, _ := Get(test, "child.baz")
	getChildRefBaz, _ := Get(test, "childRef.baz")
	c.Assert(getFoo, Equals, 5)
	c.Assert(getChildBaz, Equals, "baz")
	c.Assert(getChildRefBaz, Equals, "bing")

	_, err := Get(test, "asdf")
	c.Assert(err, Not(Equals), nil)
}

func (s *TestSuite) TestGetMap(c *C) {
	test := map[string]interface{}{
		"foo": "5",
		"child": map[string]int{
			"a": 1,
			"b": 1,
		},
	}

	getFoo, _ := Get(test, "foo")
	getChild, _ := Get(test, "child.a")
	c.Assert(getFoo, Equals, "5")
	c.Assert(getChild, Equals, 1)
}

// func (s *TestSuite) TestSetStruct(c *C) {
// 	test := &TestStruct{
// 		Foo: 5,
// 		Child: ChildStruct{
// 			Baz: "baz",
// 		},
// 		ChildRef: &ChildStruct{
// 			Baz: "bing",
// 		},
// 	}

// 	// err := reflections.SetField(test.ChildRef, "Baz", "boop")
// 	// log.Println(err)

// 	Set(test, "foo", 6)
// 	c.Assert(test.Foo, Equals, 6)

// 	err := Set(test, "child.baz", "bar")
// 	c.Assert(err, Not(IsNil))

// 	Set(test, "childRef.baz", "boof")
// 	c.Assert(test.ChildRef.Baz, Equals, "boof")

// }

func (s *TestSuite) TestSetMap(c *C) {
	test := map[string]interface{}{
		"foo": "5",
		"child": map[string]int{
			"a": 1,
			"b": 1,
		},
	}

	Set(test, "foo", "7")
	c.Assert(test["foo"], Equals, "7")

	Set(test, "child.a", 2)

	// Assert
	mp := test["child"].(map[string]int)
	c.Assert(mp["a"], Equals, 2)
	// log.Println(i)
}
