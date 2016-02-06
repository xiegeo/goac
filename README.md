# goAC

### Why?
I need to replace a static user access control database file with something users can self serve to share permissions.

### How does control propagate in a graph?

![example control graph, see raw .md file for DOT formate](http://g.gravizo.com/g? digraph G {
	A [shape=box];B [shape=box];C [shape=box];D [shape=box];
	A -> B ; A -> C ; "-g" -> C; 
	B -> f;
	C -> D; D -> g[style=dotted];
	A -> f [color=blue]; A -> g [color=blue]; 
	C -> g [color=blue,style=dotted];
})

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
![A -> B (p)](http://g.gravizo.com/g?digraph G {
	node [shape=box];
    A -> B [label="p"];
	{ rank=same; A B }
  }
)

"A -> B (p)"  means "A have control over p in B" or "permissions of A is a superset of the intersection between B and p". In parameterized assignments only the stated rules gets passed, this also means that B can have permissions that A does not have, and any deny rules over A does not also cover B.

Only the following parameterization is currently planned for implementation:

	{
		"name": "Alice",
		"assignments": [...],
		"levels": [{
			"elevate": "Carol",
			"over": "g",
			"level": {
				"appName": 2
			}
			"comments":{
				"note":"grant service called appName on server g with level 2 access to Carol"
			}
		}]
	}
	
This style of parameterization is intended to work with access defined by levels (higher is more access) or feature sets (1 is allow, one name per feature).

### Usage
![Shows an user request been authenticatied, then check by goAC, before the request is sent to app internals.](http://g.gravizo.com/g?
actor User;
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
hide footbox;
)

TODO: go library API

### User Self Serve Interface
- Update user's JSON config file. A graphical user interface (GUI) that shows an editable graph is good to have.
- A GUI that shows the user it's relevent subgraph, including allow and deny rules others have placed. It lets the user see and test who he have control over, and who among them can do what.

The GUI will be an Ajax web app. It will start with JSON text file edit, and graph properties viewing by tables. Actural editable graphs will be expermented on later.

## Related Reading

https://github.com/mikespook/gorbac (no self serve, no permission sharing by users)

https://en.wikipedia.org/wiki/Role-based_access_control

https://wiki.evolveum.com/display/midPoint/Identity+Management