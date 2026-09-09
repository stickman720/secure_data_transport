<h1>Endpoints</h1>

<h2>aouth:</h2>

<h3>[POST] : /api/v1/register</h3>

<code>
input:
    {
        "username" : "str username",
        "password" : "str password",
        "email"    : "str email"
    }

response:
    status 200 
    {
        "message" : "user created",
        "user"    : "userstruct user" 
    }
    status 409
    {
        "error" : "this username is alredy taken"
    }
    status 422
    {
        "error" : "a user with this email is alredy registered"
    }
    status 500
    {
        "error" : "internal server error"
    }
</code>