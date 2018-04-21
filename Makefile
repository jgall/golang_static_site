

build:
	docker build -t static_file_server:latest .
run:
	docker run -d -v "$(pwd)"/test:/public -p 80:8000 static_fs --entry=/public
