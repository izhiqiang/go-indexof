# Purpose

Map the local directory to the web side. If it is a file, you can download it directly. If it is a directory, you can enter the directory

# Using nginx

~~~
worker_processes auto;
events {
    worker_connections 1024;
}
http {
sendfile            on;
    tcp_nopush          on;
    tcp_nodelay         on;
    keepalive_timeout   65;
    types_hash_max_size 2048;
    default_type        application/octet-stream;
    server {
        listen       80 default_server;
        listen       [::]:80 default_server;
        server_name  _;
        location / {
            # 配置访问目录
            root    /;
            # 开启索引功能
            autoindex on;
            # off关闭计算文件确切大小（单位bytes），只显示大概大小（单位kb、mb、gb）
            autoindex_exact_size off;
            #显示本机时间而非 GMT 时间
            autoindex_localtime on;
        }
        error_page 404 /404.html;
            location = /40x.html {
        }
        error_page 500 502 503 504 /50x.html;
            location = /50x.html {
        }
    }
}
~~~

# Using go-indexof

~~~
git clone https://github.com/izhiqiang/go-indexof.git
cd go-indexof
vi config.json
make
chmod +x indexof
./indexof -f config.json 
~~~

# commands

~~~
//View process id
ps aux | grep indexof | grep -v grep | awk '{print $2}'

//Enable daemon startup
nohup ./indexof -f config.json > indexof.log &

//Find the process id and stop it
ps aux | grep indexof | grep -v grep | awk '{print $2}' | xargs kill -TERM
~~~

# config.json

~~~
cat > /etc/indexof/config.json << EOF
{
  "port": 8080,
  "index_of": {
    "name": "local",
    "root": "./"
  }
}
EOF
~~~

# Add startup item

~~~
cat > /etc/systemd/system/indexof.service << EOF
[Unit]
Description=Start indexof
[Service]
ExecStart= you path/indexof -f config.json
[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl start indexof
~~~

