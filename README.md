# http-hasher

For using you should compile this code before:
```sh
go build -o http-hasher
```

After that, you can use it to get hashes of site content:
```sh
./http-hasher google.com

http://google.com 3c538158403e123bf85d187c5a8966b3
```
```sh
./http-hasher https://google.com

https://google.com 552666866ce0f4e1bd6830ae5ebccd0e
```

You can use it with multiple sites:
```sh
./http-hasher https://reddit.com google.com

http://google.com 3ef7ecf85c0fac85b711eb42a0be8423
https://reddit.com 868aefdd124d4c8004d8b6bb13c68cbd
```

Also, you can process requests in parallel with `-parallel N` flag (10 threads by default):
```sh
./http-hasher -parallel 3 reddit.com google.com leetcode.com github.com

http://google.com 3d7d902d50dc3320d1d970e9b64a43df
http://github.com 22de6648b3a1bd5e9e347103674aee65
http://leetcode.com 4508b6fb73d6b7e9ccab47579d5f488b
http://reddit.com bdec04f6cad8652cbcb9360196ddbed8
```

For any help you can use `-h` flag:
```sh
./http-hasher -h

Usage of ./http-hasher:
  -parallel int
        How many parallel requests can be processed (default 10)
```