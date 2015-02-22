# easytags
Easy json/xml Tag generation tool for golang

We generally write Field names in CamelCase and we generally want them to be in snake case when marshalled to json/xml/sql etc. We use tags for this purpose. But it is a repeatative process which should be automated. 

usage :

> easytags {file_name} {tag_name} 
>example: easytags file.go json

You can also use this with go generate 
For example - In your source file, write following line 

>go:generate easytags filename.go json

And run
>go generate

This will go through all the struct declarations in your source files, and add corresponding json/xml tags with field name changed to snake case. If you have already written a tag, this tool will not change that tag.
