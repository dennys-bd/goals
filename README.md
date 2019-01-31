# goals
`goals` is a Golang/GraphQL Framework. `goals` is still in beta release.

## installation
`go get github.com/dennys-bd/goals`

## usage
`goals` has 3 main commands today

### goals init
`goals init PROJECTNAME` should start your project under your $GOPATH/src

### goals scaffold
#### scaffold graphql
- Usage:
You can generate a entire model in a syntax rails like, is recommended to use simple quote on the fields of model coz `]!` characters is reserved for bash terminals. DONT USE WHITE SPACES, it's used to separate the fields ;)
`goals s g <modelname> 'fieldName(param:ParamType,param2:ParamType2):fieldType! relationName:type:[ModelName]!'`

Eventually you can create your model inserting it's fields in line by line using the verbose flag (-v)
`goals s g <modelname> -v`

Console will ask you to insert the first field you can write something like this: (there is no problem using whitespaces where)
`fieldName(param: ParamType, param2: ParamType2): FieldType!`

then console will print informations about the field and you can confirm or erase this field
if confirmed it will asks to create another field until you send a blank line.

should create a suitable model, schema and resolver structures for your new type  
For now you should use simple quote in your attributes, -Use Spaces to separate the attributes-.  
Flags:
  * `--json`
   It creates a model genarating json value for each attribute of model.
  * `--no-model`
   Don't create the model. (Use it if you alread have a model)
  * `--no-schema`
   Don't create the schema. (Use it if you alread have a schema)
  * `--no-resolver`
   Don't create the resolver. (Use it if you alread have a resolver)
  * `--resolver` or `-r`
   Change the resolver name for your model.
  * `--verbose` or `-v`
   Is a way of creating your params for the attribute, and create each attribute separately. (You can use spaces now)
   NO QUITE READY YET.

#### scaffold auth
`goals s auth`
should create a private resolver, with it's auth structures.

covering commun graphql types: String, Int, Boolean, ID, Float, and Time from `github.com/graph-gophers/graphql-go`.  
Every other type will be treated as Scalar if you don't specifically declare type before modelName

### goals runserver
`goals r`
Start your server with some goals pattern configurations.  
Flags:
  * `--port` or `-p`
   Change the port to serv your goals application
  * `--env-port`
   Infer the port from the environment variable `PORT` (The `--port` flag is stronger if used together, please don't.)
  * `--env` or `-e`
   Starts your goals application with the configurations for specified environment - `goals r -e production`
  * `--verbose` or `-v`
   Verbose, right? for now it only prints the schemas


## todo

* [x] command init
* [x] command scaffold
* [x] command scaffold graphql
* [x] command scaffold auth
* [x] goals core to facilitate serving pages
* [x] create runserver command
* [x] environment separation ready (dotEnv)
* [x] accept params to resolver on scaffolding model
* [ ] scalar Image
* [ ] database integration:
  * [ ] save model directive (databases or gateway)
  * [ ] automagic model basics attributes (id, created_at, updated_at)
  * [ ] flag to dont create model basics attributes on goals s graphql
  * [ ] command makemigrations
  * [ ] command migrate
* [ ] flag to separate application in database directive ou gateway in init
* [ ] flag to separate model in databased ou delivered from external api in scaffold
* [ ] Write tests for check writing files
* [ ] versoning
  * [ ] check `go get`
  * [ ] use goals of vendor package primarily 
* [ ] goals core
* [ ] goals admin
