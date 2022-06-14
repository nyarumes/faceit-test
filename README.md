# faceit-test

Run it locally (with golang on machine):
- make run

Run it dockerized:
- docker run ...

Run tests:
- make tests

Run test coverage:
- make tests/coverage

Thoughts about my choices:
- gorilla mux is best choice to be close to std library and be sure we have fully tested http in production
- bolt nice way to start and test our business logic if we wanna test the service locally without cluster options but suppose in view of structure to expand it to other databases
- std rpc is easy way to expose the service in rpc world without supporting another schema (in case if we suppose gRPC as alternative)

Let's think about improves:
- It can expand user domain to subfolders that represents entities as separated packages to get clear view on each entity (handlers, repositores, models) like example: "user can have roles"
- Main can isolate RPC and HTTP init handlers to domains itself to avoid connect on one entity
- When we want to be clustered then repository package should be growes with sub packages to each database we want support and additional layer like "service" to unify access to any count of database we want to support (it started with Bolt to "easy way" start develop business logic without complex solutions)
- Integrate confuguration to be more flexlible
- Expand logger to unify place as start point and use ctx to control trace through requests