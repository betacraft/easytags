# easytags
Easy json/xml Tag generation tool for golang

[![Build Status](https://travis-ci.org/betacraft/easytags.svg?branch=master)](https://travis-ci.org/rainingclouds/easytags)

We generally write Field names in CamelCase (aka pascal case) and we generally want them to be in snake case (camel and pascal case are supported as well) when marshalled to json/xml/sql etc. We use tags for this purpose. But it is a repeatative process which should be automated.

usage :

> easytags {file_name} {tag_name_1:case_1, tag_name_2:case_2}

> example: easytags file.go

You can also use this with go generate
For example - In your source file, write following line

> go:generate easytags $GOFILE json,xml,sql

And run
> go generate

This will go through all the struct declarations in your source files, and add corresponding json/xml/sql tags with field name changed to snake case. If you have already written tag with "-" value, this tool will not change that tag.

Now supports Go modules.

![Screencast with Go Generate](https://media.giphy.com/media/26n6G34sQ4hV8HMgo/giphy.gif)
