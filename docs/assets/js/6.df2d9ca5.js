(window.webpackJsonp=window.webpackJsonp||[]).push([[6],{193:function(s,t,a){"use strict";a.r(t);var e=a(0),n=Object(e.a)({},function(){var s=this,t=s.$createElement,a=s._self._c||t;return a("ContentSlotsDistributor",{attrs:{"slot-key":s.$parent.slotKey}},[a("h1",{attrs:{id:"adding-to-the-project"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#adding-to-the-project","aria-hidden":"true"}},[s._v("#")]),s._v(" Adding to The Project")]),s._v(" "),a("p",[s._v("The "),a("code",[s._v("add")]),s._v(" command lets you define a new "),a("strong",[s._v("resource")]),s._v(" or "),a("strong",[s._v("function group")]),s._v(".")]),s._v(" "),a("h2",{attrs:{id:"adding-a-simple-resource"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#adding-a-simple-resource","aria-hidden":"true"}},[s._v("#")]),s._v(" Adding a Simple Resource")]),s._v(" "),a("p",[s._v("The CRUDL functions will be generated in the "),a("code",[s._v("functions/resourcename/")]),s._v(" folder and the "),a("code",[s._v("serverless.yml")]),s._v(" will be updated accordingly.")]),s._v(" "),a("p",[s._v("For example, let's create a course resource, which has a couple of simple attributes. "),a("code",[s._v("cd")]),s._v(" into the project directory and run:")]),s._v(" "),a("div",{staticClass:"language- line-numbers-mode"},[a("pre",{pre:!0,attrs:{class:"language-text"}},[a("code",[s._v('mug add resource course -a "name,subtitle,description"\n')])]),s._v(" "),a("div",{staticClass:"line-numbers-wrapper"},[a("span",{staticClass:"line-number"},[s._v("1")]),a("br")])]),a("p",[s._v("This will generate the following files/ subdirectories in the folder "),a("code",[s._v("functions/course/")]),s._v(":")]),s._v(" "),a("div",{staticClass:"language- line-numbers-mode"},[a("pre",{pre:!0,attrs:{class:"language-text"}},[a("code",[s._v("-rw-r--r-- course.json\n-rw-r--r-- course.go\n-rw-r--r-- serverless.yml\n-rw-r--r-- create/main.go\n-rw-r--r-- delete/main.go\n-rw-r--r-- list/main.go\n-rw-r--r-- read/main.go\n-rw-r--r-- update/main.go\n")])]),s._v(" "),a("div",{staticClass:"line-numbers-wrapper"},[a("span",{staticClass:"line-number"},[s._v("1")]),a("br"),a("span",{staticClass:"line-number"},[s._v("2")]),a("br"),a("span",{staticClass:"line-number"},[s._v("3")]),a("br"),a("span",{staticClass:"line-number"},[s._v("4")]),a("br"),a("span",{staticClass:"line-number"},[s._v("5")]),a("br"),a("span",{staticClass:"line-number"},[s._v("6")]),a("br"),a("span",{staticClass:"line-number"},[s._v("7")]),a("br"),a("span",{staticClass:"line-number"},[s._v("8")]),a("br")])]),a("p",[s._v("The resource definition is kept track of in the "),a("code",[s._v("course.json")]),s._v(" like this. I use this as a development step to eventually enable definition through "),a("code",[s._v("json")]),s._v(" input:")]),s._v(" "),a("div",{staticClass:"language-json line-numbers-mode"},[a("pre",{pre:!0,attrs:{class:"language-json"}},[a("code",[a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("{")]),s._v("\n  "),a("span",{pre:!0,attrs:{class:"token property"}},[s._v('"name"')]),a("span",{pre:!0,attrs:{class:"token operator"}},[s._v(":")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('"course"')]),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v(",")]),s._v("\n  "),a("span",{pre:!0,attrs:{class:"token property"}},[s._v('"type"')]),a("span",{pre:!0,attrs:{class:"token operator"}},[s._v(":")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('"Course"')]),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v(",")]),s._v("\n  "),a("span",{pre:!0,attrs:{class:"token property"}},[s._v('"ident"')]),a("span",{pre:!0,attrs:{class:"token operator"}},[s._v(":")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('"course"')]),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v(",")]),s._v("\n  "),a("span",{pre:!0,attrs:{class:"token property"}},[s._v('"attributes"')]),a("span",{pre:!0,attrs:{class:"token operator"}},[s._v(":")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("[")]),s._v("\n    "),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("{")]),s._v("\n      "),a("span",{pre:!0,attrs:{class:"token property"}},[s._v('"name"')]),a("span",{pre:!0,attrs:{class:"token operator"}},[s._v(":")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('"id"')]),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v(",")]),s._v("\n      "),a("span",{pre:!0,attrs:{class:"token property"}},[s._v('"ident"')]),a("span",{pre:!0,attrs:{class:"token operator"}},[s._v(":")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('"id"')]),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v(",")]),s._v("\n      "),a("span",{pre:!0,attrs:{class:"token property"}},[s._v('"goType"')]),a("span",{pre:!0,attrs:{class:"token operator"}},[s._v(":")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('"string"')]),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v(",")]),s._v("\n      "),a("span",{pre:!0,attrs:{class:"token property"}},[s._v('"awsType"')]),a("span",{pre:!0,attrs:{class:"token operator"}},[s._v(":")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('"S"')]),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v(",")]),s._v("\n      "),a("span",{pre:!0,attrs:{class:"token property"}},[s._v('"hash"')]),a("span",{pre:!0,attrs:{class:"token operator"}},[s._v(":")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token boolean"}},[s._v("true")]),s._v("\n    "),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("}")]),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v(",")]),s._v("\n    "),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("{")]),s._v("\n      "),a("span",{pre:!0,attrs:{class:"token property"}},[s._v('"name"')]),a("span",{pre:!0,attrs:{class:"token operator"}},[s._v(":")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('"name"')]),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v(",")]),s._v("\n      "),a("span",{pre:!0,attrs:{class:"token property"}},[s._v('"ident"')]),a("span",{pre:!0,attrs:{class:"token operator"}},[s._v(":")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('"name"')]),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v(",")]),s._v("\n      "),a("span",{pre:!0,attrs:{class:"token property"}},[s._v('"goType"')]),a("span",{pre:!0,attrs:{class:"token operator"}},[s._v(":")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('"string"')]),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v(",")]),s._v("\n      "),a("span",{pre:!0,attrs:{class:"token property"}},[s._v('"awsType"')]),a("span",{pre:!0,attrs:{class:"token operator"}},[s._v(":")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('"S"')]),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v(",")]),s._v("\n      "),a("span",{pre:!0,attrs:{class:"token property"}},[s._v('"hash"')]),a("span",{pre:!0,attrs:{class:"token operator"}},[s._v(":")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token boolean"}},[s._v("false")]),s._v("\n    "),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("}")]),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v(",")]),s._v("\n    "),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("{")]),s._v("\n      "),a("span",{pre:!0,attrs:{class:"token property"}},[s._v('"name"')]),a("span",{pre:!0,attrs:{class:"token operator"}},[s._v(":")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('"subtitle"')]),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v(",")]),s._v("\n      "),a("span",{pre:!0,attrs:{class:"token property"}},[s._v('"ident"')]),a("span",{pre:!0,attrs:{class:"token operator"}},[s._v(":")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('"subtitle"')]),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v(",")]),s._v("\n      "),a("span",{pre:!0,attrs:{class:"token property"}},[s._v('"goType"')]),a("span",{pre:!0,attrs:{class:"token operator"}},[s._v(":")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('"string"')]),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v(",")]),s._v("\n      "),a("span",{pre:!0,attrs:{class:"token property"}},[s._v('"awsType"')]),a("span",{pre:!0,attrs:{class:"token operator"}},[s._v(":")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('"S"')]),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v(",")]),s._v("\n      "),a("span",{pre:!0,attrs:{class:"token property"}},[s._v('"hash"')]),a("span",{pre:!0,attrs:{class:"token operator"}},[s._v(":")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token boolean"}},[s._v("false")]),s._v("\n    "),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("}")]),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v(",")]),s._v("\n    "),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("{")]),s._v("\n      "),a("span",{pre:!0,attrs:{class:"token property"}},[s._v('"name"')]),a("span",{pre:!0,attrs:{class:"token operator"}},[s._v(":")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('"description"')]),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v(",")]),s._v("\n      "),a("span",{pre:!0,attrs:{class:"token property"}},[s._v('"ident"')]),a("span",{pre:!0,attrs:{class:"token operator"}},[s._v(":")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('"description"')]),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v(",")]),s._v("\n      "),a("span",{pre:!0,attrs:{class:"token property"}},[s._v('"goType"')]),a("span",{pre:!0,attrs:{class:"token operator"}},[s._v(":")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('"string"')]),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v(",")]),s._v("\n      "),a("span",{pre:!0,attrs:{class:"token property"}},[s._v('"awsType"')]),a("span",{pre:!0,attrs:{class:"token operator"}},[s._v(":")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('"S"')]),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v(",")]),s._v("\n      "),a("span",{pre:!0,attrs:{class:"token property"}},[s._v('"hash"')]),a("span",{pre:!0,attrs:{class:"token operator"}},[s._v(":")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token boolean"}},[s._v("false")]),s._v("\n    "),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("}")]),s._v("\n  "),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("]")]),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v(",")]),s._v("\n  "),a("span",{pre:!0,attrs:{class:"token property"}},[s._v('"nested"')]),a("span",{pre:!0,attrs:{class:"token operator"}},[s._v(":")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token null keyword"}},[s._v("null")]),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v(",")]),s._v("\n  "),a("span",{pre:!0,attrs:{class:"token property"}},[s._v('"imports"')]),a("span",{pre:!0,attrs:{class:"token operator"}},[s._v(":")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("[")]),s._v("\n    "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('"github.com/gofrs/uuid"')]),s._v("\n  "),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("]")]),s._v("\n"),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("}")]),s._v("\n")])]),s._v(" "),a("div",{staticClass:"line-numbers-wrapper"},[a("span",{staticClass:"line-number"},[s._v("1")]),a("br"),a("span",{staticClass:"line-number"},[s._v("2")]),a("br"),a("span",{staticClass:"line-number"},[s._v("3")]),a("br"),a("span",{staticClass:"line-number"},[s._v("4")]),a("br"),a("span",{staticClass:"line-number"},[s._v("5")]),a("br"),a("span",{staticClass:"line-number"},[s._v("6")]),a("br"),a("span",{staticClass:"line-number"},[s._v("7")]),a("br"),a("span",{staticClass:"line-number"},[s._v("8")]),a("br"),a("span",{staticClass:"line-number"},[s._v("9")]),a("br"),a("span",{staticClass:"line-number"},[s._v("10")]),a("br"),a("span",{staticClass:"line-number"},[s._v("11")]),a("br"),a("span",{staticClass:"line-number"},[s._v("12")]),a("br"),a("span",{staticClass:"line-number"},[s._v("13")]),a("br"),a("span",{staticClass:"line-number"},[s._v("14")]),a("br"),a("span",{staticClass:"line-number"},[s._v("15")]),a("br"),a("span",{staticClass:"line-number"},[s._v("16")]),a("br"),a("span",{staticClass:"line-number"},[s._v("17")]),a("br"),a("span",{staticClass:"line-number"},[s._v("18")]),a("br"),a("span",{staticClass:"line-number"},[s._v("19")]),a("br"),a("span",{staticClass:"line-number"},[s._v("20")]),a("br"),a("span",{staticClass:"line-number"},[s._v("21")]),a("br"),a("span",{staticClass:"line-number"},[s._v("22")]),a("br"),a("span",{staticClass:"line-number"},[s._v("23")]),a("br"),a("span",{staticClass:"line-number"},[s._v("24")]),a("br"),a("span",{staticClass:"line-number"},[s._v("25")]),a("br"),a("span",{staticClass:"line-number"},[s._v("26")]),a("br"),a("span",{staticClass:"line-number"},[s._v("27")]),a("br"),a("span",{staticClass:"line-number"},[s._v("28")]),a("br"),a("span",{staticClass:"line-number"},[s._v("29")]),a("br"),a("span",{staticClass:"line-number"},[s._v("30")]),a("br"),a("span",{staticClass:"line-number"},[s._v("31")]),a("br"),a("span",{staticClass:"line-number"},[s._v("32")]),a("br"),a("span",{staticClass:"line-number"},[s._v("33")]),a("br"),a("span",{staticClass:"line-number"},[s._v("34")]),a("br"),a("span",{staticClass:"line-number"},[s._v("35")]),a("br"),a("span",{staticClass:"line-number"},[s._v("36")]),a("br"),a("span",{staticClass:"line-number"},[s._v("37")]),a("br"),a("span",{staticClass:"line-number"},[s._v("38")]),a("br"),a("span",{staticClass:"line-number"},[s._v("39")]),a("br")])]),a("p",[s._v("The course's struct in the "),a("code",[s._v("course.go")]),s._v(" looks like this:")]),s._v(" "),a("div",{staticClass:"language-go line-numbers-mode"},[a("pre",{pre:!0,attrs:{class:"language-go"}},[a("code",[a("span",{pre:!0,attrs:{class:"token comment"}},[s._v("// Course defines the Course model")]),s._v("\n"),a("span",{pre:!0,attrs:{class:"token keyword"}},[s._v("type")]),s._v(" Course "),a("span",{pre:!0,attrs:{class:"token keyword"}},[s._v("struct")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("{")]),s._v("\n\tID          "),a("span",{pre:!0,attrs:{class:"token builtin"}},[s._v("string")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('`json:"id"`')]),s._v("\n\tName        "),a("span",{pre:!0,attrs:{class:"token builtin"}},[s._v("string")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('`json:"name"`')]),s._v("\n\tSubtitle    "),a("span",{pre:!0,attrs:{class:"token builtin"}},[s._v("string")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('`json:"subtitle"`')]),s._v("\n\tDescription "),a("span",{pre:!0,attrs:{class:"token builtin"}},[s._v("string")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('`json:"description"`')]),s._v("\n"),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("}")]),s._v("\n")])]),s._v(" "),a("div",{staticClass:"line-numbers-wrapper"},[a("span",{staticClass:"line-number"},[s._v("1")]),a("br"),a("span",{staticClass:"line-number"},[s._v("2")]),a("br"),a("span",{staticClass:"line-number"},[s._v("3")]),a("br"),a("span",{staticClass:"line-number"},[s._v("4")]),a("br"),a("span",{staticClass:"line-number"},[s._v("5")]),a("br"),a("span",{staticClass:"line-number"},[s._v("6")]),a("br"),a("span",{staticClass:"line-number"},[s._v("7")]),a("br")])]),a("h3",{attrs:{id:"identifying-attribute"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#identifying-attribute","aria-hidden":"true"}},[s._v("#")]),s._v(" Identifying Attribute")]),s._v(" "),a("p",[s._v("As you can see an "),a("code",[s._v("ID")]),s._v(" attribute of type "),a("code",[s._v("string")]),s._v(" is automatically added. This "),a("code",[s._v("ID")]),s._v(" will also be used as the hash key for the table.")]),s._v(" "),a("p",[s._v("It is important to note, that DynamoDB is not best friends with "),a("strong",[s._v("UUID")]),s._v(" types. The marshall/ unmarshall process is still a little buggy, so at this point it is safer to just store "),a("strong",[s._v("UUIDs")]),s._v(" as strings and marshall/ unmarshall them in the functions code, instead of letting DynamoDB do that work.")]),s._v(" "),a("p",[s._v("If wish to define your own identifying attribute, you can do so by specifying the required flags of the command. For details see the "),a("router-link",{attrs:{to:"/commands/"}},[s._v("Commands Reference")]),s._v(".")],1),s._v(" "),a("h3",{attrs:{id:"general-attributes"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#general-attributes","aria-hidden":"true"}},[s._v("#")]),s._v(" General Attributes")]),s._v(" "),a("p",[s._v("The type for an attribute is by default "),a("code",[s._v("string")]),s._v(" if none is provided. E.g., if you would like to add a number just provide the golang type with the attribute's name separated by a colon (😃 like this:")]),s._v(" "),a("div",{staticClass:"language- line-numbers-mode"},[a("pre",{pre:!0,attrs:{class:"language-text"}},[a("code",[s._v('mug add resource course -a "name,subtitle,description,price:float32"\n')])]),s._v(" "),a("div",{staticClass:"line-numbers-wrapper"},[a("span",{staticClass:"line-number"},[s._v("1")]),a("br")])]),a("p",[s._v("The "),a("code",[s._v("course.go")]),s._v(" file will also contain all the wrapper methods to interact with the database, just in case you want to have another solution as DynamoDB as persistence layer. The "),a("code",[s._v("main.go")]),s._v(" files in the function subdirectories will contain the actual lambda functionality.")]),s._v(" "),a("h2",{attrs:{id:"complex-resource-definition-with-nested-objects"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#complex-resource-definition-with-nested-objects","aria-hidden":"true"}},[s._v("#")]),s._v(" Complex Resource Definition with Nested Objects")]),s._v(" "),a("p",[s._v("With Dynamo DB being a NoSQL database you certainly cannot use relationships like you may be used to from relational databases like MySQL or PostgreSQL. Usually you overcome this by deciding which entities you work with (querying, writing, etc.) and embedding all related information.")]),s._v(" "),a("p",[s._v("With mug you can easily generate such resources with a little more complex attribute definition to the "),a("code",[s._v("-a")]),s._v(" attributes flag:")]),s._v(" "),a("ul",[a("li",[a("code",[s._v("{}")]),s._v(" wraps nested objects ("),a("code",[s._v("1-1")]),s._v(" relationship)")]),s._v(" "),a("li",[a("code",[s._v("[]")]),s._v(" wraps slices of objects ("),a("code",[s._v("1-n")]),s._v(" relationship).")])]),s._v(" "),a("p",[s._v("For example, let's generate a user, who has an address and may have multiple enrollments to courses: "),a("strong",[s._v("Please make sure to wrap the attribute definition with "),a("code",[s._v('"')]),s._v("(double quotes) to ensure the recursive parsing works!")])]),s._v(" "),a("div",{staticClass:"language- line-numbers-mode"},[a("pre",{pre:!0,attrs:{class:"language-text"}},[a("code",[s._v('mug add resource user -a "name,isActive:bool,email,\naddress:{street,zip,city},\nenrollments:[courseID,startDate:time.Time,endDate:time.Time]"\n')])]),s._v(" "),a("div",{staticClass:"line-numbers-wrapper"},[a("span",{staticClass:"line-number"},[s._v("1")]),a("br"),a("span",{staticClass:"line-number"},[s._v("2")]),a("br"),a("span",{staticClass:"line-number"},[s._v("3")]),a("br")])]),a("p",[s._v("The generated go struct(s) will look like this:")]),s._v(" "),a("div",{staticClass:"language-go line-numbers-mode"},[a("pre",{pre:!0,attrs:{class:"language-go"}},[a("code",[a("span",{pre:!0,attrs:{class:"token comment"}},[s._v("// User defines the User model")]),s._v("\n"),a("span",{pre:!0,attrs:{class:"token keyword"}},[s._v("type")]),s._v(" User "),a("span",{pre:!0,attrs:{class:"token keyword"}},[s._v("struct")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("{")]),s._v("\n\tID       "),a("span",{pre:!0,attrs:{class:"token builtin"}},[s._v("string")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('`json:"id"`')]),s._v("\n\tName     "),a("span",{pre:!0,attrs:{class:"token builtin"}},[s._v("string")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('`json:"name"`')]),s._v("\n\tIsActive "),a("span",{pre:!0,attrs:{class:"token builtin"}},[s._v("bool")]),s._v("   "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('`json:"is_active"`')]),s._v("\n\tEmail    "),a("span",{pre:!0,attrs:{class:"token builtin"}},[s._v("string")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('`json:"email"`')]),s._v("\n\n\tAddress     Address       "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('`json:"address"`')]),s._v("\n\tEnrollments "),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("[")]),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("]")]),s._v("Enrollment  "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('`json:"enrollments"`')]),s._v("\n"),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("}")]),s._v("\n\n"),a("span",{pre:!0,attrs:{class:"token comment"}},[s._v("// Address defines the Address model")]),s._v("\n"),a("span",{pre:!0,attrs:{class:"token keyword"}},[s._v("type")]),s._v(" Address "),a("span",{pre:!0,attrs:{class:"token keyword"}},[s._v("struct")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("{")]),s._v("\n\tStreet "),a("span",{pre:!0,attrs:{class:"token builtin"}},[s._v("string")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('`json:"street"`')]),s._v("\n\tZip    "),a("span",{pre:!0,attrs:{class:"token builtin"}},[s._v("string")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('`json:"zip"`')]),s._v("\n\tCity   "),a("span",{pre:!0,attrs:{class:"token builtin"}},[s._v("string")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('`json:"city"`')]),s._v("\n"),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("}")]),s._v("\n\n"),a("span",{pre:!0,attrs:{class:"token comment"}},[s._v("// Enrollment defines the Enrollments model")]),s._v("\n"),a("span",{pre:!0,attrs:{class:"token keyword"}},[s._v("type")]),s._v(" Enrollment "),a("span",{pre:!0,attrs:{class:"token keyword"}},[s._v("struct")]),s._v(" "),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("{")]),s._v("\n\tCourseID  "),a("span",{pre:!0,attrs:{class:"token builtin"}},[s._v("string")]),s._v("    "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('`json:"course_id"`')]),s._v("\n\tStartDate time"),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v(".")]),s._v("Time "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('`json:"start_date"`')]),s._v("\n\tEndDate   time"),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v(".")]),s._v("Time "),a("span",{pre:!0,attrs:{class:"token string"}},[s._v('`json:"end_date"`')]),s._v("\n"),a("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("}")]),s._v("\n")])]),s._v(" "),a("div",{staticClass:"line-numbers-wrapper"},[a("span",{staticClass:"line-number"},[s._v("1")]),a("br"),a("span",{staticClass:"line-number"},[s._v("2")]),a("br"),a("span",{staticClass:"line-number"},[s._v("3")]),a("br"),a("span",{staticClass:"line-number"},[s._v("4")]),a("br"),a("span",{staticClass:"line-number"},[s._v("5")]),a("br"),a("span",{staticClass:"line-number"},[s._v("6")]),a("br"),a("span",{staticClass:"line-number"},[s._v("7")]),a("br"),a("span",{staticClass:"line-number"},[s._v("8")]),a("br"),a("span",{staticClass:"line-number"},[s._v("9")]),a("br"),a("span",{staticClass:"line-number"},[s._v("10")]),a("br"),a("span",{staticClass:"line-number"},[s._v("11")]),a("br"),a("span",{staticClass:"line-number"},[s._v("12")]),a("br"),a("span",{staticClass:"line-number"},[s._v("13")]),a("br"),a("span",{staticClass:"line-number"},[s._v("14")]),a("br"),a("span",{staticClass:"line-number"},[s._v("15")]),a("br"),a("span",{staticClass:"line-number"},[s._v("16")]),a("br"),a("span",{staticClass:"line-number"},[s._v("17")]),a("br"),a("span",{staticClass:"line-number"},[s._v("18")]),a("br"),a("span",{staticClass:"line-number"},[s._v("19")]),a("br"),a("span",{staticClass:"line-number"},[s._v("20")]),a("br"),a("span",{staticClass:"line-number"},[s._v("21")]),a("br"),a("span",{staticClass:"line-number"},[s._v("22")]),a("br"),a("span",{staticClass:"line-number"},[s._v("23")]),a("br"),a("span",{staticClass:"line-number"},[s._v("24")]),a("br")])]),a("p",[s._v("Of course, you can also reference the "),a("code",[s._v("ID")]),s._v(" in an object, however, you will have to manage this in your own code - getting the "),a("code",[s._v("ID")]),s._v(" or multiple "),a("code",[s._v("IDs")]),s._v(" for "),a("code",[s._v("m-n")]),s._v(" relationships and fetching the referenced objects afterwards.")]),s._v(" "),a("h2",{attrs:{id:"adding-function-groups-functions-etc"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#adding-function-groups-functions-etc","aria-hidden":"true"}},[s._v("#")]),s._v(" Adding Function Groups, Functions etc.")]),s._v(" "),a("p",[s._v("Adding function groups or functions is very similar. You can also remove resources, function groups or functions from the previous.")]),s._v(" "),a("p",[a("strong",[s._v("Have a look at the appropriate command syntax in the "),a("router-link",{attrs:{to:"/commands/"}},[s._v("Commands Reference")]),s._v(".")],1)])])},[],!1,null,null,null);t.default=n.exports}}]);