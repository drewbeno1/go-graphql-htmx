# Sample Go, graphQL, HTMX app

step 1 initialize empty project: 
```
go mod init github.com/drewbeno1/go-graphql-htmx
```

step 2 get gql packages and initialize gql code generator:
```
go get github.com/99designs/gqlgen github.com/google/uuid 
gqlgen init
```

step 3 cleanup the auto generated resolver setup:
```
# delete schema.resolvers.go file
# remove "resolvers" block from gqlgen.yml cause we want to control our own resolvers
```

step 4 start on your app by writing your main.go and your resolvers.go and index.html

step 5:
```
# anytime you make a schema change, run:
gqlgen generate
```


