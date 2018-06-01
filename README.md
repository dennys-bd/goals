# goals
`goals` is a Golang/GraphQL Boilerplate Generator maybe a Framework.. `goals` is still in beta release.

# use
`goals` has 2 commands today

`goals init [PROJECTNAME]` should start your project unde your $GOPATH

`goals s [MODELNAME] attributeName:AttributeType! relationName:type:[ModelName]!`
should create a suitable model, schema and resolver structures for your new type

covering commun graphql types: String, Int, Boolean, ID, Float, and Time from graph-gophers/graphql-go every other type will be treated as Scalar if you don't specifically declare type before modelName


# todo

* [x] command init
* [x] command scaffold
* [ ] Write goalsfile
* [ ] retrieving project from goalsfile
* [ ] more directory options in goals init
* [ ] test goals init
* [ ] remove templates from scaffold
* [ ] create templates package?
* [ ] Create integration with gorm
* [ ] Flag to resolver name on scaffold model
* [ ] Flag to separate application in database directive ou gateway in init
* [ ] Flag to separate model in databased ou delivered from external api in scaffold
* [ ] authentication
* [ ] lib/schema should contains: Model Name, Resolver name, Type of retrival
* [ ] dotEnv
* [ ] cors
* [ ] command runserver
* [ ] command migrate
* [ ] command s migration
* [ ] goals core
* [ ] hide package gqltype under goals core??
* [ ] Brew installation