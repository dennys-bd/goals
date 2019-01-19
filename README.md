# goals
`goals` is a Golang/GraphQL Framework. `goals` is still in beta release.

## installation
`go get github.com/dennys-bd/goals`

## usage
`goals` has 3 main commands today

### goals init
`goals init PROJECTNAME` should start your project under your $GOPATH/src

### goals scaffold
`goals s g MODELNAME attributeName:AttributeType! relationName:type:[ModelName]!`
should create a suitable model, schema and resolver structures for your new type

`goals s auth`
should create a private resolver, with it's auth structures.

covering commun graphql types: String, Int, Boolean, ID, Float, and Time from `github.com/graph-gophers/graphql-go`.  
Every other type will be treated as Scalar if you don't specifically declare type before modelName

### goals runserver
`goals r`
Start your server with some goals pattern configurations.  
You can change the port `goals r -p 3000`  
You can infer the port from a environment variable `goals r --env-port`  
You can set verbose on start to see check your schema `goals r -v`  
You can set your `GOALS_ENV` environment variable with the flag `--env` or `-e`

## todo

* [x] command init
* [x] command scaffold
* [x] command scaffold graphql
* [x] command scaffold auth
* [x] goals core to facilitate serving pages
* [x] create runserver command
* [x] environment separation ready (dotEnv)
* [ ] accept params to resolver on scaffolding model
* [ ] Ennable cors options
* [ ] Choose to serve graphilQL on CoreOpts
* [ ] Create flags --no-model, --no-resolver, --no-schema on goals s graphql
* [ ] database integration
* [ ] lib/schema should contains: Model Name, Resolver name, Type of retrival
* [ ] automagic model basics attributes (id, created_at, updated_at)
* [ ] flag to dont create model basics attributes on goals s graphql
* [ ] command makemigrations
* [ ] command migrate
* [ ] flag to separate application in database directive ou gateway in init
* [ ] flag to separate model in databased ou delivered from external api in scaffold
* [ ] Write tests for check writing files
* [ ] goals core
* [ ] goals admin
