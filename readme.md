# go file watcher

watch directory and run command

## 格式

`go_file_watcher --config=test.conf`

## test.conf例子

```
{"command_list":[
    {
        "pattern": "*.py",
        "command": "sudo supervisorctl restart redis"
    },
    {
        "command": "echo things changed"
    }
]
"directory": "/home/vagrant/code/code1"
}
```

