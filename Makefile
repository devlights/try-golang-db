prepare: \
	_go_get \
	_download_sqlite3_database

_go_get:
	go mod download
	go install honnef.co/go/tools/cmd/staticcheck@latest

_download_sqlite3_database:
	@if [ ! -e "chinook.db" ]; then\
		wget https://www.sqlitetutorial.net/wp-content/uploads/2018/03/chinook.zip;\
		unzip -o chinook.zip;\
		rm -f chinook.zip;\
	fi