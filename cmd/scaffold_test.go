package cmd

import "testing"

type ScaffoldTest struct {
	attribute string
	typeName  string
	isModel   bool
	output    string
}

var schemaTest = []ScaffoldTest{
	{"attrCool", "string", false, "attrCool: String\n"},
	{"attrCool", "bool!", false, "attrCool: Boolean!\n"},
	{"attrCool", "Boolean!", false, "attrCool: Boolean!\n"},
	{"attrCool", "[id!]", false, "attrCool: [ID!]\n"},
	{"attrCool", "[json]!", false, "attrCool: [Json]!\n"},
	{"attrCool", "[time]", false, "attrCool: [Time]\n"},
	{"attrCool", "[int!]!", false, "attrCool: [Int!]!\n"},
	{"attrCool", "[float!]!", false, "attrCool: [Float!]!\n"},
	{"userCool", "myUser!", true, "userCool: MyUser!\n"},
	{"userCool", "MyUser", true, "userCool: MyUser\n"},
	{"userCool", "[myUser!]", true, "userCool: [MyUser!]\n"},
	{"userCool", "[MyUser!]", true, "userCool: [MyUser!]\n"},
	{"userCool", "[myUser]!", true, "userCool: [MyUser]!\n"},
	{"userCool", "[MyUser!]!", true, "userCool: [MyUser!]!\n"},
}

func TestGetSchemaLine(t *testing.T) {
	for _, test := range schemaTest {
		if result := getSchemaLine(test.attribute, test.typeName); result != test.output {
			t.Errorf(`getSchemaLine("%s", "%s"), wanted "%s", got "%s"`, test.attribute, test.typeName, test.output, result)
		}
	}
}

var modelTest = []ScaffoldTest{
	{"attrCool", "string", false, "AttrCool *string `json:\"attr_cool,omitempty\"`\n"},
	{"attrCool", "bool!", false, "AttrCool bool `json:\"attr_cool\"`\n"},
	{"attrCool", "Boolean!", false, "AttrCool bool `json:\"attr_cool\"`\n"},
	{"attrCool", "[id!]", false, "AttrCool *[]graphql.ID `json:\"attr_cool,omitempty\"`\n"},
	{"attrCool", "[json]!", false, "AttrCool []*scalar.Json `json:\"attr_cool,omitempty\"`\n"},
	{"attrCool", "[time]", false, "AttrCool *[]*time.Time `json:\"attr_cool,omitempty\"`\n"},
	{"attrCool", "[int!]!", false, "AttrCool []int32 `json:\"attr_cool,omitempty\"`\n"},
	{"attrCool", "[float!]!", false, "AttrCool []float64 `json:\"attr_cool,omitempty\"`\n"},
	{"userCool", "myUser!", true, "UserCool MyUser `json:\"-\"`\n"},
	{"userCool", "MyUser", true, "UserCool *MyUser `json:\"-\"`\n"},
	{"userCool", "[myUser!]", true, "UserCool *[]MyUser `json:\"-\"`\n"},
	{"userCool", "[myUser]!", true, "UserCool []*MyUser `json:\"-\"`\n"},
	{"userCool", "[myUser!]!", true, "UserCool []MyUser `json:\"-\"`\n"},
}

func TestGetModelLine(t *testing.T) {
	for _, test := range modelTest {
		if result := getModelLine(test.attribute, test.typeName, test.isModel); result != test.output {
			t.Errorf(`getSchemaLine("%s", "%s"), wanted "%s", got "%s"`, test.attribute, test.typeName, test.output, result)
		}
	}
}

var resolverTest = []ScaffoldTest{
	{"attrCool", "string", false, `func (r *{{.resolver}}) AttrCool() *string {
	return r.{{.abbreviation}}.AttrCool
	}
`},
	// {"attrCool", "bool!", false, "AttrCool bool `json:\"attr_cool\"`\n"},
	// {"attrCool", "Boolean!", false, "AttrCool bool `json:\"attr_cool\"`\n"},
	// {"attrCool", "[id!]", false, "AttrCool *[]graphql.ID `json:\"attr_cool,omitempty\"`\n"},
	// {"attrCool", "[json]!", false, "AttrCool []*scalar.Json `json:\"attr_cool,omitempty\"`\n"},
	// {"attrCool", "[time]", false, "AttrCool *[]*time.Time `json:\"attr_cool,omitempty\"`\n"},
	// {"attrCool", "[int!]!", false, "AttrCool []int32 `json:\"attr_cool,omitempty\"`\n"},
	// {"attrCool", "[float!]!", false, "AttrCool []float64 `json:\"attr_cool,omitempty\"`\n"},
	// {"userCool", "myUser!", true, "UserCool MyUser `json:\"-\"`\n"},
	// {"userCool", "MyUser", true, "UserCool *MyUser `json:\"-\"`\n"},
	// {"userCool", "[myUser!]", true, "UserCool *[]MyUser `json:\"-\"`\n"},
	// {"userCool", "[myUser]!", true, "UserCool []*MyUser `json:\"-\"`\n"},
	// {"userCool", "[myUser!]!", true, "UserCool []MyUser `json:\"-\"`\n"},
}

func TestGetResolverLine(t *testing.T) {
	for _, test := range resolverTest {
		if result := getResolverLine(test.attribute, test.typeName, test.isModel); result != test.output {
			t.Errorf(`getResolver("%s", "%s", "%t"), wanted "%s", got "%s"`, test.attribute, test.typeName, test.isModel, test.output, result)
		}
	}
}
