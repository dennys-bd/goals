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
* [x] remove page from server.go
* [x] create getDateInFormat method in model/helper
* [ ] Remove gqltype set it under goals core
* [ ] Goals core create and serve schemas (open and closed)
* [ ] make scaffold a parent command with s model a sub command
* [ ] create scaffold auth command
* [ ] auto generate json scalar on goals init
* [ ] create environments separation
* [ ] create runserver command
* [ ] time type should be string on graphql with format opptions
* [ ] auto generate getter in
* [ ] fix goals scaffold model syntax
* [ ] accept params to resolver on scaffolding model
* [ ] dotEnv
* [ ] cors
* [ ] database integration
* [ ] lib/schema should contains: Model Name, Resolver name, Type of retrival
* [ ] automagic model basics attributes (id, created_at, updated_at)
* [ ] flag to dont create model basics attributes on goals s model
* [ ] command makemigrations
* [ ] command migrate
* [ ] flag to separate application in database directive ou gateway in init
* [ ] flag to separate model in databased ou delivered from external api in scaffold
* [ ] Write tests for check writing files
* [ ] goals core
* [ ] goals admin
