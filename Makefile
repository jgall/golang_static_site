

build:
	docker build -t static_file_server:latest .
run:
	docker run -v "${pwd}"/test:/public -p 8000:8000 static_file_server:latest --entry=/public