# Gitolite conf parser validate

 ! Work in progress
 
 The parser takes Gitolite configuration, parses it and should allow validation and optimization. Along with these goals this also served as a test for writing a lexer and a parser in Go.
 
 ## Usage
 
    cat gitolite.conf | go run cmd/validator/main.go
 
 ## License 
 
 MIT License