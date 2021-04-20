## Demo
___

There's a demo hosted currently hosted at the following url:
```
https://elbwgrs.space
```

There's some data that has already been preseeded. To access the modification related routes, you'll need an api key. Here's one:
```
UY94nF1VsfMZp+Uv
```


## Endpoints
___

All modification routes are protected. To create or delete comments, you'll need to provide the api key in the Authorisation header with the bearer prefix like so:

```
Authorization: Bearer <api_key>
```

POST /org/:org/comments/\
- Creates a comment scoped to an organisation.
GET /org/:org/comments/\
- Returns all comments scoped to an organisation.
GET /org/:org/members/\
- Returns the list of members in an organisation ordered by the number of followers in descending order.
DELETE /org/:org/comments/\
- Deletes all comments under an organisation.

## Schema
___
### Organisations
Stores the organisations.\
Has the following columns:
* id
   * serial, primary key.
* name
   * text, not null unique.
   * Name of the organisation.
* created_at
   * timestamptz, defaults to current_timestamp.
   * When the row was created.
* updated_at
   * timestamptz, defaults to current_timestamp.
   * When the row was updated.
* deleted_at
   * timestamptz, nullable.
   * When the row was deleted.

### Users
Stores the users.\
Has the following columns:
* id
   * serial, primary key.
* username
   * text, not null unique.
   * Usename of the user.
* password_hash
   * text, not null.
   * Bcrypt hash of the user's password
* api_key
   * text, not null.
   * User's api key for access to the api.
* created_at
   * timestamptz, defaults to current_timestamp.
   * When the row was created.
* deleted_at
   * timestamptz, nullable.
   * When the row was deleted.

### OrganisationComments
Stores the comments associated with an organisation.\
Has the following columns:
* id
   * serial, primary key.
* comment
   * text, not null.
   * The comment.
* organisation_id
   * int, not null references organisations(id).
   * The organisation this comment is associated with.
* created_by
   * int, not null references user(id).
   * The user that posted this comment
* created_at
   * timestamptz, defaults to current_timestamp.
   * When the row was created.
* deleted_at
   * timestamptz, nullable.
   * When the row was deleted.

### UserOrganisations
Stores the user and the organisations they're tied to.\
Has the following columns:
* user_id
   * int, not null references user(id).
   * The user.
* organisation_id
   * int, not null references organisations(id).
   * The organisation.
* created_at
   * timestamptz, defaults to current_timestamp.
   * When the row was created.
* The primary key is a composite of (user_id, organisation_id). Is unique.

### UserFollowers
Stores the users and their followers.\
Has the following columns:
* followee_id
   * int, not null references user(id).
   * The user that is being followed.
* follower_id
   * int, not null references user(id).
   * The user that is following.
* The primary key is a composite of (follower_id, followee_id). Is unique.
* Has check constraint (followee_id != follower_id). This prevents a scenario where a user follows themselves.

## Project Structure
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
* [mockery] (https://github.com/vektra/mockery) - automatic mock generation for testing
* [migrate] (https://github.com/golang-migrate/migrate) - schema migrations

## Infrastructure
___

Currently the application is being deployed in an elastic beanstalk through a docker compose file. In there, there's 3 components:
1. A Proxy, we're using Traefik in this case.
2. A members app
3. A comments app
To speed up the existing development process, the identical application is used in both services. Path based routing is employed on the proxy end to route the traffic to the relevant services. The proxy also handles TLS termination.

The applications are configured through environment variables that are being set through 2 avenues. Elastic beanstalk environment settings for options that aren't secrets. Secrets are stored on AWS SSM parameter store and encrypted with KMS. The instance profiles have the relevant permissions to fetch them and it is done during start up. The database used is Postgres hosted on RDS.

## How to deploy to Production
___

First, build and tag the image. We're using the a stable tag for convenience, i.e latest.
```
docker build -t 461960015384.dkr.ecr.ap-southeast-1.amazonaws.com/github:latest .
```
We're using ECR to store the images.

Then run the `create-app-bundle.sh` script in scripts like so.
```
./scripts/create-app-bundle.sh
```
The scripts does the following:
1. Copy production compose file into a tmp folder.
2. Zips it up. The produced zip file will have the current commit hash.
3. Uploads it to S3.
4. Registers the version with the App created on elastic beanstalk.

After the bundle is uploaded. Login to elastic beanstalk, find the commit hash of the version you want to deploy and deploy it to the production environment.

## How to run the tests
___

Inside the `docker-compose.yml` file, there's a service named `test-db`. Start `test-db` like so:
```
docker-compose up -d test-db
```

Run the migrations like so:
```
migrate -path migrations -database "postgres://postgres:postgres@localhost:5433/main_test?sslmode=disable" up 
```

To run the testsuite, run this command:
```
DB_URL="postgres://postgres:postgres@localhost:5433/main_test?sslmode=disable" go test -v ./...  
```
