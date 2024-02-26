EasyUptime is a tool written in Go to monitor uptime of a website, returing the http response status code and millisecond response time. It is currently configured using SQLite as the database and GORM, but any database GORM supports can be used. The only other package used is httprouter for the API endpoints.

A very basic frontend written in React is included as well, but you can also use Postman or Bruno, etc... to hit the endpoints.

TODO:
- Add tests
- Add graphs for the frontend - uptime trends
- Pagination
- Write a scheduler to run the getDomainStatus function on an interval