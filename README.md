# go-meetup

This is the result of a live-coding at a private meetup. The application uses the Rick and Morty API to retriev a number of random characters.<br />
To try it out, start the application (`startup.sh`) and run `http://localhost:8080/random-characters` and pass the query parameter `count` to specify the number of characters (1 - 50).<br />
Or use user `docker build -t rick-and-morty:1 .` to create a docker image and then `docker run --name rick-and-morty -p 8080:8080 rick-and-morty:1` to start a container.

The application uses the following features:
- HTTP handlers
- Middleware functions
- HTML templates
- Channels
- Structured logging