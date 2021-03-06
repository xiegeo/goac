# goAC

### Why?
I need to replace a static user access control database file with something users can self serve to share permissions.

### How does control propagate in a graph?

![digraph G {
	A [shape=box];B [shape=box];C [shape=box];D [shape=box];
	A -> B ; A -> C ; "-g" -> C; 
	B -> f;
	C -> D; D -> g[style=dotted];
	A -> f [color=blue]; A -> g [color=blue]; 
	C -> g [color=blue,style=dotted];
}](http://g.gravizo.com/svg?digraph%20G%20%7B%0A%09A%20%5Bshape%3Dbox%5D%3BB%20%5Bshape%3Dbox%5D%3BC%20%5Bshape%3Dbox%5D%3BD%20%5Bshape%3Dbox%5D%3B%0A%09A%20-%3E%20B%20%3B%20A%20-%3E%20C%20%3B%20%22-g%22%20-%3E%20C%3B%20%0A%09B%20-%3E%20f%3B%0A%09C%20-%3E%20D%3B%20D%20-%3E%20g%5Bstyle%3Ddotted%5D%3B%0A%09A%20-%3E%20f%20%5Bcolor%3Dblue%5D%3B%20A%20-%3E%20g%20%5Bcolor%3Dblue%5D%3B%20%0A%09C%20-%3E%20g%20%5Bcolor%3Dblue%2Cstyle%3Ddotted%5D%3B)

Vertices A, B, C, and D are users; f and g are allowed actions; and -g (not g) denies access to g. Black arrows are explicit permission assignments, and blue arrows are derived (not all shown). Dotted arrows are disabled (by -g). 

A, B, C ... does not all have to be users, they can also model any other organizational unit. "A -> B" means "A have control over B" or "permissions of A is a superset of B" (within reason, non transferable permissions should not be modeled. For example, A can not pretend to be B and have actions recorded as B did them). 

Rules:
- Allow propagate up. ( A -> B, B -> f, therefore A -> f)
- Deny propagate down. ( -g -> C, C -> D, therefore -g -> D)
- Deny if not allowed.
- Deny disables allow.

Edge cases:
- Allow propagate up even when blocked. ( A -> C, -g -> C, C -> g, therefore C -> g is disabled, but A -> g  is allowed.) If denying a higher user is intended, then that should be set explicitly.
- It is possible to create cycles in the graph. There is no guarantee against doing so when modifying the graph to allow partial updates without locking the whole graph. By rule, vertices in a cycle will have the same set of permissions. Graph search algorithms must remember visited to avoid going in circles.



It will do:
- Role-based access control, without the complexity. It is all many to many relationships anyways, so the algorithm just need a directed graph. 
- Discretionary access control, allowed actions move up the graph. Capabilities is also core to how it works and explained later.
- Mandatory access control, blocked action move down the graph.

It does not do:
- Authenticate users. There are many ways to Authenticate user actions, this library only require that names be unique identifiers. Pre and post authentication function calls should be kept separate, only post auth APIs exist in this library.
- Provisioning. This library should make it easy to do, but does not implement any network protocols itself.
- Auditing. Again, just make it easy by exposing the right APIs.

Serialization from and to JSON files so that existing tools that provision and audit text files can be used.

## Capabilities / Users Sharing Permissions

	[{
		"name": "Admin",
		"assignments": [{
			"elevate": "Alice",
			"over": "g",
			"comments": {
				"note": "this tells future you why this assignment is made.",
				"createdOn": "2016.02.02",
				"foobar": "comments is a map of string to string, a good place to add metadata for auditing"
			}
		}]
	}, {
		"name": "Alice",
		"assignments": [{
			"elevate": "Bob",
			"over": "g"
		}, {
			"elevate": "-g",
			"over": "Bob"
		}]
	}]
	
The above is a concatenation of 2 files, supplied by users. Admin is a built-in all powerful account, it gives the control of g over to Alice. Alice also shares g with Bob using her first assignment. Notice that Alice also tries to disable g for Bob, but fails as Alice does not have control over Bob. In this contrived example, g end up controled by Admin, Alice, and Bob. If Admin also assigns Alice over Bob, then Bob will be blocked from g as deny rules have priority.

### Parameterized

[See wiki](https://github.com/xiegeo/goac/wiki/Parameterized-Assignments)

### Usage
![actor User;
participant "Authentication" as A;
participant "goAC" as G;
participant "App Internals" as I;
User -> A: TLS protected request;
activate A;
A -> G: Request verified to User;
activate G;
G -> A: Allow or Deny request;
deactivate G;
A -> User: if Denied, return error;
A -> I: if Allowed, send request;
activate I;
hide footbox;](http://g.gravizo.com/svg?actor%20User%3B%0Aparticipant%20%22Authentication%22%20as%20A%3B%0Aparticipant%20%22goAC%22%20as%20G%3B%0Aparticipant%20%22App%20Internals%22%20as%20I%3B%0AUser%20-%3E%20A%3A%20TLS%20protected%20request%3B%0Aactivate%20A%3B%0AA%20-%3E%20G%3A%20Request%20verified%20to%20User%3B%0Aactivate%20G%3B%0AG%20-%3E%20A%3A%20Allow%20or%20Deny%20request%3B%0Adeactivate%20G%3B%0AA%20-%3E%20User%3A%20if%20Denied%2C%20return%20error%3B%0AA%20-%3E%20I%3A%20if%20Allowed%2C%20send%20request%3B%0Aactivate%20I%3B%0Ahide%20footbox%3B
)

see [goac_test.go](goac_test.go) for examples

### User Self Serve Interface
- Update user's JSON config file. A graphical user interface (GUI) that shows an editable graph is good to have.
- A GUI that shows the user it's relevent subgraph, including allow and deny rules others have placed. It lets the user see and test who he have control over, and who among them can do what.

The GUI will be an Ajax web app. It will start with JSON text file edit, and graph properties viewing by tables. Actural editable graphs will be expermented on later.

## Related Reading

https://github.com/mikespook/gorbac (no self serve, no permission sharing by users)

https://en.wikipedia.org/wiki/Role-based_access_control

https://wiki.evolveum.com/display/midPoint/Identity+Management
