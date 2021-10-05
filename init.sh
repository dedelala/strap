#!/bin/bash
##Usage of init.sh:
##  init.sh <module-path>
##
##  Reinitializes this repository from scratch using your new module path and
##  pushes an initial commit. Assumes the git url is ssh and an appropriately
##  named bare repository exists. Probably only works with github or gitlab.
##  Feel free to use it or don't. Removes itself on completion.

module=$1

if [[ -z "$module" ]]; then
        awk -F'#' '/^##/{print $3}' "$0"
        exit 2
fi

project=$(basename "$module")
git_url="git@$(sed 's#/#:#' <<< "$module").git"

printf "module: %s\nproject: %s\ngit_url: %s\n...continue? [y/n] " \
        "$module" "$project" "$git_url"

read -r confirm
if [[ $confirm != "y" ]]; then
        printf "will not continue, bye\n"
        exit 1
fi

die() {
        printf "oh noes! init failed at %s\n" $*
        exit 1
}

rm -rf .git
git init --initial-branch=main                              || die "git init"
git remote add origin "$git_url"                            || die "git add origin"

go mod edit -module "$module"                               || die "go mod edit"

sed -i"" -e "s/strap/$project/g" .gitignore index.html      || die "s/strap/$project/"

cat > README.md <<< "# $project"

git add .gitignore README.md go.mod *.go index.html static/ || die "git add files"
git commit -m "initialized from github.com/dedelala/strap"  || die "initial commit"
git push -u origin main                                     || die "git push"

printf "%s successfully initialized\n" "$module"
rm "$0"
