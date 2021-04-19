
# Project Structure
## Folder Structure
___
* controllers 
    * Contains the endpoints and routing.
* entities 
    * Contains the domain models, the objects that are created after data retrieval from the database.
* middleware 
    * Contains the middleware, any common functionality that is shared between multiple controllers.
* migrations 
    * Contains the datastore migrations, in this case that would be SQL.
* mocks 
    * Contains the interface mocks used for testing.
* repositories 
    * Contains the repositories, components that are used to load and store data to and from a datastore.
* requests 
    * Contains the request structs. These are simple components that are used as an intermediary for incoming requests. There goals are:
        1. To provide structure to incoming requests. 
        2. Perform the first pass validation of incoming data.
* responses 
    * Contains the response structs. Similar in nature to requests. The goals are: 
        1. To provide structure outgoing responses. 
        2. Allow us to decouple entities from their representations that are returned to users. 
* scripts 
    * Contains miscellaneous scripts that are used for testing and deployment.
* server 
    * Contains the app container where the individual components are hooked up and assembled together.
* testing_utils 
    * Contains testing related utility functions. Used as a catch all until a more suitable place is found for the code.
* utils 
    * Contains misc utility functions. Used as a catch all until a more suitable place is found for the code.

## Third Party Libraries
___
These are the third party libraries that are used. A brief description is included.
* [Gin](https://github.com/gin-gonic/gin) - web framework
* [sqlx](https://github.com/jmoiron/sqlx) - sql extensions 
* [testify](https://github.com/stretchr/testify) - testing utilities
* [crypto](https://github.com/golang/crypto) - provide bcrypt support for password hashing
* [pq](https://github.com/lib/pq) - postgres driver
* [uuid](https://github.com/google/uuid) - used to generate API keys, currently only being used for testing.