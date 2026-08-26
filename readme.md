## Instagram ghost followers finder
Lightweight script to find ghost follwers in instagram.

### How to run?
- Login to instagram from a browser, keep network tab open, click on followers
- inspect the requests, find a request like `https://www.instagram.com/api/v1/friendships/<user_id>/followers/?count=12`
- copy the user id from the url, paste it as value for the `userId` variable in `main` fn
- right click on the request -> Copy value -> Copy request headers
- Remove the Accept-Encoding header
- paste the rest of headers as value for the `headers` variable.
- Output will be written to the file `output.txt`