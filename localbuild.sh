#! /bin/bash

set -e 

appName="openlist-with-extensions"
builtAt="$(date +'%F %T %z')"
gitAuthor="jenken827 <jenken827@gmail.com>"
gitCommit=$(git log --pretty=format:"%h" -1)
version=$(git describe --long --tags --dirty --always)
# extra="alist with extensions by jenken827"
webVersion="follow backend"  #$(wget -qO- -t1 -T2 "https://api.github.com/repos/alist-org/alist-web/releases/latest" | grep "tag_name" | head -n 1 | awk -F ":" '{print $2}' | sed 's/\"//g;s/,//g;s/ //g')
ldflags="\
-w -s \
-X 'github.com/OpenListTeam/OpenList/internal/conf.BuiltAt=$builtAt' \
-X 'github.com/OpenListTeam/OpenList/internal/conf.GitAuthor=$gitAuthor' \
-X 'github.com/OpenListTeam/OpenList/internal/conf.GitCommit=$gitCommit' \
-X 'github.com/OpenListTeam/OpenList/internal/conf.Version=$version' \
-X 'github.com/OpenListTeam/OpenList/internal/conf.WebVersion=$webVersion' \
"
# -X 'github.com/alist-org/alist/v3/internal/conf.Extra=$extra' \
go build -ldflags="$ldflags" -o ./openlist

# mv openlist ~/servers/alist/