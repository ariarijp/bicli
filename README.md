# bicli

Create short URLs with Bitly API.

## Usage

```shell
$ bicli -h
Usage of bicli:
  -conf string
    	Config file (default "config.toml")
  -init
    	Create config file
  -sep string
    	Output separator (default ",")
  -sleep-msec uint
    	Sleep time for each request (default 1000)
  -urls string
    	Long URL urls file (default "urls.csv")

$ cat config.toml
login = "YOUR_BITLY_LOGIN_NAME"
api-key = "YOUR_BITLY_API_KEY"

$ cat urls.csv.example
https://google.com
https://www.facebook.com
https://twitter.com

$ bicli -urls=urls.csv.example
1,http://bit.ly/2jMLb5C,https://google.com
2,http://bit.ly/2jMRdTR,https://www.facebook.com
3,http://bit.ly/2jMVG90,https://twitter.com
```

## License

MIT

## Author

[ariarijp](https://github.com/ariarijp)
