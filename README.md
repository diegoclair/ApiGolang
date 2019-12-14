# ApiGolang
ApiGolang project RV<br>


# Start:
<b>go run main.go</b><br>

# Docker:
If you are using docker, you need to uncomment <b>#DB_HOST=fullstack-postgres</b> in the <b>.env</b> file and comment the <b>DB_HOST=localhost</b>.<br><br>
Then run: <b>docker-compose up</b>

# Using API

## 1 - First of all you need to create a new user:

You can use a Postman to send a "POST" Json to http://localhost:8080/users, as the following example:<br><br>
{<br>
	"fullname":"User testing test",<br>
	"password":"password",<br>
	"birthdate":"1990-12-14",<br>
	"email":"user@test.com"<br>
}<br>
## 2 - Now you need to login to get the JWT token:
Send a "POST" Json to http://localhost:8080/login, as the following example:<br><br>
{	<br>
	"email": "user@test.com",<br>
	"password": "password"<br>
}<br><br>
Then you'll get the token, something like that: <br>
<b>eyJhbGclOiJIUsI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3HpemVkIjp0cnVlLCJ1c2VyX2lkIjozfQ.PBHX7OM9yzJaiYWYoRoXWChDYp-5zL38jxaDnDLudLs</b>
<br><br>
## 3 - Put the header auhotization:
Now you need to put the header authorization with Type Bearer Token, then put the token that you got in the step 2.
<br>

## 4 - Use the following urls to use full API:
### GET:
http://localhost:8080/users<br>
http://localhost:8080/buys<br>
http://localhost:8080/sales<br><br>

### GET by id:
http://localhost:8080/users/{id}<br>
http://localhost:8080/buys/{id}<br>
http://localhost:8080/sales/{id}<br>
To get all buys and sales by user:<br>
http://localhost:8080/reports/id/{id}<br><br>

### GET by date
To get all buys and sales by some date like 2019-12-14:<br>
http://localhost:8080/reports/date/{date}<br><br>

### POST:

#### Users:
http://localhost:8080/users<br> 
To users send a Json like that: The same way as the step 1.<br><br>

#### Buys:
http://localhost:8080/buys<br>
To buys send a Json like that:<br>
{	<br>
	"bitcoin_amount": "0.45",<br>
	"author_id": 1<br>
}<br>
<b>Observe that the "author_id" need to be the same that you got when you create the user</b><br><br>

#### Sales:
http://localhost:8080/sales<br><br>
To sales send a Json like that:<br>
{	<br>
	"bitcoin_amount": "0.45",<br>
	"author_id": 1<br>
}<br>
<b>Observe that the "author_id" need to be the same that you got when you create the user</b><br><br>



# Credits
https://levelup.gitconnected.com/crud-restful-api-with-go-gorm-jwt-postgres-mysql-and-testing-460a85ab7121
