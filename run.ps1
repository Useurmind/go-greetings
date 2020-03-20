param(
    $password
)

go build ./app
./app.exe "host=localhost port=5432 user=postgres password=$password dbname=gogreeting sslmode=disable"