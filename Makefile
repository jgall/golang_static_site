

build:
	docker build -t static_file_server:latest .
run:
	docker run -d -v /static_web:/public -p 80:80 -p 443:443 static_file_server:latest --dir=/public --host=${HOST}
