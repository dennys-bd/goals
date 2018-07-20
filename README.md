# goals
`goals` is a Golang/GraphQL Boilerplate Generator maybe a Framework.. `goals` is still in beta release.

## installation
`go get github.com/dennys-bd/goals`

## usage
`goals` has 2 commands today

`goals init [PROJECTNAME]` should start your project under your $GOPATH/src

`goals s [MODELNAME] attributeName:AttributeType! relationName:type:[ModelName]!`
should create a suitable model, schema and resolver structures for your new type

covering commun graphql types: String, Int, Boolean, ID, Float, and Time from graph-gophers/graphql-go every other type will be treated as Scalar if you don't specifically declare type before modelName


## todo

* [x] command init
* [x] command scaffold
* [x] write goalsfile
* [x] retrieving project from goalsfile
* [x] name fix on create project
* [x] more directory options in goals init
* [x] fix import model on resolver
* [x] uppercase if the attribute name is id (model, resolver)
* [x] import graphql and scalar in model and resolver if needed
* [x] create templates package
* [x] unit test for project.go
* [x] replace goals command text
* [x] flag to resolver name on scaffold model
* [ ] automagic model basics attributes (id, created_at, updated_at)
* [ ] flag to dont create model basics attributes
* [ ] authentication
* [ ] dotEnv
* [ ] cors
* [ ] goals core
* [ ] command runserver
* [ ] relations should always be pointer?
* [ ] lib/schema should contains: Model Name, Resolver name, Type of retrival
* [ ] flag to separate model in databased ou delivered from external api in scaffold
* [ ] command migrate
* [ ] command s migration
* [ ] create integration with gorm
* [ ] flag to separate application in database directive ou gateway in init
* [ ] Write tests for check writing files
* [ ] Brew installation
* [ ] hide package gqltype under goals core??
