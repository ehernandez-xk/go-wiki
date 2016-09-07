# go-wiki

Drone Status

[![Build Status](https://drone.io/github.com/ehernandez-xk/go-wiki/status.png)](https://drone.io/github.com/ehernandez-xk/go-wiki/latest)

This is small wiki that helps to undertand some Go packages
- fmt
- html/template
- io/ioutil
- net/http
- regexp

### Directories
* templates/

  Contains templates that the wiki uses.

* data/

  Contains the pages created by the wiki (edit|view|save)
  
### How to use it
```go run wiki.go```

After run this program you can go to http://localhost:8080/TestPage


*Note*

If you want to see the step by step to how the wiki was created.

```git log --oneline --graph --decorate```

you will see:

* 3f935ae 16 - Added valid HTML for templates.
* 832d2c5 15 - New rootHandler to redirect to /view/home
* 3f3240b 14 - Organize the templates and data in new directories.
* f50fb2f 13 - Create makeHandler, to remove repeated code, removed getTitle and import errors.
* ....
*

checkout to the commit you want to test and run

```git checkout f50fb2f```


This tutorial comes from:

https://golang.org/doc/articles/wiki/
