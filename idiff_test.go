package idiff

import (
	"testing"

	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

type MenuItem struct {
	Value   string
	Onclick string
}

type Popup struct {
	MenuItem []MenuItem
}

type Data struct {
	ID        string
	Value     string
	List1     []string
	List2     []string
	List3     []string
	Popup     Popup
	Func      func()
	Recursive *Data
	Map       map[string]string
}

var teststruct1 = Data{
	ID:    "file1",
	Value: "value1",
	List1: []string{"hey", "there"},
	List2: []string{},
	List3: []string{"removed"},
	Popup: Popup{
		MenuItem: []MenuItem{
			{Value: "New", Onclick: "CreateNewDoc()"},
			{Value: "Open", Onclick: "OpenDoc()"},
			{Value: "Close", Onclick: "CloseDoc()"},
		},
	},
	Func:      func() {},
	Recursive: &Data{ID: "rec1"},
	Map:       map[string]string{"key1": "value2"},
}

var teststruct2 = Data{
	ID:    "file2",
	Value: "value2",
	List1: []string{"changed1", "changed2"},
	List2: []string{"added"},
	List3: []string{},
	Popup: Popup{
		MenuItem: []MenuItem{
			{Value: "Newier", Onclick: "CreateNewDoc()"},
			{Value: "Open it", Onclick: "OpenDoc()"},
			{Value: "Close it", Onclick: "CloseDoc()"},
		},
	},
	Recursive: &Data{ID: "rec2"},
	Map:       map[string]string{"key2": "value2"},
}

func TestStructDiffFormatTest1to2(t *testing.T) {
	a := assertions.New(t)

	diff, equal := Diff(teststruct1, teststruct2)
	a.So(equal, should.BeFalse)
	a.So(FormatTest(diff), should.Equal, `idiff.Data.List2[0]: added: "added"
idiff.Data.Map["key2"]: added: "value2"
idiff.Data.List3[0]: removed: "removed"
idiff.Data.Map["key1"]: removed: "value2"
idiff.Data.ID: got: "file2", expected: "file1"
idiff.Data.Value: got: "value2", expected: "value1"
idiff.Data.List1[0]: got: "changed1", expected: "hey"
idiff.Data.List1[1]: got: "changed2", expected: "there"
idiff.Data.Popup.MenuItem[0].Value: got: "Newier", expected: "New"
idiff.Data.Popup.MenuItem[1].Value: got: "Open it", expected: "Open"
idiff.Data.Popup.MenuItem[2].Value: got: "Close it", expected: "Close"
idiff.Data.Func: got: nil, expected: not nil
idiff.Data.Recursive.ID: got: "rec2", expected: "rec1"`)
}

func TestStructDiffFormatTest2to1(t *testing.T) {
	a := assertions.New(t)

	diff, equal := Diff(teststruct2, teststruct1)
	a.So(equal, should.BeFalse)
	a.So(FormatTest(diff), should.Equal, `idiff.Data.List3[0]: added: "removed"
idiff.Data.Map["key1"]: added: "value2"
idiff.Data.List2[0]: removed: "added"
idiff.Data.Map["key2"]: removed: "value2"
idiff.Data.ID: got: "file1", expected: "file2"
idiff.Data.Value: got: "value1", expected: "value2"
idiff.Data.List1[0]: got: "hey", expected: "changed1"
idiff.Data.List1[1]: got: "there", expected: "changed2"
idiff.Data.Popup.MenuItem[0].Value: got: "New", expected: "Newier"
idiff.Data.Popup.MenuItem[1].Value: got: "Open", expected: "Open it"
idiff.Data.Popup.MenuItem[2].Value: got: "Close", expected: "Close it"
idiff.Data.Func: got: not nil, expected: nil
idiff.Data.Recursive.ID: got: "rec1", expected: "rec2"`)
}

func TestSimple(t *testing.T) {
	a := assertions.New(t)

	diff, equal := Diff("hey", "there")
	a.So(equal, should.BeFalse)
	a.So(FormatTest(diff), should.Equal, `string: got: "there", expected: "hey"`)

	diff, equal = Diff("hey", "hey")
	a.So(equal, should.BeTrue)
	a.So(FormatTest(diff), should.Equal, ``)
}
