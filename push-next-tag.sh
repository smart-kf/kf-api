#!/bin/bash

VERSION=$(git describe --tags --abbrev=0)
echo "current tag: $VERSION"


# 提取主版本号、次版本号和修订号
MAJOR=$(echo $VERSION | awk -F. '{print $1}' | sed 's/v//')
MINOR=$(echo $VERSION | awk -F. '{print $2}')
PATCH=$(echo $VERSION | awk -F. '{print $3}')



# 增加修订号
PATCH=$((PATCH + 1))

# 组合新的版本号
NEW_VERSION="v$MAJOR.$MINOR.$PATCH"


echo "new tag $NEW_VERSION"

git tag $NEW_VERSION
git push -u origin $NEW_VERSION