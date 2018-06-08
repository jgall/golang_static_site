#!/bin/bash
docker rm -f static_site
docker run -d --rm -p 80:80 -p 443:443 -v /static_web_dir:/public -v /.tlsCacheDir:/.tlsCacheDir --name static_site johnwgallagher/golang_static_site --httpPort=80 --httpsPort=443 --dir=/public --host=droplet.johng.site --https=true --cache-dir /.tlsCacheDir