# GOP
## Git
# Object
# Proxy

Continuing on from revision1.md, this is the updated design doc, 
GOP, or Git Object Proxy is a generic tool that allows for git objects to be retrieved, processed and cached. It can be used as a generic package or can be started as a service.

## Requirements
- Generic
    - Must be easily extendable or completely file type agnostic
- Ability to easily import from and git service (github/gitlab/bitbucket)
    - Must be platform agnostic
- Ability to import from multiple versions of a specification without name clashes
- Ability to be language agnostic
    - Can't assume that the repository has any other code in it
    - Current solution requires a go.mod file to exist in target repo

## Nice to haves
- Generic 
- Fast
- Ability to construct a tree of every sysl file in existence
    - Having a central proxy to track the git tags of every repo that uses sysl can track dependencies and build import graphs of every sysl module in existence
- Not needing to worry about git credentials if running on an internal network
    - If The central proxy has access to all the git repos in an organisation, anyone with the repo can edit and build the source code without worrying about git credentials. This also opens up the possibility for something like sysl-catalog to be a service that has access to all the repos that the sysl proxy has access to.

## Architecture
There are 4 main parts:
    - Client
    - Proxy
    - VCS
    - Data store
These responsibilities have been split into 3 Interfaces:
    - Retriever
    - Processor
    - Cacher
```go

type Cacher interface {
	Cache(res *pbmod.Object) (err error)
}
```
#### Object type
Defined as 
```go 
type Object struct {
	Repo      string  `json:"repo"`
	Resource  string  `json:"resource"`
	Version   string  `json:"version"`
	Content   string  `json:"content"`
	Imported  bool    `json:"imported"`
	Processed *string `json:"processed,omitempty"`
}
```
Where:
 - `Repo` is the target repository (eg `github.com/joshcarp/gop`)
 - `Resource` is the target file (eg `pbmod.sysl`)
 - `Version` is the target commit hash (eg `165081dd92025fb5cae3fef575eca1ad9521e4cc`)
 - `Content` is the content of the file that will be returned
 - `Imported` is a flag to indicate whether any post processing jobs are required (for example is the Object should be cached by later steps)
 - `Processed` is an optional field that holds the result of any `Processor.Process` steps


## Retriever interface
```go
type Retriever interface {
	Retrieve(res *pbmod.Object) (err error)
}
```
- Used to retrieve a specified resource
 
Required Object fields:
    - `Repo`
    - `Resource`
    - `Version`
Should populate:
    - `Content`
    - `Imported` should be set tp true so that later caching steps will not execute
    - `Processed` if applicable 
    
Examples:
    - [../gop/retriever/retriever_fs/fs.go]() Retrieve from a filesystem
    - [../gop/retriever/retriever_gcs/gcs.go]() Retrieve from google cloud storage
    - [../gop/retriever/retriever_git/git.go]() Retrieve from a git repository

#### Processor interface
```go
type Processor interface {
	Process(pre *pbmod.Object) (err error)
}
```
- Used to process an object and populate the `Object.Processed` field
Required Object fields:
    - `Content`
Should populate:
    - `Processed`
Examples:
    - [../gop/processor/processor_sysl/sysl.go]() Process a sysl source file into parsed sysl protobuf json  

#### Cacher interface
```go
type Processor interface {
	Process(pre *pbmod.Object) (err error)
}
```
- Used to cache the given Object to a data source
Required Object fields:
    - `Repo`
    - `Resource`
    - `Version`
    - `Content` and/or `Processed` 
Should populate: Read only: Doesn't populate
Examples:
    - [../gop/cacher/cacher_fs/fs.go]() Cache a file to a file system
    - [../gop/cacher/cacher_gcs/gcs.go]() Cache a file to a gcs bucket

## Sequence Diagram

<img src="https://www.plantuml.com/plantuml/png/dP51JuGm48Nl_8evOWBbJdH3ieddZVu2B4yabivcMiZwxqqfAhYOu9wObkatyzwhdA_53xr9ZgQ3zPGVw2Hy-IZf2LuwZ13rLQNJdun6xPJclc1f2y6PYzVEGE7Ygn5oboIryQHh_OOc8QB82-1dpmBP9BTEgvT1lyDVupDQydzEwYoiuHoQE9U8vX4B5G8_YDr5M2yR_VW_8BvR4ex12b7J9sMd7br6QyTW7AZdPZ0WorlcvOTWIqdQiCMLGulE-poFdT_tCGrBywhJtKffBCe_HD43dUAPHSrLkbu_q60tmwPVybljvfptfXh0iwT1MunrNnotH571DaDlFW40">