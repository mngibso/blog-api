#blog-api
This is the server side of a blogging application using a RESTful api and Mongodb.
This API enables the consumer to:
- Create, update and delete users 
- Create, update and delete blog posts

Basic authentication is used for protected endpoints.  A user can only update and delete themselves. A user manage only the blog postings they created.

##Setup:

1. An accessible MongoDB is required to run the api.  [MongoDB Atlas](https://www.mongodb.com/cloud/atlas) provides a free tier that allows you to deploy a cloud MongoDB.
    - Once your database is up, you can obtain the connection string from the "Connect"
dialog on the Clusters page.  e.g. `mongodb+srv://blog_api:<password>@cluster0.pmalk.mongodb.net/<dbname>?retryWrites=true&w=majority"`
  
2. Clone the repository

    ```
    $ git clone https://github.com/mngibso/blog-api.git
    ```

3. Run the API. The MONGDB_URI environment variable is required for the application to run.  It should be set to the connection string obtained above.
    ``` 
   $ cd blog-api
   $ MONGODB_URI="mongodb+srv://blog_api:<password>@cluster0.pmalk.mongodb.net/<dbname>?retryWrites=true&w=majority" go run main.go
   ...
   [GIN-debug] Listening and serving HTTP on :8080
    ```

##Run Unit Test

The API was designed to allow unit tests to be written.  I've written one unit test for the create post handler as an example.

```
go test ./...
```
    
##Testing

1. Create user "allen" with password "password"
    ```
   curl --location --request POST 'http://localhost:8080/v1/user' \
    --header 'Content-Type: application/json' \
    --header 'Authorization: Basic Og==' \
    --data-raw '{"username":"allen", "firstName":"Allen", "lastName":"Smith", "email":"al@blog.com", "password":"password"}'
    ```

2. User "allen" creates two blog posts using Basic Authentication

    ```
    curl --location --request POST 'http://localhost:8080/v1/post' \
    --header 'Content-Type: application/json' \
    --header 'Authorization: Basic YWxsZW46cGFzc3dvcmQ=' \
    --data-raw '{
        "title": "Allen'\''s First Blog Post",
        "body": "This is a blog post about allen'\''s stuff."
    }'

    curl --location --request POST 'http://localhost:8080/v1/post' \
    --header 'Content-Type: application/json' \
    --header 'Authorization: Basic YWxsZW46cGFzc3dvcmQ=' \
    --data-raw '{
        "title": "Allen'\''s Second Blog Post",
        "body": "More stuff about Allen Smith"
    }'
    ```

3. Create user "barbra" with password "password".  Add a blog entry for Barbra

    ```
    curl --location --request POST 'http://localhost:8080/v1/user' \
    --header 'Content-Type: application/json' \
    --header 'Authorization: Basic Og==' \
    --data-raw '{"username":"barbra", "firstName":"Barbra", "lastName":"Jones", "email":"barb@blog.com", "password":"password"}'
    
    curl --location --request POST 'http://localhost:8080/v1/post' \
    --header 'Content-Type: application/json' \
    --header 'Authorization: Basic YmFyYnJhOnBhc3N3b3Jk' \
    --data-raw '{
        "title": "All about Barbra",
        "body": "What has Barbra been doing?"
    }'
    ```

4. List all users and blog posts
    ```
    curl --location --request GET 'http://localhost:8080/v1/user'
    curl --location --request GET 'http://localhost:8080/v1/post'
    ```

5. List all Allen's blog posts
    ```
    curl --location --request GET 'http://localhost:8080/v1/post?username=allen'
    ```

6. Barbra updates Barbra's blog post and views the changes 

    ```
    curl --location --request PUT 'localhost:8080/v1/post/<postID>' \
    --header 'Content-Type: application/json' \
    --header 'Authorization: Basic YmFyYnJhOnBhc3N3b3Jk' \
    --data-raw '{"title":"ALL ABOUT BARBRA!!!!","username":"barbra","body":"The story of Barbra'\''s Life...." ,"createdAt":1594242530 ,"id":"<postID>"}'

    curl --location --request GET 'http://localhost:8080/v1/post<postID>'
    ```

7. Allen tries to update Barbra's blog post, results in a 403
    ```
    curl --location --request PUT 'localhost:8080/v1/post/5f076f99b51e987b586c1d6a' \
    --header 'Content-Type: application/json' \
    --header 'Authorization: Basic YWxsZW46cGFzc3dvcmQ=' \
    --data-raw '{"title":"HACKED!!!","username":"barbra","body":"HACKED!!!" ,"createdAt":1594242530 ,"id":"5f076f99b51e987b586c1d6a"}'
    ```
   
8. Allen deletes his user.  When a user is deleted, all of the user's blog posts are removed as well.

    ```
    curl --location --request DELETE 'http://localhost:8080/v1/user/allen' \
    --header 'Authorization: Basic YWxsZW46cGFzc3dvcmQ='

    curl --location --request GET 'http://localhost:8080/v1/post'
    ```
