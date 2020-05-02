module github.com/b4fun/migrated/example

go 1.14

require (
	github.com/b4fun/migrated v0.0.0-20190310063311-c086b214816a
	github.com/golang-migrate/migrate/v4 v4.6.0
	github.com/mattn/go-sqlite3 v1.10.0
	github.com/spf13/cobra v0.0.7
)

replace github.com/golang/lint v0.0.0-20190409202823-959b441ac422 => github.com/golang/lint v0.0.0-20190409202823-5614ed5bae6fb75893070bdc0996a68765fdd275

replace github.com/b4fun/migrated => ../
