#!/bin/sh
set -eu

project_dir=$(CDPATH= cd -- "$(dirname -- "$0")/../../../.." && pwd)
build_dir=$(mktemp -d)
trap 'rm -rf "$build_dir"' EXIT HUP INT TERM

case $(uname -s) in
    Darwin) extension=dylib ;;
    *) extension=so ;;
esac

cd "$project_dir"
go build -mod=mod -buildmode=c-shared -o "$build_dir/libfrontage.$extension" ./lib/client/io
cc -std=c11 -Wall -Wextra -Werror -fsanitize=address \
    -I "$build_dir" lib/client/io/testdata/free_views.c \
    -L "$build_dir" -lfrontage -Wl,-rpath,"$build_dir" \
    -o "$build_dir/free_views"
DYLD_LIBRARY_PATH="$build_dir" "$build_dir/free_views"
