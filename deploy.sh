#!/usr/bin/env sh

# abort on errors
set -e

# build
npm run docs:build

# navigate into the build output directory
cd docs/.vuepress/dist

git init
git add -A
git commit -m 'update documentation'

git push -f git@github.com:crolly/mug.git master:gh-pages

cd -