gosh is a toy shell I implemented to play with Golang


## Requirements

```
go get -u github.com/nsf/termbox-go
```

## Features
- Builtin functions: `ciao`, `bye`
- Suggestion of command options
- Possibility to extend suggestions to other commands

## Support new command suggestions

Put a file named as the command you want to extend and follow the format of the provided commands, i.e. [ls][ls-config] 
## License

This project is release under CC0 license ( Creative Commons Zero)

The termbox-go library and Golang follows its own licenses


[ls-config]: configs/ls
