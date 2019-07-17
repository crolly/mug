# Adding to The Project

The `add` command lets you define a new **resource** or **function group**. 

## Adding a Simple Resource

The CRUDL functions will be generated in the `functions/resourcename/` folder and the `serverless.yml` will be updated accordingly.

For example, let's create a course resource, which has a couple of simple attributes. `cd` into the project directory and run:
```
mug add resource course -a "name,subtitle,description"
```

This will generate the following files/ subdirectories in the folder `functions/course/`:
```
-rw-r--r-- course.json
-rw-r--r-- course.go
-rw-r--r-- serverless.yml
-rw-r--r-- create/main.go
-rw-r--r-- delete/main.go
-rw-r--r-- list/main.go
-rw-r--r-- read/main.go
-rw-r--r-- update/main.go
```

The resource definition is kept track of in the `course.json` like this. I use this as a development step to eventually enable definition through `json` input:

```json
{
  "name": "course",
  "type": "Course",
  "ident": "course",
  "attributes": [
    {
      "name": "id",
      "ident": "id",
      "goType": "string",
      "awsType": "S",
      "hash": true
    },
    {
      "name": "name",
      "ident": "name",
      "goType": "string",
      "awsType": "S",
      "hash": false
    },
    {
      "name": "subtitle",
      "ident": "subtitle",
      "goType": "string",
      "awsType": "S",
      "hash": false
    },
    {
      "name": "description",
      "ident": "description",
      "goType": "string",
      "awsType": "S",
      "hash": false
    }
  ],
  "nested": null,
  "imports": [
    "github.com/gofrs/uuid"
  ]
}
```
The course's struct in the `course.go` looks like this:
```go
// Course defines the Course model
type Course struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Subtitle    string `json:"subtitle"`
	Description string `json:"description"`
}
```

### Identifying Attribute

As you can see an `ID` attribute of type `string` is automatically added. This `ID` will also be used as the hash key for the table.

It is important to note, that DynamoDB is not best friends with **UUID** types. The marshall/ unmarshall process is still a little buggy, so at this point it is safer to just store **UUIDs** as strings and marshall/ unmarshall them in the functions code, instead of letting DynamoDB do that work.

If wish to define your own identifying attribute, you can do so by specifying the required flags of the command. For details see the [Commands Reference](/commands/).

### General Attributes

The type for an attribute is by default `string` if none is provided. E.g., if you would like to add a number just provide the golang type with the attribute's name separated by a colon (:) like this:
```
mug add resource course -a "name,subtitle,description,price:float32"
```
The `course.go` file will also contain all the wrapper methods to interact with the database, just in case you want to have another solution as DynamoDB as persistence layer. The `main.go` files in the function subdirectories will contain the actual lambda functionality.

## Complex Resource Definition with Nested Objects

With Dynamo DB being a NoSQL database you certainly cannot use relationships like you may be used to from relational databases like MySQL or PostgreSQL. Usually you overcome this by deciding which entities you work with (querying, writing, etc.) and embedding all related information. 

With mug you can easily generate such resources with a little more complex attribute definition to the `-a` attributes flag:

* `{}` wraps nested objects (`1-1` relationship)
* `[]` wraps slices of objects (`1-n` relationship).

For example, let's generate a user, who has an address and may have multiple enrollments to courses: **Please make sure to wrap the attribute definition with `"`(double quotes) to ensure the recursive parsing works!** 

```
mug add resource user -a "name,isActive:bool,email,
address:{street,zip,city},
enrollments:[courseID,startDate:time.Time,endDate:time.Time]"
```

The generated go struct(s) will look like this:
```go
// User defines the User model
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
	Email    string `json:"email"`

	Address     Address       `json:"address"`
	Enrollments []Enrollment  `json:"enrollments"`
}

// Address defines the Address model
type Address struct {
	Street string `json:"street"`
	Zip    string `json:"zip"`
	City   string `json:"city"`
}

// Enrollment defines the Enrollments model
type Enrollment struct {
	CourseID  string    `json:"course_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}
```

Of course, you can also reference the `ID` in an object, however, you will have to manage this in your own code - getting the `ID` or multiple `IDs` for `m-n` relationships and fetching the referenced objects afterwards.

## Adding Function Groups, Functions etc.

Adding function groups or functions is very similar. You can also remove resources, function groups or functions from the previous.

**Have a look at the appropriate command syntax in the [Commands Reference](/commands/).**

