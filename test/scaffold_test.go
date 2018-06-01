package test

import (
	"goals/cmd"
	"testing"
)

type ScaffoldTest struct {
	attribute string
	typeName  string
	isModel   bool
	output    string
}

var schemaTest = []ScaffoldTest{
	{"attrCool", "string", false, "	attrCool: String\n"},
	{"attrCool", "bool!", false, "	attrCool: Boolean!\n"},
	{"attrCool", "Boolean!", false, "	attrCool: Boolean!\n"},
	{"attrCool", "[id!]", false, "	attrCool: [ID!]\n"},
	{"attrCool", "[json]!", false, "	attrCool: [Json]!\n"},
	{"attrCool", "[time]", false, "	attrCool: [Time]\n"},
	{"attrCool", "[int!]!", false, "	attrCool: [Int!]!\n"},
	{"attrCool", "[float!]!", false, "	attrCool: [Float!]!\n"},
	{"userCool", "myUser!", true, "	userCool: MyUser!\n"},
	{"userCool", "MyUser", true, "	userCool: MyUser\n"},
	{"userCool", "[myUser!]", true, "	userCool: [MyUser!]\n"},
	{"userCool", "[MyUser!]", true, "	userCool: [MyUser!]\n"},
	{"userCool", "[myUser]!", true, "	userCool: [MyUser]!\n"},
	{"userCool", "[MyUser!]!", true, "	userCool: [MyUser!]!\n"},
}

func TestGetSchemaLine(t *testing.T) {
	for _, test := range schemaTest {
		if result := cmd.GetSchemaLine(test.attribute, test.typeName); result != test.output {
			t.Errorf(`getSchemaLine("%s", "%s"), wanted "%s", got "%s"`, test.attribute, test.typeName, test.output, result)
		}
	}
}

var modelTest = []ScaffoldTest{
	{"attrCool", "string", false, "	AttrCool *string `json:\"attr_cool,omitempty\"`\n"},
	{"attrCool", "bool!", false, "	AttrCool bool `json:\"attr_cool\"`\n"},
	{"attrCool", "Boolean!", false, "	AttrCool bool `json:\"attr_cool\"`\n"},
	{"attrCool", "[id!]", false, "	AttrCool *[]graphql.ID `json:\"attr_cool,omitempty\"`\n"},
	{"attrCool", "[json]!", false, "	AttrCool []*scalar.Json `json:\"attr_cool,omitempty\"`\n"},
	{"attrCool", "[time]", false, "	AttrCool *[]*time.Time `json:\"attr_cool,omitempty\"`\n"},
	{"attrCool", "[int!]!", false, "	AttrCool []int32 `json:\"attr_cool,omitempty\"`\n"},
	{"attrCool", "[float!]!", false, "	AttrCool []float64 `json:\"attr_cool,omitempty\"`\n"},
	{"userCool", "myUser!", true, "	UserCool MyUser `json:\"-\"`\n"},
	{"userCool", "MyUser", true, "	UserCool *MyUser `json:\"-\"`\n"},
	{"userCool", "[myUser!]", true, "	UserCool *[]MyUser `json:\"-\"`\n"},
	{"userCool", "[myUser]!", true, "	UserCool []*MyUser `json:\"-\"`\n"},
	{"userCool", "[myUser!]!", true, "	UserCool []MyUser `json:\"-\"`\n"},
}

func TestGetModelLine(t *testing.T) {
	for _, test := range modelTest {
		if result := cmd.GetModelLine(test.attribute, test.typeName, test.isModel); result != test.output {
			t.Errorf(`getSchemaLine("%s", "%s"), wanted "%s", got "%s"`, test.attribute, test.typeName, test.output, result)
		}
	}
}

var resolverTest = []ScaffoldTest{
	{"attrCool", "string", false, `func (r *{{.resolver}}) AttrCool() *string {
	return r.{{.abbreviation}}.AttrCool
}
`},
	{"attrCool", "bool!", false, `func (r *{{.resolver}}) AttrCool() bool {
	return r.{{.abbreviation}}.AttrCool
}
`},
	{"attrCool", "Boolean!", false, `func (r *{{.resolver}}) AttrCool() bool {
	return r.{{.abbreviation}}.AttrCool
}
`},
	{"attrCool", "[id!]", false, `func (r *{{.resolver}}) AttrCool() *[]graphql.ID {
	return r.{{.abbreviation}}.AttrCool
}
`},
	{"attrCool", "[json]!", false, `func (r *{{.resolver}}) AttrCool() []*scalar.Json {
	return r.{{.abbreviation}}.AttrCool
}
`},
	{"attrCool", "[time]", false, `func (r *{{.resolver}}) AttrCool() *[]*time.Time {
	return r.{{.abbreviation}}.AttrCool
}
`},
	{"attrCool", "[int!]!", false, `func (r *{{.resolver}}) AttrCool() []int32 {
	return r.{{.abbreviation}}.AttrCool
}
`},
	{"attrCool", "[float!]!", false, `func (r *{{.resolver}}) AttrCool() []float64 {
	return r.{{.abbreviation}}.AttrCool
}
`},
	{"userCool", "myUser!", true, `func (r *{{.resolver}}) UserCool() *myUserResolver {
	return &myUserResolver{&r.{{.abbreviation}}.UserCool}
}
`},
	{"userCool", "MyUser", true, `func (r *{{.resolver}}) UserCool() *myUserResolver {
	return &myUserResolver{r.{{.abbreviation}}.UserCool}
}
`},
	{"userCool", "[myUser!]", true, `func (r *{{.resolver}}) UserCool() *[]*myUserResolver {
	if r.{{.abbreviation}}.UserCool == nil {
		return nil
	}
	slice := *r.{{.abbreviation}}.UserCool

	l := make([]*myUserResolver, len(slice))
	for i := range l {
		l[i] = &myUserResolver{&slice[i]}
	}

	return &l
}
`},
	// []*MyUser
	{"userCool", "[myUser]!", true, `func (r *{{.resolver}}) UserCool() []*myUserResolver {
	slice := r.{{.abbreviation}}.UserCool

	l := make([]*myUserResolver, len(slice))
	for i := range l {
		l[i] = &myUserResolver{slice[i]}
	}

	return l
}
`},
	// []MyUser
	{"userCool", "[myUser!]!", true, `func (r *{{.resolver}}) UserCool() []*myUserResolver {
	slice := r.{{.abbreviation}}.UserCool

	l := make([]*myUserResolver, len(slice))
	for i := range l {
		l[i] = &myUserResolver{&slice[i]}
	}

	return l
}
`},
}

func TestGetResolverLine(t *testing.T) {
	for _, test := range resolverTest {
		if result := cmd.GetResolverLine(test.attribute, test.typeName, test.isModel); result != test.output {
			t.Errorf(`getResolver("%s", "%s", "%t"), wanted "%s", got "%s"`, test.attribute, test.typeName, test.isModel, test.output, result)
		}
	}
}
