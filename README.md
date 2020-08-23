# Twitter Contest REST API
A restful service which fetches data from twitter apis and runs some business logic to identify most active retweeter

## Instructions to run the API:
The API server is configured to listen to the port **8080**
### Create a Twitter App

To create a Twitter App, head over to <http://apps.twitter.com/>, sign in, then create a new app.

### Authenticate with the Twitter API

Before you use the REST API, you need to authenticate with the Twitter API. To authenticate with the Twitter API (Application-only Authentication), head to the "Keys and Access Tokens" tab and get your Consumer Key and Secret.

|                 |                          |
|-----------------|--------------------------|
| Consumer Key    | `Lx.... your key here`   |
| Consumer Secret | `ED... your secret here` |

Now place the obtained Consumer Key and Consumer Secret in the TODO.keys.json file and follow the instructions mentioned in that file.


#### Endpoints:
- ```sh
  GET /twitter/retweets/:user_handle/max  
  ```
  Picks up winners who has retweeted the most number of times across the last 100 tweets of a given Twitter Handle and then prints out the winner's username and the   count of retweets.
  
- ```sh
  GET /twitter/tweet/:user_handle/latest  
  ```
  Prints out the latest tweets of a given Twitter Handle.
  
### To run the application using docker
* Requirements : `Docker`
* #### RUN `docker build --tag twittercontest .`
* #### RUN `docker run -it -p 8080:8080 twittercontest`
* The Application is up and running at [localhost:8080](http://localhost:8080)   

### To run the application locally
* Requirements : `Golang`

* `cd` to the root directory
* #### RUN  `go mod download`
* #### RUN  `go run main.go`
* The Application is up and running at [localhost:8080](http://localhost:8080) 
