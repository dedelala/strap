# strap

A skeleton go web server with embedded bootstrap and font awesome!


## why

I seem to be making a few of these lately and I thought, why not make a
template.


## what

What's in the box?
- Skeleton go webserver that serves `index.html` as a template (`main.go`)
- Request logger middleware and some stuff to help handle content negotiation (`util.go`)
- [Bootstrap 5.1.0](https://getbootstrap.com/docs/5.0/getting-started/introduction/) bundle version
- [JQuery 3.6.0](https://jquery.com/) slim version
- [Font Awesome Free 5.15.4](https://fontawesome.com/v5.15/how-to-use/on-the-web/setup/hosting-font-awesome-yourself) solid font


## how

1. Sort yourself out a new, empty repository somewhere, let's say I've created `github.com/dedelala/my-cool-project`
1. `git clone https://github.com/dedelala/strap.git my-cool-project`
1. `cd my-cool-project`
1. `./init.sh github.com/dedelala/my-cool-project`

The script assumes you push git over ssh and probably only works with github or
gitlab. Otherwise, you're on your own.
